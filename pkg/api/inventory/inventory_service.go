package inventory

import (
	"errors"
	"fmt"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"sync"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryService interface {
	IsExistsGuid(shopID string, guidFixed string) (bool, error)
	CreateInBatch(shopID string, authUsername string, inventories []models.Inventory) (models.InventoryBulkImport, error)
	CreateWithGuid(shopID string, authUsername string, guidFixed string, inventory models.Inventory) (string, error)
	CreateInventory(shopID string, authUsername string, inventory models.Inventory) (string, string, error)
	UpdateInventory(shopID string, guid string, authUsername string, inventory models.Inventory) error
	DeleteInventory(shopID string, guid string, username string) error
	InfoInventory(shopID string, guid string) (models.InventoryInfo, error)
	InfoMongoInventory(id string) (models.InventoryInfo, error)
	SearchInventory(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error)
	LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, paginate.PaginationData, error)
	UpdateProductCategory(shopID string, authUsername string, catId string, guid []string) error
}

type InventoryService struct {
	invRepo   IInventoryRepository
	invMqRepo IInventoryMQRepository
	cacheRepo mastersync.IMasterSyncCacheRepository
}

func NewInventoryService(inventoryRepo IInventoryRepository, inventoryMqRepo IInventoryMQRepository, cacheRepo mastersync.IMasterSyncCacheRepository) InventoryService {
	return InventoryService{
		invRepo:   inventoryRepo,
		invMqRepo: inventoryMqRepo,
		cacheRepo: cacheRepo,
	}
}

func (svc InventoryService) IsExistsGuid(shopID string, guidFixed string) (bool, error) {

	findDoc, err := svc.invRepo.FindByGuid(shopID, guidFixed)

	if err != nil {
		return false, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return false, nil
	}

	return true, nil
}

func (svc InventoryService) CreateInBatch(shopID string, authUsername string, inventories []models.Inventory) (models.InventoryBulkImport, error) {

	createDataList := []models.InventoryDoc{}
	duplicateDataList := []models.Inventory{}

	payloadInventoryList, payloadDuplicateInventoryList := filterDuplicateInventory(inventories)

	itemCodeGuidList := []string{}
	for _, inventory := range payloadInventoryList {
		itemCodeGuidList = append(itemCodeGuidList, inventory.ItemGuid)
	}

	findItemGuid, err := svc.invRepo.FindByItemCodeGuid(shopID, itemCodeGuidList)

	if err != nil {
		return models.InventoryBulkImport{}, err
	}

	duplicateDataList, createDataList = preparePayloadDataInventory(shopID, authUsername, findItemGuid, payloadInventoryList)

	updateSuccessDataList, updateFailDataList := updateOnDuplicateInventory(shopID, authUsername, duplicateDataList, svc.invRepo)

	if len(createDataList) > 0 {
		err = svc.invRepo.CreateInBatch(createDataList)

		if err != nil {
			return models.InventoryBulkImport{}, err
		}
	}
	createDataKey := []string{}

	for _, inv := range createDataList {
		createDataKey = append(createDataKey, inv.ItemGuid)

		// reply kafka
		if svc.invMqRepo != nil {
			err = svc.invMqRepo.Create(inv.InventoryData)

			if err != nil {
				return models.InventoryBulkImport{}, err
			}
		}
	}

	payloadDuplicateDataKey := []string{}
	for _, inv := range payloadDuplicateInventoryList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, inv.ItemGuid)
	}

	updateDataKey := []string{}
	for _, inv := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, inv.ItemGuid)
		// reply kafka
		if svc.invMqRepo != nil {
			err = svc.invMqRepo.Update(inv.InventoryData)

			if err != nil {
				return models.InventoryBulkImport{}, err
			}
		}
	}

	updateFailDataKey := []string{}
	for _, inv := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, inv.ItemGuid)
	}

	svc.saveMasterSync(shopID)

	return models.InventoryBulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func filterDuplicateInventory(inventories []models.Inventory) (itemTemp []models.Inventory, itemDuplicate []models.Inventory) {
	tempFilterDict := map[string]models.Inventory{}
	for _, inventory := range inventories {
		if _, ok := tempFilterDict[inventory.ItemGuid]; ok {
			itemDuplicate = append(itemDuplicate, inventory)

		}
		tempFilterDict[inventory.ItemGuid] = inventory
	}

	for _, inventory := range tempFilterDict {
		itemTemp = append(itemTemp, inventory)
	}

	return itemTemp, itemDuplicate
}

func updateOnDuplicateInventory(shopID string, authUsername string, duplicateDataList []models.Inventory, repo IInventoryRepository) ([]models.InventoryDoc, []models.Inventory) {
	updateSuccessDataList := []models.InventoryDoc{}
	updateFailDataList := []models.Inventory{}
	for _, inv := range duplicateDataList {
		findDoc, err := repo.FindByItemGuid(shopID, inv.ItemGuid)

		if err != nil || findDoc.ID == primitive.NilObjectID {
			updateFailDataList = append(updateFailDataList, inv)
			continue
		}

		findDoc.Inventory = inv

		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()
		findDoc.LastUpdatedAt = time.Now()

		err = repo.Update(shopID, findDoc.GuidFixed, findDoc)

		if err != nil {
			updateFailDataList = append(updateFailDataList, inv)
			continue
		}

		updateSuccessDataList = append(updateSuccessDataList, findDoc)
	}
	return updateSuccessDataList, updateFailDataList
}

func preparePayloadDataInventory(shopID string, authUsername string, findItemGuid []models.InventoryItemGuid, payloadInventoryList []models.Inventory) ([]models.Inventory, []models.InventoryDoc) {
	createDataList := []models.InventoryDoc{}
	duplicateDataList := []models.Inventory{}
	tempItemGuidDict := make(map[string]bool)

	for _, itemGuid := range findItemGuid {
		tempItemGuidDict[itemGuid.ItemGuid] = true
	}

	for _, inventory := range payloadInventoryList {

		if _, ok := tempItemGuidDict[inventory.ItemGuid]; ok {
			duplicateDataList = append(duplicateDataList, inventory)
		} else {
			newGuid := utils.NewGUID()

			invDoc := models.InventoryDoc{}

			invDoc.GuidFixed = newGuid
			invDoc.ShopID = shopID
			invDoc.Inventory = inventory

			invDoc.CreatedBy = authUsername
			invDoc.CreatedAt = time.Now()
			invDoc.LastUpdatedAt = time.Now()

			createDataList = append(createDataList, invDoc)
		}
	}
	return duplicateDataList, createDataList
}

func (svc InventoryService) CreateWithGuid(shopID string, authUsername string, guidFixed string, inventory models.Inventory) (string, error) {

	newGuid := guidFixed

	invDoc := models.InventoryDoc{}

	invDoc.GuidFixed = newGuid
	invDoc.ShopID = shopID
	invDoc.Inventory = inventory

	invDoc.CreatedBy = authUsername
	invDoc.CreatedAt = time.Now()
	invDoc.LastUpdatedAt = time.Now()

	mongoIdx, err := svc.invRepo.Create(invDoc)

	if err != nil {
		return "", err
	}

	if svc.invMqRepo != nil {
		err = svc.invMqRepo.Create(invDoc.InventoryData)

		if err != nil {
			return "", err
		}
	}

	svc.saveMasterSync(shopID)

	return mongoIdx, nil
}

// func (svc InventoryService) CreateBulk(shopID string, authUsername string, guidFixed string, inventories []models.Inventory) (error) {

// 	newGuid := guidFixed

// 	invDocList := make([]models.InventoryDoc, len(inventories))

// 	for index, inv := range inventories {
// 		invDoc := models.InventoryDoc{}

// 		invDoc.GuidFixed = newGuid
// 		invDoc.ShopID = shopID
// 		invDoc.Inventory = inv

// 		invDoc.CreatedBy = authUsername
// 		invDoc.CreatedAt = time.Now()

// 		invDocList[index] = invDoc
// 	}

// 	mongoIdx, err := svc.invRepo.Create(invDoc)

// 	if err != nil {
// 		return "", err
// 	}

// 	if svc.invMqRepo != nil {
// 		err = svc.invMqRepo.Create(invDoc.InventoryData)

// 		if err != nil {
// 			return "", err
// 		}
// 	}

// 	return mongoIdx, nil
// }

func (svc InventoryService) CreateInventory(shopID string, authUsername string, inventory models.Inventory) (string, string, error) {

	newGuid := utils.NewGUID()
	mongoIdx, err := svc.CreateWithGuid(shopID, authUsername, newGuid, inventory)
	return mongoIdx, newGuid, err
}

func (svc InventoryService) UpdateInventory(shopID string, guid string, authUsername string, inventory models.Inventory) error {

	findDoc, err := svc.invRepo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Inventory = inventory

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()
	findDoc.LastUpdatedAt = time.Now()

	err = svc.invRepo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	err = svc.invMqRepo.Update(findDoc.InventoryData)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc InventoryService) DeleteInventory(shopID string, guid string, username string) error {

	err := svc.invRepo.Delete(shopID, guid, username)

	if err != nil {
		return err
	}

	docIndentity := models.Identity{
		ShopID:    shopID,
		GuidFixed: guid,
	}

	err = svc.invMqRepo.Delete(docIndentity)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc InventoryService) InfoMongoInventory(id string) (models.InventoryInfo, error) {
	start := time.Now()

	idx, err := primitive.ObjectIDFromHex(id)
	findDoc, err := svc.invRepo.FindByID(idx)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.InventoryInfo{}, err
	}

	elapsed := time.Since(start)
	fmt.Printf("mongo :: pure id :: %s\n", elapsed)

	return findDoc.InventoryInfo, nil
}

func (svc InventoryService) InfoInventory(shopID string, guid string) (models.InventoryInfo, error) {

	findDoc, err := svc.invRepo.FindByGuid(shopID, guid)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.InventoryInfo{}, err
	}

	// if findDoc.ID == primitive.NilObjectID {
	// 	return models.InventoryInfo{}, nil
	// }

	return findDoc.InventoryInfo, nil
}

func (svc InventoryService) SearchInventory(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.invRepo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.InventoryInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc InventoryService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, paginate.PaginationData, error) {

	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.InventoryDeleteActivity
	var pagination1 paginate.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.invRepo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.InventoryActivity
	var pagination2 paginate.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.invRepo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
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

func (svc InventoryService) UpdateProductCategory(shopID string, authUsername string, catId string, guids []string) error {

	for _, guid := range guids {

		findDoc, err := svc.invRepo.FindByItemGuid(shopID, guid)

		if err != nil {
			return err
		}

		if findDoc.ID == primitive.NilObjectID {
			return errors.New("document not found")
		}

		findDoc.CategoryGuid = catId
		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()

		err = svc.invRepo.Update(shopID, findDoc.GuidFixed, findDoc)

		if err != nil {
			return err
		}

		err = svc.invMqRepo.Update(findDoc.InventoryData)

		if err != nil {
			return err
		}
	}

	if len(guids) > 0 {
		svc.saveMasterSync(shopID)
	}

	return nil
}

func (svc InventoryService) saveMasterSync(shopID string) {
	err := svc.cacheRepo.Save(shopID)

	if err != nil {
		fmt.Println("save inventory master cache error :: " + err.Error())
	}
}
