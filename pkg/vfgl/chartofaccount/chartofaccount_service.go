package chartofaccount

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IChartOfAccountService interface {
	Create(shopID string, authUsername string, doc vfgl.ChartOfAccount) (string, error)
	Update(guid string, shopID string, authUsername string, doc vfgl.ChartOfAccount) error
	Delete(guid string, shopID string, authUsername string) error
	Info(guid string, shopID string) (vfgl.ChartOfAccountInfo, error)
	Search(shopID string, q string, page int, limit int) ([]vfgl.ChartOfAccountInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []vfgl.ChartOfAccount) (models.BulkImport, error)
}

type ChartOfAccountService struct {
	repo   ChartOfAccountRepository
	mqRepo ChartOfAccountMQRepository
}

func NewChartOfAccountService(repo ChartOfAccountRepository, mqRepo ChartOfAccountMQRepository) ChartOfAccountService {
	return ChartOfAccountService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc ChartOfAccountService) Create(shopID string, authUsername string, doc vfgl.ChartOfAccount) (string, error) {
	newGuidFixed := utils.NewGUID()

	docData := vfgl.ChartOfAccountDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ChartOfAccount = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.mqRepo.Create(docData)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc ChartOfAccountService) Update(guid string, shopID string, authUsername string, doc vfgl.ChartOfAccount) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ChartOfAccount = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc ChartOfAccountService) Delete(guid string, shopID string, authUsername string) error {
	err := svc.repo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc ChartOfAccountService) Info(guid string, shopID string) (vfgl.ChartOfAccountInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return vfgl.ChartOfAccountInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return vfgl.ChartOfAccountInfo{}, errors.New("document not found")
	}

	return findDoc.ChartOfAccountInfo, nil

}

func (svc ChartOfAccountService) Search(shopID string, q string, page int, limit int) ([]vfgl.ChartOfAccountInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"accountcode",
		"accountname",
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []vfgl.ChartOfAccountInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ChartOfAccountService) SaveInBatch(shopID string, authUsername string, dataList []vfgl.ChartOfAccount) (models.BulkImport, error) {

	createDataList := []vfgl.ChartOfAccountDoc{}
	duplicateDataList := []vfgl.ChartOfAccount{}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[vfgl.ChartOfAccount](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.AccountCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "accountcode", itemCodeGuidList)

	if err != nil {
		return models.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.AccountCode)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[vfgl.ChartOfAccount, vfgl.ChartOfAccountDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc vfgl.ChartOfAccount) vfgl.ChartOfAccountDoc {
			newGuid := utils.NewGUID()

			dataDoc := vfgl.ChartOfAccountDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.ChartOfAccount = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[vfgl.ChartOfAccount, vfgl.ChartOfAccountDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (vfgl.ChartOfAccountDoc, error) {
			return svc.repo.FindByDocIndentiryGuid(shopID, "accountcode", guid)
		},
		func(doc vfgl.ChartOfAccountDoc) bool {
			if doc.AccountCode != "" {
				return true
			}
			return false
		},
		func(shopID string, authUsername string, data vfgl.ChartOfAccount, doc vfgl.ChartOfAccountDoc) error {

			doc.ChartOfAccount = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			err = svc.repo.Update(doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(createDataList)

		if err != nil {
			return models.BulkImport{}, err
		}

		svc.mqRepo.CreateInBatch(createDataList)

		if err != nil {
			return models.BulkImport{}, err
		}
	}

	createDataKey := []string{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.AccountCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.AccountCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		svc.mqRepo.Update(doc)
		updateDataKey = append(updateDataKey, doc.AccountCode)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, svc.getDocIDKey(doc))
	}

	return models.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc ChartOfAccountService) getDocIDKey(doc vfgl.ChartOfAccount) string {
	return doc.AccountCode
}
