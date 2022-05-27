package kitchen

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

type IKitchenService interface {
	CreateKitchen(shopID string, authUsername string, doc restaurant.Kitchen) (string, error)
	UpdateKitchen(shopID string, guid string, authUsername string, doc restaurant.Kitchen) error
	DeleteKitchen(shopID string, guid string, authUsername string) error
	InfoKitchen(shopID string, guid string) (restaurant.KitchenInfo, error)
	SearchKitchen(shopID string, q string, page int, limit int) ([]restaurant.KitchenInfo, mongopagination.PaginationData, error)
	LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []restaurant.Kitchen) (models.BulkImport, error)
}

type KitchenService struct {
	repo KitchenRepository
}

func NewKitchenService(repo KitchenRepository) KitchenService {

	return KitchenService{
		repo: repo,
	}
}

func (svc KitchenService) CreateKitchen(shopID string, authUsername string, doc restaurant.Kitchen) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := restaurant.KitchenDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Kitchen = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc KitchenService) UpdateKitchen(shopID string, guid string, authUsername string, doc restaurant.Kitchen) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Kitchen = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc KitchenService) DeleteKitchen(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc KitchenService) InfoKitchen(shopID string, guid string) (restaurant.KitchenInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return restaurant.KitchenInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return restaurant.KitchenInfo{}, errors.New("document not found")
	}

	return findDoc.KitchenInfo, nil

}

func (svc KitchenService) SearchKitchen(shopID string, q string, page int, limit int) ([]restaurant.KitchenInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	for i := range [5]bool{} {
		searchCols = append(searchCols, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []restaurant.KitchenInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc KitchenService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, mongopagination.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []restaurant.KitchenDeleteActivity
	var pagination1 mongopagination.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.repo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []restaurant.KitchenActivity
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

func (svc KitchenService) SaveInBatch(shopID string, authUsername string, dataList []restaurant.Kitchen) (models.BulkImport, error) {

	createDataList := []restaurant.KitchenDoc{}
	duplicateDataList := []restaurant.Kitchen{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[restaurant.Kitchen](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList = importdata.PreparePayloadData[restaurant.Kitchen, restaurant.KitchenDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc restaurant.Kitchen) restaurant.KitchenDoc {
			newGuid := utils.NewGUID()

			dataDoc := restaurant.KitchenDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Kitchen = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[restaurant.Kitchen, restaurant.KitchenDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (restaurant.KitchenDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc restaurant.KitchenDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data restaurant.Kitchen, doc restaurant.KitchenDoc) error {

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

func (svc KitchenService) getDocIDKey(doc restaurant.Kitchen) string {
	return doc.Code
}
