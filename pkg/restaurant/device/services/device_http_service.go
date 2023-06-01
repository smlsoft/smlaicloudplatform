package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/device/models"
	"smlcloudplatform/pkg/restaurant/device/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDeviceHttpService interface {
	CreateDevice(shopID string, authUsername string, doc models.Device) (string, error)
	UpdateDevice(shopID string, guid string, authUsername string, doc models.Device) error
	DeleteDevice(shopID string, guid string, authUsername string) error
	DeleteDeviceByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoDevice(shopID string, guid string) (models.DeviceInfo, error)
	SearchDevice(shopID string, pageable micromodels.Pageable) ([]models.DeviceInfo, mongopagination.PaginationData, error)
	SearchDeviceStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.DeviceInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Device) (common.BulkImport, error)

	GetModuleName() string
}

type DeviceHttpService struct {
	repo repositories.IDeviceRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.DeviceActivity, models.DeviceDeleteActivity]
}

func NewDeviceHttpService(repo repositories.IDeviceRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *DeviceHttpService {

	insSvc := &DeviceHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.DeviceActivity, models.DeviceDeleteActivity](repo)

	return insSvc
}

func (svc DeviceHttpService) CreateDevice(shopID string, authUsername string, doc models.Device) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.DeviceDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Device = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc DeviceHttpService) UpdateDevice(shopID string, guid string, authUsername string, doc models.Device) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Device = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc DeviceHttpService) DeleteDevice(shopID string, guid string, authUsername string) error {

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

func (svc DeviceHttpService) DeleteDeviceByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc DeviceHttpService) InfoDevice(shopID string, guid string) (models.DeviceInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.DeviceInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DeviceInfo{}, errors.New("document not found")
	}

	return findDoc.DeviceInfo, nil

}

func (svc DeviceHttpService) SearchDevice(shopID string, pageable micromodels.Pageable) ([]models.DeviceInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"name1",
		"name2",
		"name3",
		"name4",
		"name5",
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchInFields, pageable)

	if err != nil {
		return []models.DeviceInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc DeviceHttpService) SearchDeviceStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.DeviceInfo, int, error) {
	searchInFields := []string{
		"code",
		"name1",
		"name2",
		"name3",
		"name4",
		"name5",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.DeviceInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc DeviceHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Device) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Device](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Device, models.DeviceDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Device) models.DeviceDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.DeviceDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Device = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Device, models.DeviceDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.DeviceDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.DeviceDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Device, doc models.DeviceDoc) error {

			doc.Device = data
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

func (svc DeviceHttpService) getDocIDKey(doc models.Device) string {
	return doc.Code
}

func (svc DeviceHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc DeviceHttpService) GetModuleName() string {
	return "device"
}
