package services

import (
	"errors"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"smlcloudplatform/pkg/vfgl/chartofaccount/models"
	"smlcloudplatform/pkg/vfgl/chartofaccount/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IChartOfAccountHttpService interface {
	Create(shopID string, authUsername string, doc models.ChartOfAccount) (string, error)
	Update(guid string, shopID string, authUsername string, doc models.ChartOfAccount) error
	Delete(guid string, shopID string, authUsername string) error
	Info(guid string, shopID string) (models.ChartOfAccountInfo, error)
	Search(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ChartOfAccountInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ChartOfAccount) (common.BulkImport, error)
}

type ChartOfAccountHttpService struct {
	repo   repositories.ChartOfAccountRepository
	mqRepo repositories.ChartOfAccountMQRepository
}

func NewChartOfAccountHttpService(repo repositories.ChartOfAccountRepository, mqRepo repositories.ChartOfAccountMQRepository) ChartOfAccountHttpService {
	return ChartOfAccountHttpService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc ChartOfAccountHttpService) Create(shopID string, authUsername string, doc models.ChartOfAccount) (string, error) {
	newGuidFixed := utils.NewGUID()

	docData := models.ChartOfAccountDoc{}
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

func (svc ChartOfAccountHttpService) Update(guid string, shopID string, authUsername string, doc models.ChartOfAccount) error {

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

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}
	svc.mqRepo.Update(findDoc)
	if err != nil {
		return err
	}
	return nil
}

func (svc ChartOfAccountHttpService) Delete(guid string, shopID string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}
	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.mqRepo.Delete(findDoc)
	if err != nil {
		return err
	}
	return nil
}

func (svc ChartOfAccountHttpService) Info(guid string, shopID string) (models.ChartOfAccountInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ChartOfAccountInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ChartOfAccountInfo{}, errors.New("document not found")
	}

	return findDoc.ChartOfAccountInfo, nil

}

func (svc ChartOfAccountHttpService) Search(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ChartOfAccountInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"accountcode",
		"accountname",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.ChartOfAccountInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ChartOfAccountHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.ChartOfAccount) (common.BulkImport, error) {

	createDataList := []models.ChartOfAccountDoc{}
	duplicateDataList := []models.ChartOfAccount{}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.ChartOfAccount](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.AccountCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "accountcode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.AccountCode)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[models.ChartOfAccount, models.ChartOfAccountDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.ChartOfAccount) models.ChartOfAccountDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ChartOfAccountDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.ChartOfAccount = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.ChartOfAccount, models.ChartOfAccountDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ChartOfAccountDoc, error) {
			return svc.repo.FindByDocIndentiryGuid(shopID, "accountcode", guid)
		},
		func(doc models.ChartOfAccountDoc) bool {
			if doc.AccountCode != "" {
				return true
			}
			return false
		},
		func(shopID string, authUsername string, data models.ChartOfAccount, doc models.ChartOfAccountDoc) error {

			doc.ChartOfAccount = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			err = svc.repo.Update(shopID, doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(createDataList)

		if err != nil {
			return common.BulkImport{}, err
		}

		svc.mqRepo.CreateInBatch(createDataList)

		if err != nil {
			return common.BulkImport{}, err
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

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc ChartOfAccountHttpService) getDocIDKey(doc models.ChartOfAccount) string {
	return doc.AccountCode
}
