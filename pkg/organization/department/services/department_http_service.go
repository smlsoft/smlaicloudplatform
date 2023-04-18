package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/organization/department/models"
	"smlcloudplatform/pkg/organization/department/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IDepartmentHttpService interface {
	CreateDepartment(shopID string, authUsername string, doc models.Department) (string, error)
	UpdateDepartment(shopID string, guid string, authUsername string, doc models.Department) error
	DeleteDepartment(shopID string, guid string, authUsername string) error
	DeleteDepartmentByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoDepartment(shopID string, guid string) (models.DepartmentInfo, error)
	InfoDepartmentByCode(shopID string, code string) (models.DepartmentInfo, error)
	SearchDepartment(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DepartmentInfo, mongopagination.PaginationData, error)
	SearchDepartmentStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.DepartmentInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Department) (common.BulkImport, error)

	GetModuleName() string
}

type DepartmentHttpService struct {
	repo repositories.IDepartmentRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.DepartmentActivity, models.DepartmentDeleteActivity]
}

func NewDepartmentHttpService(repo repositories.IDepartmentRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *DepartmentHttpService {

	insSvc := &DepartmentHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.DepartmentActivity, models.DepartmentDeleteActivity](repo)

	return insSvc
}

func (svc DepartmentHttpService) CreateDepartment(shopID string, authUsername string, doc models.Department) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.DepartmentDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Department = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc DepartmentHttpService) UpdateDepartment(shopID string, guid string, authUsername string, doc models.Department) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.Department = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc DepartmentHttpService) DeleteDepartment(shopID string, guid string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc DepartmentHttpService) DeleteDepartmentByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc DepartmentHttpService) InfoDepartment(shopID string, guid string) (models.DepartmentInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.DepartmentInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.DepartmentInfo{}, errors.New("document not found")
	}

	return findDoc.DepartmentInfo, nil
}

func (svc DepartmentHttpService) InfoDepartmentByCode(shopID string, code string) (models.DepartmentInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)

	if err != nil {
		return models.DepartmentInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.DepartmentInfo{}, errors.New("document not found")
	}

	return findDoc.DepartmentInfo, nil
}

func (svc DepartmentHttpService) SearchDepartment(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DepartmentInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.DepartmentInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc DepartmentHttpService) SearchDepartmentStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.DepartmentInfo, int, error) {
	searchInFields := []string{
		"code",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.DepartmentInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc DepartmentHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Department) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Department](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Code)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "code", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Department, models.DepartmentDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Department) models.DepartmentDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.DepartmentDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Department = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Department, models.DepartmentDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.DepartmentDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.DepartmentDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Department, doc models.DepartmentDoc) error {

			doc.Department = data
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
		createDataKey = append(createDataKey, doc.Code)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Code)
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

func (svc DepartmentHttpService) getDocIDKey(doc models.Department) string {
	return doc.Code
}

func (svc DepartmentHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc DepartmentHttpService) GetModuleName() string {
	return "department"
}
