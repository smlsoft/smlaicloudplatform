package services

import (
	"errors"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IJournalHttpService interface {
	CreateJournal(shopID string, authUsername string, doc models.Journal) (string, error)
	UpdateJournal(guid string, shopID string, authUsername string, doc models.Journal) error
	DeleteJournal(guid string, shopID string, authUsername string) error
	InfoJournal(guid string, shopID string) (models.JournalInfo, error)
	SearchJournal(shopID string, q string, page int, limit int, sort map[string]int) ([]models.JournalInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Journal) (common.BulkImport, error)
}

type JournalHttpService struct {
	repo   repositories.JournalRepository
	mqRepo repositories.JournalMqRepository
}

func NewJournalHttpService(repo repositories.JournalRepository, mqRepo repositories.JournalMqRepository) JournalHttpService {

	return JournalHttpService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc JournalHttpService) CreateJournal(shopID string, authUsername string, doc models.Journal) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentiryGuid(shopID, "docno", doc.DocNo)

	if err != nil {
		return "", err
	}

	if findDoc.DocNo != "" {
		return "", errors.New("DocNo is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.JournalDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Journal = doc

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

func (svc JournalHttpService) UpdateJournal(guid string, shopID string, authUsername string, doc models.Journal) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Journal = doc

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

func (svc JournalHttpService) DeleteJournal(guid string, shopID string, authUsername string) error {

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

func (svc JournalHttpService) InfoJournal(guid string, shopID string) (models.JournalInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.JournalInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.JournalInfo{}, errors.New("document not found")
	}

	return findDoc.JournalInfo, nil

}

func (svc JournalHttpService) SearchJournal(shopID string, q string, page int, limit int, sort map[string]int) ([]models.JournalInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.JournalInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc JournalHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Journal) (common.BulkImport, error) {

	createDataList := []models.JournalDoc{}
	duplicateDataList := []models.Journal{}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Journal](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.DocNo)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "docno", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.DocNo)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[models.Journal, models.JournalDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Journal) models.JournalDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.JournalDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Journal = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Journal, models.JournalDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.JournalDoc, error) {
			return svc.repo.FindByDocIndentiryGuid(shopID, "docno", guid)
		},
		func(doc models.JournalDoc) bool {
			if doc.DocNo != "" {
				return true
			}
			return false
		},
		func(shopID string, authUsername string, data models.Journal, doc models.JournalDoc) error {

			doc.Journal = data
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
		createDataKey = append(createDataKey, doc.DocNo)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.DocNo)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		svc.mqRepo.Update(doc)
		updateDataKey = append(updateDataKey, doc.DocNo)
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

func (svc JournalHttpService) getDocIDKey(doc models.Journal) string {
	return doc.DocNo
}
