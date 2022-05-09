package shopzone

import (
	"errors"
	"fmt"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"sync"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopZoneService interface {
	CreateShopZone(shopID string, authUsername string, doc restaurant.ShopZone) (string, error)
	UpdateShopZone(guid string, shopID string, authUsername string, doc restaurant.ShopZone) error
	DeleteShopZone(guid string, shopID string, authUsername string) error
	InfoShopZone(guid string, shopID string) (restaurant.ShopZoneInfo, error)
	SearchShopZone(shopID string, q string, page int, limit int) ([]restaurant.ShopZoneInfo, mongopagination.PaginationData, error)
	LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []restaurant.ShopZone) (models.BulkImport, error)
}

type ShopZoneService struct {
	repo IShopZoneRepository
}

func NewShopZoneService(
	repo IShopZoneRepository,
) ShopZoneService {

	return ShopZoneService{
		repo: repo,
	}
}

func (svc ShopZoneService) CreateShopZone(shopID string, authUsername string, doc restaurant.ShopZone) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := restaurant.ShopZoneDoc{}
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
	return newGuidFixed, nil
}

func (svc ShopZoneService) UpdateShopZone(guid string, shopID string, authUsername string, doc restaurant.ShopZone) error {

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

	err = svc.repo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopZoneService) DeleteShopZone(guid string, shopID string, authUsername string) error {
	err := svc.repo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopZoneService) InfoShopZone(guid string, shopID string) (restaurant.ShopZoneInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return restaurant.ShopZoneInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return restaurant.ShopZoneInfo{}, errors.New("document not found")
	}

	return findDoc.ShopZoneInfo, nil

}

func (svc ShopZoneService) SearchShopZone(shopID string, q string, page int, limit int) ([]restaurant.ShopZoneInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"code",
	}

	for i := range [5]bool{} {
		searchCols = append(searchCols, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []restaurant.ShopZoneInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ShopZoneService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, mongopagination.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []restaurant.ShopZoneDeleteActivity
	var pagination1 mongopagination.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.repo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []restaurant.ShopZoneActivity
	var pagination2 mongopagination.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.repo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return models.LastActivity{}, pagination1, err1
	}

	if err2 != nil {
		return models.LastActivity{}, pagination2, err2
	}

	lastActivity := models.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	pagination := pagination1

	if pagination.Total < pagination2.Total {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc ShopZoneService) SaveInBatch(shopID string, authUsername string, dataList []restaurant.ShopZone) (models.BulkImport, error) {

	createDataList := []restaurant.ShopZoneDoc{}
	duplicateDataList := []restaurant.ShopZone{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[restaurant.ShopZone](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadCategoryList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Code)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "code", itemCodeGuidList)

	if err != nil {
		return models.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[restaurant.ShopZone, restaurant.ShopZoneDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc restaurant.ShopZone) restaurant.ShopZoneDoc {
			newGuid := utils.NewGUID()

			dataDoc := restaurant.ShopZoneDoc{}

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

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[restaurant.ShopZone, restaurant.ShopZoneDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (restaurant.ShopZoneDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc restaurant.ShopZoneDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data restaurant.ShopZone, doc restaurant.ShopZoneDoc) error {

			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(createDataList)

		if err != nil {
			return models.BulkImport{}, err
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

	return models.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc ShopZoneService) getDocIDKey(doc restaurant.ShopZone) string {
	return doc.Code
}
