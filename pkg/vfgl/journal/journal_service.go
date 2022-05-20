package journal

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

type IJournalService interface {
	CreateJournal(shopID string, authUsername string, doc vfgl.Journal) (string, error)
	UpdateJournal(guid string, shopID string, authUsername string, doc vfgl.Journal) error
	DeleteJournal(guid string, shopID string, authUsername string) error
	InfoJournal(guid string, shopID string) (vfgl.JournalInfo, error)
	SearchJournal(shopID string, q string, page int, limit int) ([]vfgl.JournalInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []vfgl.Journal) (models.BulkImport, error)
}

type JournalService struct {
	repo   JournalRepository
	mqRepo JournalMqRepository
}

func NewJournalService(repo JournalRepository, mqRepo JournalMqRepository) JournalService {

	return JournalService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc JournalService) CreateJournal(shopID string, authUsername string, doc vfgl.Journal) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := vfgl.JournalDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Journal = doc

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

func (svc JournalService) UpdateJournal(guid string, shopID string, authUsername string, doc vfgl.Journal) error {

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

	err = svc.repo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc JournalService) DeleteJournal(guid string, shopID string, authUsername string) error {
	err := svc.repo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc JournalService) InfoJournal(guid string, shopID string) (vfgl.JournalInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return vfgl.JournalInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return vfgl.JournalInfo{}, errors.New("document not found")
	}

	return findDoc.JournalInfo, nil

}

func (svc JournalService) SearchJournal(shopID string, q string, page int, limit int) ([]vfgl.JournalInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []vfgl.JournalInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc JournalService) SaveInBatch(shopID string, authUsername string, dataList []vfgl.Journal) (models.BulkImport, error) {

	createDataList := []vfgl.JournalDoc{}
	duplicateDataList := []vfgl.Journal{}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[vfgl.Journal](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Docno)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "docno", itemCodeGuidList)

	if err != nil {
		return models.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Docno)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[vfgl.Journal, vfgl.JournalDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc vfgl.Journal) vfgl.JournalDoc {
			newGuid := utils.NewGUID()

			dataDoc := vfgl.JournalDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Journal = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[vfgl.Journal, vfgl.JournalDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (vfgl.JournalDoc, error) {
			return svc.repo.FindByDocIndentiryGuid(shopID, "docno", guid)
		},
		func(doc vfgl.JournalDoc) bool {
			if doc.Docno != "" {
				return true
			}
			return false
		},
		func(shopID string, authUsername string, data vfgl.Journal, doc vfgl.JournalDoc) error {

			doc.Journal = data
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
		createDataKey = append(createDataKey, doc.Docno)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Docno)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		svc.mqRepo.Update(doc)
		updateDataKey = append(updateDataKey, doc.Docno)
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

func (svc JournalService) getDocIDKey(doc vfgl.Journal) string {
	return doc.Docno
}
