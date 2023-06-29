package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/staff/models"
	"smlcloudplatform/pkg/restaurant/staff/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStaffHttpService interface {
	CreateStaff(shopID string, authUsername string, doc models.Staff) (string, error)
	UpdateStaff(shopID string, guid string, authUsername string, doc models.Staff) error
	DeleteStaff(shopID string, guid string, authUsername string) error
	DeleteStaffByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoStaff(shopID string, guid string) (models.StaffInfo, error)
	SearchStaff(shopID string, pageable micromodels.Pageable) ([]models.StaffInfo, mongopagination.PaginationData, error)
	SearchStaffStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.StaffInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Staff) (common.BulkImport, error)

	GetModuleName() string
}

type StaffHttpService struct {
	repo repositories.IStaffRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.StaffActivity, models.StaffDeleteActivity]
}

func NewStaffHttpService(repo repositories.IStaffRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *StaffHttpService {

	insSvc := &StaffHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.StaffActivity, models.StaffDeleteActivity](repo)

	return insSvc
}

func (svc StaffHttpService) CreateStaff(shopID string, authUsername string, doc models.Staff) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.StaffDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Staff = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc StaffHttpService) UpdateStaff(shopID string, guid string, authUsername string, doc models.Staff) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Staff = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc StaffHttpService) DeleteStaff(shopID string, guid string, authUsername string) error {

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

func (svc StaffHttpService) DeleteStaffByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc StaffHttpService) InfoStaff(shopID string, guid string) (models.StaffInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.StaffInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.StaffInfo{}, errors.New("document not found")
	}

	return findDoc.StaffInfo, nil

}

func (svc StaffHttpService) SearchStaff(shopID string, pageable micromodels.Pageable) ([]models.StaffInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchInFields, pageable)

	if err != nil {
		return []models.StaffInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc StaffHttpService) SearchStaffStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.StaffInfo, int, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.StaffInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc StaffHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Staff) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Staff](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Staff, models.StaffDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Staff) models.StaffDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.StaffDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Staff = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Staff, models.StaffDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.StaffDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.StaffDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Staff, doc models.StaffDoc) error {

			doc.Staff = data
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

func (svc StaffHttpService) getDocIDKey(doc models.Staff) string {
	return doc.Code
}

func (svc StaffHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc StaffHttpService) GetModuleName() string {
	return "restaurant-staff"
}
