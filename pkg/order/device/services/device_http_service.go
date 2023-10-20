package services

import (
	"context"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/order/device/models"
	"smlcloudplatform/pkg/order/device/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IDeviceHttpService interface {
	CreateDevice(shopID string, authUsername string, doc models.Device) (string, error)
	UpdateDevice(shopID string, guid string, authUsername string, doc models.Device) error
	DeleteDevice(shopID string, guid string, authUsername string) error
	DeleteDeviceByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoDevice(shopID string, guid string) (models.DeviceInfo, error)
	InfoDeviceByCode(shopID string, code string) (models.DeviceInfo, error)
	SearchDevice(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DeviceInfo, mongopagination.PaginationData, error)
	SearchDeviceStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DeviceInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Device) (common.BulkImport, error)

	GetModuleName() string
}

type DeviceHttpService struct {
	repo repositories.IDeviceRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.DeviceActivity, models.DeviceDeleteActivity]
	contextTimeout time.Duration
}

func NewDeviceHttpService(
	repo repositories.IDeviceRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,

	contextTimeout time.Duration,
) *DeviceHttpService {

	insSvc := &DeviceHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,

		contextTimeout: contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.DeviceActivity, models.DeviceDeleteActivity](repo)

	return insSvc
}

func (svc DeviceHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc DeviceHttpService) CreateDevice(shopID string, authUsername string, doc models.Device) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "id", doc.ID)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("ID is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.DeviceDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Device = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", err
	}

	go func() {
		svc.saveMasterSync(shopID)
	}()

	return newGuidFixed, nil
}

func (svc DeviceHttpService) UpdateDevice(shopID string, guid string, authUsername string, doc models.Device) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.Device = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

	if err != nil {
		return err
	}

	go func() {
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc DeviceHttpService) DeleteDevice(shopID string, guid string, authUsername string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(ctx, shopID, guid, authUsername)
	if err != nil {
		return err
	}

	go func() {
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc DeviceHttpService) DeleteDeviceByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	go func() {
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc DeviceHttpService) InfoDevice(shopID string, guid string) (models.DeviceInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.DeviceInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.DeviceInfo{}, errors.New("document not found")
	}

	return findDoc.DeviceInfo, nil
}

func (svc DeviceHttpService) InfoDeviceByCode(shopID string, code string) (models.DeviceInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "id", code)

	if err != nil {
		return models.DeviceInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.DeviceInfo{}, errors.New("document not found")
	}

	return findDoc.DeviceInfo, nil
}

func (svc DeviceHttpService) SearchDevice(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DeviceInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"id",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.DeviceInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc DeviceHttpService) SearchDeviceStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DeviceInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"id",
	}

	selectFields := map[string]interface{}{}

	/*
		if langCode != "" {
			selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
		} else {
			selectFields["names"] = 1
		}
	*/

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.DeviceInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc DeviceHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Device) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Device](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.ID)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(ctx, shopID, "id", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.ID)
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
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "id", guid)
		},
		func(doc models.DeviceDoc) bool {
			return doc.ID != ""
		},
		func(shopID string, authUsername string, data models.Device, doc models.DeviceDoc) error {

			doc.Device = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			err = svc.repo.Update(ctx, shopID, doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(ctx, createDataList)

		if err != nil {
			return common.BulkImport{}, err
		}

	}

	createDataKey := []string{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.ID)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.ID)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.ID)
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
	return doc.ID
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
