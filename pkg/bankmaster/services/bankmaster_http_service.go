package services

import (
	"errors"
	"smlcloudplatform/pkg/bankmaster/models"
	"smlcloudplatform/pkg/bankmaster/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IBankMasterHttpService interface {
	CreateBankMaster(shopID string, authUsername string, doc models.BankMaster) (string, error)
	UpdateBankMaster(guid string, shopID string, authUsername string, doc models.BankMaster) error
	DeleteBankMaster(guid string, shopID string, authUsername string) error
	InfoBankMaster(guid string, shopID string) (models.BankMasterInfo, error)
	SearchBankMaster(shopID string, q string, page int, limit int, sort map[string]int) ([]models.BankMasterInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.BankMaster) (common.BulkImport, error)
}

type BankMasterHttpService struct {
	repo repositories.BankMasterRepository
}

func NewBankMasterHttpService(repo repositories.BankMasterRepository) BankMasterHttpService {

	return BankMasterHttpService{
		repo: repo,
	}
}

func (svc BankMasterHttpService) CreateBankMaster(shopID string, authUsername string, doc models.BankMaster) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentiryGuid(shopID, "bankcode", doc.BankCode)

	if err != nil {
		return "", err
	}

	if findDoc.BankCode != "" {
		return "", errors.New("BankCode is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.BankMasterDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.BankMaster = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc BankMasterHttpService) UpdateBankMaster(guid string, shopID string, authUsername string, doc models.BankMaster) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.BankMaster = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc BankMasterHttpService) DeleteBankMaster(guid string, shopID string, authUsername string) error {

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

	return nil
}

func (svc BankMasterHttpService) InfoBankMaster(guid string, shopID string) (models.BankMasterInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.BankMasterInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.BankMasterInfo{}, errors.New("document not found")
	}

	return findDoc.BankMasterInfo, nil

}

func (svc BankMasterHttpService) SearchBankMaster(shopID string, q string, page int, limit int, sort map[string]int) ([]models.BankMasterInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"bankcode",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.BankMasterInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc BankMasterHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.BankMaster) (common.BulkImport, error) {

	createDataList := []models.BankMasterDoc{}
	duplicateDataList := []models.BankMaster{}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.BankMaster](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.BankCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "bankcode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.BankCode)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[models.BankMaster, models.BankMasterDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.BankMaster) models.BankMasterDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.BankMasterDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.BankMaster = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.BankMaster, models.BankMasterDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.BankMasterDoc, error) {
			return svc.repo.FindByDocIndentiryGuid(shopID, "bankcode", guid)
		},
		func(doc models.BankMasterDoc) bool {
			if doc.BankCode != "" {
				return true
			}
			return false
		},
		func(shopID string, authUsername string, data models.BankMaster, doc models.BankMasterDoc) error {

			doc.BankMaster = data
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

	}

	createDataKey := []string{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.BankCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.BankCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.BankCode)
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

func (svc BankMasterHttpService) getDocIDKey(doc models.BankMaster) string {
	return doc.BankCode
}
