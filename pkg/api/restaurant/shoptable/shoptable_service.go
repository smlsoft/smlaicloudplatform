package shoptable

import (
	"errors"
	"fmt"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"sync"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopTableService interface {
	CreateShopTable(shopID string, authUsername string, doc restaurant.ShopTable) (string, error)
	UpdateShopTable(guid string, shopID string, authUsername string, doc restaurant.ShopTable) error
	DeleteShopTable(guid string, shopID string, authUsername string) error
	InfoShopTable(guid string, shopID string) (restaurant.ShopTableInfo, error)
	SearchShopTable(shopID string, q string, page int, limit int) ([]restaurant.ShopTableInfo, mongopagination.PaginationData, error)
	LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []restaurant.ShopTable) (models.BulkImport, error)
}

type ShopTableService struct {
	crudRepo     repositories.CrudRepository[restaurant.ShopTableDoc]
	searchRepo   repositories.SearchRepository[restaurant.ShopTableInfo]
	guidRepo     repositories.GuidRepository[restaurant.ShopTableItemGuid]
	activityRepo repositories.ActivityRepository[restaurant.ShopTableActivity, restaurant.ShopTableDeleteActivity]
}

func NewShopTableService(
	crudRepo repositories.CrudRepository[restaurant.ShopTableDoc],
	searchRepo repositories.SearchRepository[restaurant.ShopTableInfo],
	guidRepo repositories.GuidRepository[restaurant.ShopTableItemGuid],
	activityRepo repositories.ActivityRepository[restaurant.ShopTableActivity, restaurant.ShopTableDeleteActivity],
) ShopTableService {
	return ShopTableService{
		crudRepo:     crudRepo,
		searchRepo:   searchRepo,
		guidRepo:     guidRepo,
		activityRepo: activityRepo,
	}
}

func (svc ShopTableService) CreateShopTable(shopID string, authUsername string, doc restaurant.ShopTable) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := restaurant.ShopTableDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ShopTable = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.crudRepo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc ShopTableService) UpdateShopTable(guid string, shopID string, authUsername string, doc restaurant.ShopTable) error {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ShopTable = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.crudRepo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopTableService) DeleteShopTable(guid string, shopID string, authUsername string) error {
	err := svc.crudRepo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopTableService) InfoShopTable(guid string, shopID string) (restaurant.ShopTableInfo, error) {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return restaurant.ShopTableInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return restaurant.ShopTableInfo{}, errors.New("document not found")
	}

	return findDoc.ShopTableInfo, nil

}

func (svc ShopTableService) SearchShopTable(shopID string, q string, page int, limit int) ([]restaurant.ShopTableInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"code",
	}

	for i := range [5]bool{} {
		searchCols = append(searchCols, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.searchRepo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []restaurant.ShopTableInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ShopTableService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, mongopagination.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []restaurant.ShopTableDeleteActivity
	var pagination1 mongopagination.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.activityRepo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []restaurant.ShopTableActivity
	var pagination2 mongopagination.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.activityRepo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
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

func (svc ShopTableService) SaveInBatch(shopID string, authUsername string, dataList []restaurant.ShopTable) (models.BulkImport, error) {

	createDataList := []restaurant.ShopTableDoc{}
	duplicateDataList := []restaurant.ShopTable{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[restaurant.ShopTable](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadCategoryList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Number)
	}

	findItemGuid, err := svc.guidRepo.FindInItemGuid(shopID, "code", itemCodeGuidList)

	if err != nil {
		return models.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[restaurant.ShopTable, restaurant.ShopTableDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc restaurant.ShopTable) restaurant.ShopTableDoc {
			newGuid := utils.NewGUID()

			dataDoc := restaurant.ShopTableDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.ShopTable = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[restaurant.ShopTable, restaurant.ShopTableDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (restaurant.ShopTableDoc, error) {
			return svc.crudRepo.FindByGuid(shopID, guid)
		},
		func(doc restaurant.ShopTableDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data restaurant.ShopTable, doc restaurant.ShopTableDoc) error {

			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.crudRepo.CreateInBatch(createDataList)

		if err != nil {
			return models.BulkImport{}, err
		}
	}
	createDataKey := []string{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.Number)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateCategoryList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Number)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.Number)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, doc.Number)
	}

	return models.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc ShopTableService) getDocIDKey(doc restaurant.ShopTable) string {
	return doc.Number
}
