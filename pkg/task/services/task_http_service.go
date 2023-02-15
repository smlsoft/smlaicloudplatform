package services

import (
	"crypto/md5"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	documentImageModel "smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	repositoriesDocumentImage "smlcloudplatform/pkg/documentwarehouse/documentimage/repositories"
	servicesDocumentImage "smlcloudplatform/pkg/documentwarehouse/documentimage/services"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/task/models"
	"smlcloudplatform/pkg/task/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"strconv"
	"strings"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITaskHttpService interface {
	GenerateTaskID(shopID string, authUsername string) (string, error)
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

	UpdateTaskTotalImage(docReq documentImageModel.DocumentImageTaskChangeMessage) error
	UpdateTaskTotalRejectImage(docReq documentImageModel.DocumentImageTaskRejectMessage) error
}

type TaskHttpService struct {
	repo              repositories.ITaskRepository
	repoDocImageGroup repositoriesDocumentImage.IDocumentImageGroupRepository
	serviceDocImage   servicesDocumentImage.IDocumentImageService
	services.ActivityService[models.TaskActivity, models.TaskDeleteActivity]
}

func NewTaskHttpService(repo repositories.ITaskRepository, repoDocImageGroup repositoriesDocumentImage.IDocumentImageGroupRepository, serviceDocImage servicesDocumentImage.IDocumentImageService) *TaskHttpService {

	insSvc := &TaskHttpService{
		repo:              repo,
		repoDocImageGroup: repoDocImageGroup,
		serviceDocImage:   serviceDocImage,
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

func (svc TaskHttpService) GenerateTaskID(shopID string, authUsername string) (string, error) {

	timeNow := time.Now()

	userHash := fmt.Sprintf("%x", md5.Sum([]byte(authUsername)))
	userHashLength := len(userHash)
	userFmt := userHash[(userHashLength - 8):userHashLength]
	timeFmt := timeNow.Format("20060102")

	codeFmt := fmt.Sprintf("task%s-%s-", userFmt, timeFmt)

	findDoc, err := svc.repo.FindLastTaskByCode(shopID, codeFmt)
	if err != nil {
		return "", err
	}

	if len(findDoc.Code) == 0 {
		return codeFmt + "001", nil
	}

	lastTaskCode := findDoc.Code

	splitStrTaskCode := strings.Split(lastTaskCode, "-")

	lastNumberStr := splitStrTaskCode[len(splitStrTaskCode)-1]

	lastNumber, err := strconv.Atoi(lastNumberStr)

	if err != nil {
		return "", err
	}

	lastNumber++

	return codeFmt + fmt.Sprintf("%03d", lastNumber), nil
}

func (svc TaskHttpService) PaddingNumber(number int) string {
	return fmt.Sprintf("%03d", number)
}

func (svc TaskHttpService) CreateTask(shopID string, authUsername string, doc models.Task) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "name", doc.Name)

	if err != nil {
		return "", err
	}

	if findDoc.Code == doc.Code {
		return "", errors.New("code is duplicate")
	}

	if findDoc.Name != "" {
		return "", errors.New("name is empty")
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

	updateDoc := findDoc
	updateDoc.Task = doc
	updateDoc.Status = findDoc.Status
	updateDoc.Code = findDoc.Code

	updateDoc.UpdatedBy = authUsername
	updateDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, updateDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc TaskHttpService) UpdateTaskStatus(shopID string, taskGUID string, authUsername string, jobStatus int8) error {

	findDoc, err := svc.repo.FindByGuid(shopID, taskGUID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	if findDoc.Status >= jobStatus {
		return errors.New("task status invalid")
	}

	if jobStatus < models.TaskPending || jobStatus > models.TaskGlCompleted {
		return errors.New("task status out of range")
	}

	totalImageGroup := 0
	totalRejectImageGroup := 0
	if jobStatus == models.TaskCompleted {
		findDocImageGroups, err := svc.repoDocImageGroup.FindByTaskGUID(shopID, taskGUID)

		if err != nil {
			return err
		}

		for _, docImageGroup := range findDocImageGroups {
			if docImageGroup.Status == documentImageModel.IMAGE_REJECT {
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
	totalRejectKeyingImageGroup := 0
	if jobStatus == models.TaskGlCompleted {
		findDocImageGroups, err := svc.repoDocImageGroup.FindByTaskGUID(shopID, taskGUID)

		if err != nil {
			return err
		}

		for _, docImageGroup := range findDocImageGroups {
			if docImageGroup.Status == documentImageModel.IMAGE_REJECT_KEYING {
				totalRejectKeyingImageGroup += 1
			}

			totalImageGroup += 1
		}

		// err = svc.repoDocImageGroup.UpdateTaskIsCompletedByTaskGUID(shopID, findDoc.GuidFixed, true)
		// if err != nil {
		// 	return err
		// }
	}

	if totalRejectImageGroup > 0 || totalRejectKeyingImageGroup > 0 {
		newTaskGuidFixed := utils.NewGUID()

		docData := models.TaskDoc{}
		docData.ShopID = shopID
		docData.GuidFixed = newTaskGuidFixed
		docData.Module = findDoc.Module

		parentGUID := findDoc.GuidFixed
		if len(findDoc.ParentGUIDFixed) > 0 {
			parentGUID = findDoc.ParentGUIDFixed
		}

		if len(findDoc.ParentGUIDFixed) > 0 {
			parentGUID = findDoc.ParentGUIDFixed
			docData.Path = findDoc.Path
		} else {
			docData.Path = fmt.Sprintf("%s/%s", findDoc.Path, findDoc.GuidFixed)
		}

		taskCount, err := svc.repo.CountTaskParent(shopID, parentGUID)

		if err != nil {
			taskCount = 0
		}

		taskCount += 1
		rejectTaskName := fmt.Sprintf("%s - [%d]", findDoc.Task.Name, taskCount)

		docData.Name = rejectTaskName
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

		if jobStatus == models.TaskCompleted {
			findDocImageGroups, err := svc.repoDocImageGroup.FindByTaskGUID(shopID, taskGUID)

			if err != nil {
				return err
			}

			for _, docImageGroup := range findDocImageGroups {

				if docImageGroup.Status == documentImageModel.IMAGE_REJECT {
					newDocImgGroupGUID := utils.NewGUID()
					for _, docImage := range *docImageGroup.ImageReferences {
						docImgReq := documentImageModel.DocumentImageRequest{}

						docImgReq.DocumentImageGroupGUID = newDocImgGroupGUID
						docImgReq.TaskGUID = newTaskGuidFixed
						docImgReq.ImageURI = docImage.ImageURI
						docImgReq.Name = docImage.Name + " -- REJECT"

						gid, id, err := svc.serviceDocImage.CreateDocumentImage(shopID, authUsername, docImgReq)

						fmt.Println(gid, id)

						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
	}

	findDoc.Status = jobStatus
	// findDoc.ToTal = totalImageGroup
	// findDoc.ToTalReject = totalRejectImageGroup

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, taskGUID, findDoc)

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

func (svc TaskHttpService) UpdateTaskTotalImage(docReq documentImageModel.DocumentImageTaskChangeMessage) error {

	// findDoc, err := svc.repo.FindByGuid(docReq.ShopID, docReq.TaskGUID)

	// if err != nil {
	// 	return err
	// }

	// if findDoc.ID == primitive.NilObjectID {
	// 	return errors.New("document not found")
	// }

	// if docReq.Event == documentImageModel.TaskChangePlus {
	// 	findDoc.ToTal = findDoc.ToTal + docReq.Count
	// } else if docReq.Event == documentImageModel.TaskChangeMinus && findDoc.ToTal > 0 {
	// 	findDoc.ToTal = findDoc.ToTal - docReq.Count
	// }

	err := svc.repo.UpdateTotalDocumentImageGroup(docReq.ShopID, docReq.TaskGUID, docReq.Count)

	if err != nil {
		return err
	}

	return nil
}

func (svc TaskHttpService) UpdateTaskTotalRejectImage(docReq documentImageModel.DocumentImageTaskRejectMessage) error {

	// findDoc, err := svc.repo.FindByGuid(docReq.ShopID, docReq.TaskGUID)

	// if err != nil {
	// 	return err
	// }

	// if findDoc.ID == primitive.NilObjectID {
	// 	return errors.New("document not found")
	// }

	// currentTotal := findDoc.ToTalReject
	// if docReq.Event == documentImageModel.TaskRejectPlus {
	// 	findDoc.ToTalReject = currentTotal + docReq.Count
	// } else if docReq.Event == documentImageModel.TaskRejectMinus && currentTotal > 0 {
	// 	findDoc.ToTalReject = currentTotal - docReq.Count
	// }

	err := svc.repo.UpdateTotalRejectDocumentImageGroup(docReq.ShopID, docReq.TaskGUID, docReq.Count)

	if err != nil {
		return err
	}

	return nil
}
