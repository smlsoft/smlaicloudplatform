package services

import (
	"errors"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"smlcloudplatform/pkg/vfgl/accountgroup/models"
	"smlcloudplatform/pkg/vfgl/accountgroup/repositories"
	"time"

	common "smlcloudplatform/pkg/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAccountGroupHttpService interface {
	Create(shopID string, authUsername string, doc models.AccountGroup) (string, error)
	Update(guid string, shopID string, authUsername string, doc models.AccountGroup) error
	Delete(guid string, shopID string, authUsername string) error
	Info(guid string, shopID string) (models.AccountGroupInfo, error)
	Search(shopID string, q string, page int, limit int) ([]models.AccountGroupInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.AccountGroup) (common.BulkImport, error)
}

type AccountGroupHttpService struct {
	repo   repositories.AccountGroupMongoRepository
	mqRepo repositories.AccountGroupMqRepository
}

func NewAccountGroupHttpService(repo repositories.AccountGroupMongoRepository, mqRepo repositories.AccountGroupMqRepository) AccountGroupHttpService {

	return AccountGroupHttpService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc AccountGroupHttpService) Create(shopID string, authUsername string, doc models.AccountGroup) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.AccountGroupDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.AccountGroup = doc

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

func (svc AccountGroupHttpService) Update(guid string, shopID string, authUsername string, doc models.AccountGroup) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.AccountGroup = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc AccountGroupHttpService) Delete(guid string, shopID string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc AccountGroupHttpService) Info(guid string, shopID string) (models.AccountGroupInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.AccountGroupInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.AccountGroupInfo{}, errors.New("document not found")
	}

	return findDoc.AccountGroupInfo, nil

}

func (svc AccountGroupHttpService) Search(shopID string, q string, page int, limit int) ([]models.AccountGroupInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []models.AccountGroupInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc AccountGroupHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.AccountGroup) (common.BulkImport, error) {

	createDataList := []models.AccountGroupDoc{}
	duplicateDataList := []models.AccountGroup{}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.AccountGroup](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Code)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "docno", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.DocNo)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[models.AccountGroup, models.AccountGroupDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.AccountGroup) models.AccountGroupDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.AccountGroupDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.AccountGroup = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.AccountGroup, models.AccountGroupDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.AccountGroupDoc, error) {
			return svc.repo.FindByDocIndentiryGuid(shopID, "code", guid)
		},
		func(doc models.AccountGroupDoc) bool {
			if doc.Code != "" {
				return true
			}
			return false
		},
		func(shopID string, authUsername string, data models.AccountGroup, doc models.AccountGroupDoc) error {

			doc.AccountGroup = data
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
		createDataKey = append(createDataKey, doc.Code)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		svc.mqRepo.Update(doc)
		updateDataKey = append(updateDataKey, doc.Code)
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

func (svc AccountGroupHttpService) getDocIDKey(doc models.AccountGroup) string {
	return doc.Code
}
