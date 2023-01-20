package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/job/models"
	"smlcloudplatform/pkg/job/repositories"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IJobHttpService interface {
	CreateJob(shopID string, authUsername string, doc models.Job) (string, error)
	UpdateJob(shopID string, guid string, authUsername string, doc models.Job) error
	DeleteJob(shopID string, guid string, authUsername string) error
	DeleteJobByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoJob(shopID string, guid string) (models.JobInfo, error)
	SearchJob(shopID string, module string, pageable micromodels.Pageable) ([]models.JobInfo, mongopagination.PaginationData, error)
	SearchJobStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.JobInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Job) (common.BulkImport, error)

	GetModuleName() string
}

type JobHttpService struct {
	repo repositories.IJobRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.JobActivity, models.JobDeleteActivity]
}

func NewJobHttpService(repo repositories.IJobRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *JobHttpService {

	insSvc := &JobHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.JobActivity, models.JobDeleteActivity](repo)

	return insSvc
}

func (svc JobHttpService) CreateJob(shopID string, authUsername string, doc models.Job) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "name", doc.Name)

	if err != nil {
		return "", err
	}

	if findDoc.Name != "" {
		return "", errors.New("Name is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.JobDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Job = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc JobHttpService) UpdateJob(shopID string, guid string, authUsername string, doc models.Job) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Job = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc JobHttpService) DeleteJob(shopID string, guid string, authUsername string) error {

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

	svc.saveMasterSync(shopID)

	return nil
}

func (svc JobHttpService) DeleteJobByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc JobHttpService) InfoJob(shopID string, guid string) (models.JobInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.JobInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.JobInfo{}, errors.New("document not found")
	}

	return findDoc.JobInfo, nil

}

func (svc JobHttpService) SearchJob(shopID string, module string, pageable micromodels.Pageable) ([]models.JobInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"name",
	}

	docList, pagination, err := svc.repo.FindPageJob(shopID, module, map[string]interface{}{}, searchInFields, pageable)

	if err != nil {
		return []models.JobInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc JobHttpService) SearchJobStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.JobInfo, int, error) {
	searchInFields := []string{
		"name",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.JobInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc JobHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Job) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Job](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Name)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "name", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Name)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Job, models.JobDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Job) models.JobDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.JobDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Job = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Job, models.JobDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.JobDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "name", guid)
		},
		func(doc models.JobDoc) bool {
			return doc.Name != ""
		},
		func(shopID string, authUsername string, data models.Job, doc models.JobDoc) error {

			doc.Job = data
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
		createDataKey = append(createDataKey, doc.Name)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Name)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Name)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, svc.getDocIDKey(doc))
	}

	svc.saveMasterSync(shopID)

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc JobHttpService) getDocIDKey(doc models.Job) string {
	return doc.Name
}

func (svc JobHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc JobHttpService) GetModuleName() string {
	return "fileFolder"
}
