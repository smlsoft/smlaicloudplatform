package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	modelDocumentImage "smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	repositoriesDocumentImage "smlcloudplatform/pkg/documentwarehouse/documentimage/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/task/models"
	"smlcloudplatform/pkg/task/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITaskHttpService interface {
	CreateTask(shopID string, authUsername string, doc models.Task) (string, error)
	UpdateTask(shopID string, guid string, authUsername string, doc models.Task) error
	UpdateTaskStatus(shopID string, guid string, authUsername string, jobStatus int8) error
	DeleteTask(shopID string, guid string, authUsername string) error
	DeleteTaskByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoTask(shopID string, guid string) (models.TaskInfo, error)
	SearchTask(shopID string, module string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.TaskInfo, mongopagination.PaginationData, error)
	SearchTaskStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.TaskInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Task) (common.BulkImport, error)

	TaskStatusReference() map[int]string
	GetTaskReject(shopID string, module string, taskGUID string) ([]models.TaskInfo, error)
}

type TaskHttpService struct {
	repo              repositories.ITaskRepository
	repoDocImageGroup repositoriesDocumentImage.IDocumentImageGroupRepository
	services.ActivityService[models.TaskActivity, models.TaskDeleteActivity]
}

func NewTaskHttpService(repo repositories.ITaskRepository, repoDocImageGroup repositoriesDocumentImage.IDocumentImageGroupRepository) *TaskHttpService {

	insSvc := &TaskHttpService{
		repo:              repo,
		repoDocImageGroup: repoDocImageGroup,
	}

	insSvc.ActivityService = services.NewActivityService[models.TaskActivity, models.TaskDeleteActivity](repo)

	return insSvc
}

func (svc TaskHttpService) TaskStatusReference() map[int]string {
	return map[int]string{
		models.TaskPending:   "Pending",
		models.TaskUplaoded:  "Uploaded",
		models.TaskChecking:  "Checking",
		models.TaskCompleted: "Completed",
	}
}

func (svc TaskHttpService) CreateTask(shopID string, authUsername string, doc models.Task) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "name", doc.Name)

	if err != nil {
		return "", err
	}

	if findDoc.Name != "" {
		return "", errors.New("name is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.TaskDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Task = doc
	docData.Status = models.TaskPending

	timeNow := time.Now()

	docData.OwnerBy = authUsername
	docData.OwnerAt = timeNow

	docData.CreatedBy = authUsername
	docData.CreatedAt = timeNow

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc TaskHttpService) UpdateTask(shopID string, guid string, authUsername string, doc models.Task) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	tempStatus := findDoc.Status
	findDoc.Task = doc
	findDoc.Status = tempStatus

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc TaskHttpService) UpdateTaskStatus(shopID string, guid string, authUsername string, jobStatus int8) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	if findDoc.Status > jobStatus {
		return errors.New("task status invalid")
	}

	if jobStatus < models.TaskPending || jobStatus > models.TaskGlCompleted {
		return errors.New("task status out of range")
	}

	totalImageGroup := 0
	totalRejectImageGroup := 0
	if jobStatus == models.TaskCompleted {
		findDocImageGroups, err := svc.repoDocImageGroup.FindByTaskGUID(shopID, guid)

		if err != nil {
			return err
		}

		for _, docImageGroup := range findDocImageGroups {
			if docImageGroup.Status == modelDocumentImage.IMAGE_REJECT {
				totalRejectImageGroup += 1
			}

			totalImageGroup += 1
		}

		err = svc.repoDocImageGroup.UpdateTaskIsCompletedByTaskGUID(shopID, findDoc.GuidFixed, true)
		if err != nil {
			return err
		}
	}

	totalImageGroup = 0
	totalRejectImageGroup = 0
	if jobStatus == models.TaskGlCompleted {
		findDocImageGroups, err := svc.repoDocImageGroup.FindByTaskGUID(shopID, guid)

		if err != nil {
			return err
		}

		for _, docImageGroup := range findDocImageGroups {
			if docImageGroup.Status == modelDocumentImage.IMAGE_REJECT_KEYING {
				totalRejectImageGroup += 1
			}

			totalImageGroup += 1
		}

		// err = svc.repoDocImageGroup.UpdateTaskIsCompletedByTaskGUID(shopID, findDoc.GuidFixed, true)
		// if err != nil {
		// 	return err
		// }
	}

	if totalRejectImageGroup > 0 {
		newGuidFixed := utils.NewGUID()

		docData := models.TaskDoc{}
		docData.ShopID = shopID
		docData.GuidFixed = newGuidFixed
		docData.Module = findDoc.Module

		parentGUID := findDoc.GuidFixed

		if len(findDoc.ParentGUIDFixed) > 0 {
			parentGUID = findDoc.ParentGUIDFixed
			docData.Path = findDoc.Path
		} else {
			docData.Path = fmt.Sprintf("%s/%s", findDoc.Path, findDoc.GuidFixed)
		}

		docData.ParentGUIDFixed = parentGUID

		docData.Status = models.TaskPending

		timeNow := time.Now()

		docData.OwnerBy = findDoc.OwnerBy
		docData.OwnerAt = timeNow

		docData.RejectedBy = authUsername
		docData.RejectedAt = timeNow

		docData.CreatedBy = authUsername
		docData.CreatedAt = timeNow

		_, err = svc.repo.Create(docData)

		if err != nil {
			return err
		}
	}

	findDoc.Status = jobStatus
	findDoc.ToTal = totalImageGroup
	findDoc.ToTalReject = totalRejectImageGroup

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc TaskHttpService) DeleteTask(shopID string, guid string, authUsername string) error {

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

func (svc TaskHttpService) DeleteTaskByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc TaskHttpService) InfoTask(shopID string, guid string) (models.TaskInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.TaskInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.TaskInfo{}, errors.New("document not found")
	}

	return findDoc.TaskInfo, nil

}

func (svc TaskHttpService) GetTaskReject(shopID string, module string, taskGUID string) ([]models.TaskInfo, error) {

	docList, err := svc.repo.FindPageByTaskReject(shopID, module, taskGUID)

	if err != nil {
		return []models.TaskInfo{}, err
	}

	return docList, nil
}

func (svc TaskHttpService) SearchTask(shopID string, module string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.TaskInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"name",
	}

	docList, pagination, err := svc.repo.FindPageTask(shopID, module, filters, searchInFields, pageable)

	if err != nil {
		return []models.TaskInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc TaskHttpService) SearchTaskStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.TaskInfo, int, error) {
	searchInFields := []string{
		"name",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.TaskInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc TaskHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Task) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Task](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Task, models.TaskDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Task) models.TaskDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.TaskDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Task = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Task, models.TaskDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.TaskDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "name", guid)
		},
		func(doc models.TaskDoc) bool {
			return doc.Name != ""
		},
		func(shopID string, authUsername string, data models.Task, doc models.TaskDoc) error {

			doc.Task = data
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

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc TaskHttpService) getDocIDKey(doc models.Task) string {
	return doc.Name
}
