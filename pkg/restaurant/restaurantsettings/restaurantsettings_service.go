package restaurantsettings

import (
	"errors"
	"fmt"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/restaurantsettings/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRestaurantSettingsService interface {
	CreateRestaurantSettings(shopID string, authUsername string, doc models.RestaurantSettings) (string, error)
	UpdateRestaurantSettings(shopID string, guid string, authUsername string, doc models.RestaurantSettings) error
	DeleteByGUIDs(shopID string, authUsername string, GUIDs []string) error
	DeleteRestaurantSettings(shopID string, guid string, authUsername string) error
	InfoRestaurantSettings(shopID string, guid string) (models.RestaurantSettingsInfo, error)
	SearchRestaurantSettings(shopID string, q string, page int, limit int) ([]models.RestaurantSettingsInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.RestaurantSettings) (common.BulkImport, error)
	ListRestaurantSettingsByCode(shopID string, code string, pagable common.Pageable) ([]models.RestaurantSettingsInfo, mongopagination.PaginationData, error)

	LastActivity(shopID string, action string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, mongopagination.PaginationData, error)
	GetModuleName() string
}

type RestaurantSettingsService struct {
	repo          IRestaurantSettingsRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository

	services.ActivityService[models.RestaurantSettingsActivity, models.RestaurantSettingsDeleteActivity]
}

func NewRestaurantSettingsService(repo IRestaurantSettingsRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) RestaurantSettingsService {
	insSvc := RestaurantSettingsService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.RestaurantSettingsActivity, models.RestaurantSettingsDeleteActivity](repo)
	return insSvc
}

func (svc RestaurantSettingsService) CreateRestaurantSettings(shopID string, authUsername string, doc models.RestaurantSettings) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.RestaurantSettingsDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.RestaurantSettings = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	docData.LastUpdatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc RestaurantSettingsService) UpdateRestaurantSettings(shopID string, guid string, authUsername string, doc models.RestaurantSettings) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.RestaurantSettings = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	findDoc.LastUpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc RestaurantSettingsService) DeleteRestaurantSettings(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc RestaurantSettingsService) DeleteByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc RestaurantSettingsService) InfoRestaurantSettings(shopID string, guid string) (models.RestaurantSettingsInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.RestaurantSettingsInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.RestaurantSettingsInfo{}, errors.New("document not found")
	}

	return findDoc.RestaurantSettingsInfo, nil
}

func (svc RestaurantSettingsService) ListRestaurantSettingsByCode(shopID string, code string, pagable common.Pageable) ([]models.RestaurantSettingsInfo, mongopagination.PaginationData, error) {

	docList, pagination, err := svc.repo.FindPageFilterSort(shopID, map[string]interface{}{"code": code}, []string{}, pagable.Q, pagable.Page, pagable.Limit, pagable.Sorts)

	if err != nil {
		return []models.RestaurantSettingsInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil

}

func (svc RestaurantSettingsService) SearchRestaurantSettings(shopID string, q string, page int, limit int) ([]models.RestaurantSettingsInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"body",
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []models.RestaurantSettingsInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc RestaurantSettingsService) SaveInBatch(shopID string, authUsername string, dataList []models.RestaurantSettings) (common.BulkImport, error) {

	// createDataList := []models.RestaurantSettingsDoc{}
	// duplicateDataList := []models.RestaurantSettings{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[models.RestaurantSettings](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadCategoryList {
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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.RestaurantSettings, models.RestaurantSettingsDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.RestaurantSettings) models.RestaurantSettingsDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.RestaurantSettingsDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.RestaurantSettings = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.RestaurantSettings, models.RestaurantSettingsDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.RestaurantSettingsDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc models.RestaurantSettingsDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data models.RestaurantSettings, doc models.RestaurantSettingsDoc) error {

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
	for _, doc := range payloadDuplicateCategoryList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.Code)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, doc.Code)
	}

	svc.saveMasterSync(shopID)

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc RestaurantSettingsService) getDocIDKey(doc models.RestaurantSettings) string {
	return doc.Code
}

func (svc RestaurantSettingsService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc RestaurantSettingsService) GetModuleName() string {
	return "restaurantsettings"
}
