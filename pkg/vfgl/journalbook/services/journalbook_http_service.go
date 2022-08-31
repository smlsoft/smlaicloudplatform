package services

import (
	"errors"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"smlcloudplatform/pkg/vfgl/journalbook/models"
	"smlcloudplatform/pkg/vfgl/journalbook/repositories"
	"time"

	common "smlcloudplatform/pkg/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IJournalBookHttpService interface {
	Create(shopID string, authUsername string, doc models.JournalBook) (string, error)
	Update(guid string, shopID string, authUsername string, doc models.JournalBook) error
	Delete(guid string, shopID string, authUsername string) error
	Info(guid string, shopID string) (models.JournalBookInfo, error)
	Search(shopID string, q string, page int, limit int, sort map[string]int) ([]models.JournalBookInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.JournalBook) (common.BulkImport, error)
}

type JournalBookHttpService struct {
	repo   repositories.JournalBookMongoRepository
	mqRepo repositories.JournalBookMqRepository
}

func NewJournalBookHttpService(repo repositories.JournalBookMongoRepository, mqRepo repositories.JournalBookMqRepository) JournalBookHttpService {

	return JournalBookHttpService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc JournalBookHttpService) Create(shopID string, authUsername string, doc models.JournalBook) (string, error) {

	findDoc, err := svc.repo.FindOne(shopID, map[string]interface{}{
		"code": doc.Code,
	})

	if err != nil {
		return "", err
	}

	if len(findDoc.Code) > 0 {
		return "", errors.New("code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.JournalBookDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.JournalBook = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.mqRepo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc JournalBookHttpService) Update(guid string, shopID string, authUsername string, doc models.JournalBook) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDocCode, err := svc.repo.FindOne(shopID, map[string]interface{}{
		"code": doc.Code,
	})

	if err != nil {
		return err
	}

	if findDoc.Code != doc.Code && len(findDocCode.Code) > 0 {
		return errors.New("code is exists")
	}

	findDoc.JournalBook = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc JournalBookHttpService) Delete(guid string, shopID string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc JournalBookHttpService) Info(guid string, shopID string) (models.JournalBookInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.JournalBookInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.JournalBookInfo{}, errors.New("document not found")
	}

	return findDoc.JournalBookInfo, nil

}

func (svc JournalBookHttpService) Search(shopID string, q string, page int, limit int, sort map[string]int) ([]models.JournalBookInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.JournalBookInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc JournalBookHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.JournalBook) (common.BulkImport, error) {

	createDataList := []models.JournalBookDoc{}
	duplicateDataList := []models.JournalBook{}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.JournalBook](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList = importdata.PreparePayloadData[models.JournalBook, models.JournalBookDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.JournalBook) models.JournalBookDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.JournalBookDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.JournalBook = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.JournalBook, models.JournalBookDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.JournalBookDoc, error) {
			return svc.repo.FindByDocIndentiryGuid(shopID, "code", guid)
		},
		func(doc models.JournalBookDoc) bool {
			if doc.Code != "" {
				return true
			}
			return false
		},
		func(shopID string, authUsername string, data models.JournalBook, doc models.JournalBookDoc) error {

			doc.JournalBook = data
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

func (svc JournalBookHttpService) getDocIDKey(doc models.JournalBook) string {
	return doc.Code
}
