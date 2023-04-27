package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/organization/branch/models"
	"smlcloudplatform/pkg/organization/branch/repositories"
	businessTypeRepositories "smlcloudplatform/pkg/organization/businesstype/repositories"
	deparmentRepositories "smlcloudplatform/pkg/organization/department/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IBranchHttpService interface {
	CreateBranch(shopID string, authUsername string, doc models.Branch) (string, error)
	UpdateBranch(shopID string, guid string, authUsername string, doc models.Branch) error
	DeleteBranch(shopID string, guid string, authUsername string) error
	DeleteBranchByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoBranch(shopID string, guid string) (models.BranchInfo, error)
	InfoBranchByCode(shopID string, code string) (models.BranchInfo, error)
	SearchBranch(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.BranchInfo, mongopagination.PaginationData, error)
	SearchBranchStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.BranchInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Branch) (common.BulkImport, error)

	GetModuleName() string
}

type BranchHttpService struct {
	repo             repositories.IBranchRepository
	repoDepartment   deparmentRepositories.IDepartmentRepository
	repoBusinessType businessTypeRepositories.IBusinessTypeRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.BranchActivity, models.BranchDeleteActivity]
}

func NewBranchHttpService(repo repositories.IBranchRepository, repoDepartment deparmentRepositories.IDepartmentRepository, repoBusinessType businessTypeRepositories.IBusinessTypeRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *BranchHttpService {

	insSvc := &BranchHttpService{
		repo:             repo,
		repoDepartment:   repoDepartment,
		repoBusinessType: repoBusinessType,
		syncCacheRepo:    syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.BranchActivity, models.BranchDeleteActivity](repo)

	return insSvc
}

func (svc BranchHttpService) CreateBranch(shopID string, authUsername string, doc models.Branch) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.BranchDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Branch = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc BranchHttpService) UpdateBranch(shopID string, guid string, authUsername string, doc models.Branch) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.Branch = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc BranchHttpService) DeleteBranch(shopID string, guid string, authUsername string) error {

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

func (svc BranchHttpService) DeleteBranchByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc BranchHttpService) InfoBranch(shopID string, guid string) (models.BranchInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.BranchInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.BranchInfo{}, errors.New("document not found")
	}

	return findDoc.BranchInfo, nil
}

func (svc BranchHttpService) InfoBranchByCode(shopID string, code string) (models.BranchInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)

	if err != nil {
		return models.BranchInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.BranchInfo{}, errors.New("document not found")
	}

	return findDoc.BranchInfo, nil
}

func (svc BranchHttpService) SearchBranch(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.BranchInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.BranchInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc BranchHttpService) SearchBranchStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.BranchInfo, int, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.BranchInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc BranchHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Branch) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Branch](dataList, svc.getDocIDKey)

	itemCodeGuidList := []interface{}{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Code)
	}

	findItemGuid, err := svc.repo.FindInItemGuids(shopID, "code", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Branch, models.BranchDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Branch) models.BranchDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.BranchDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Branch = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Branch, models.BranchDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.BranchDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.BranchDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Branch, doc models.BranchDoc) error {

			doc.Branch = data
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

	createDataKey := []interface{}{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.Code)
	}

	payloadDuplicateDataKey := []interface{}{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []interface{}{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Code)
	}

	updateFailDataKey := []interface{}{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, svc.getDocIDKey(doc))
	}

	svc.saveMasterSync(shopID)

	tempCreateDataKey := svc.toSliceString(createDataKey)
	tempUpdateDataKey := svc.toSliceString(updateDataKey)
	tempUpdateFailDataKey := svc.toSliceString(updateFailDataKey)
	tempPayloadDuplicateDataKey := svc.toSliceString(payloadDuplicateDataKey)

	return common.BulkImport{
		Created:          tempCreateDataKey,
		Updated:          tempUpdateDataKey,
		UpdateFailed:     tempUpdateFailDataKey,
		PayloadDuplicate: tempPayloadDuplicateDataKey,
	}, nil
}

func (svc BranchHttpService) toSliceString(data []interface{}) []string {
	tempData := make([]string, len(data))
	for i, v := range data {
		tempData[i] = v.(string)
	}

	return tempData
}

func (svc BranchHttpService) getDocIDKey(doc models.Branch) string {
	return doc.Code
}

func (svc BranchHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc BranchHttpService) GetModuleName() string {
	return "branch"
}
