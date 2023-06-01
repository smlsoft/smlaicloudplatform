package shopzone

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/shopzone/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopZoneService interface {
	CreateShopZone(shopID string, authUsername string, doc models.ShopZone) (string, error)
	UpdateShopZone(shopID string, guid string, authUsername string, doc models.ShopZone) error
	DeleteShopZone(shopID string, guid string, authUsername string) error
	InfoShopZone(shopID string, guid string) (models.ShopZoneInfo, error)
	SearchShopZone(shopID string, pageable micromodels.Pageable) ([]models.ShopZoneInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ShopZone) (common.BulkImport, error)

	GetModuleName() string
}

type ShopZoneService struct {
	repo          IShopZoneRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository

	services.ActivityService[models.ShopZoneActivity, models.ShopZoneDeleteActivity]
}

func NewShopZoneService(repo IShopZoneRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) ShopZoneService {
	insSvc := ShopZoneService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.ShopZoneActivity, models.ShopZoneDeleteActivity](repo)
	return insSvc
}

func (svc ShopZoneService) CreateShopZone(shopID string, authUsername string, doc models.ShopZone) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.ShopZoneDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ShopZone = doc

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

func (svc ShopZoneService) UpdateShopZone(shopID string, guid string, authUsername string, doc models.ShopZone) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ShopZone = doc

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

func (svc ShopZoneService) DeleteShopZone(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ShopZoneService) InfoShopZone(shopID string, guid string) (models.ShopZoneInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ShopZoneInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ShopZoneInfo{}, errors.New("document not found")
	}

	return findDoc.ShopZoneInfo, nil

}

func (svc ShopZoneService) SearchShopZone(shopID string, pageable micromodels.Pageable) ([]models.ShopZoneInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"name1",
		"name2",
		"name3",
		"name4",
		"name5",
	}

	for i := range [5]bool{} {
		searchInFields = append(searchInFields, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchInFields, pageable)

	if err != nil {
		return []models.ShopZoneInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ShopZoneService) SaveInBatch(shopID string, authUsername string, dataList []models.ShopZone) (common.BulkImport, error) {

	// createDataList := []models.ShopZoneDoc{}
	// duplicateDataList := []models.ShopZone{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[models.ShopZone](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.ShopZone, models.ShopZoneDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.ShopZone) models.ShopZoneDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ShopZoneDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.ShopZone = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.ShopZone, models.ShopZoneDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ShopZoneDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc models.ShopZoneDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data models.ShopZone, doc models.ShopZoneDoc) error {

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

func (svc ShopZoneService) getDocIDKey(doc models.ShopZone) string {
	return doc.Code
}

func (svc ShopZoneService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc ShopZoneService) GetModuleName() string {
	return "shopzone"
}
