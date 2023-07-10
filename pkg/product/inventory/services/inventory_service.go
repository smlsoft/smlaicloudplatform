package services

import (
	"context"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/inventory/models"
	"smlcloudplatform/pkg/product/inventory/repositories"
	"smlcloudplatform/pkg/utils"
	"sync"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryService interface {
	IsExistsGuid(shopID string, guidFixed string) (bool, error)
	CreateInBatch(shopID string, authUsername string, inventories []models.Inventory) (models.InventoryBulkImport, error)
	CreateWithGuid(ctx context.Context, shopID string, authUsername string, guidFixed string, inventory models.Inventory) (string, error)
	CreateInventory(shopID string, authUsername string, inventory models.Inventory) (string, string, error)
	UpdateInventory(shopID string, findDoc models.InventoryDoc, authUsername string, inventory models.Inventory) error
	UpdateInventoryByGuidfixed(shopID string, guid string, authUsername string, inventory models.Inventory) error
	UpdateInventoryByItemCode(shopID string, itemCode string, authUsername string, inventory models.Inventory) error
	DeleteInventory(shopID string, guid string, username string) error
	SaveInventory(shopID string, authUsername string, inventory models.Inventory) error
	InfoInventory(shopID string, guid string) (models.InventoryInfo, error)
	InfoInventoryItemCode(shopID string, itemCode string) (models.InventoryInfo, error)
	InfoMongoInventory(id string) (models.InventoryInfo, error)
	InfoInventoryBarcode(shopID string, barcode string) (models.InventoryInfo, error)
	SearchInventory(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.InventoryInfo, mongopagination.PaginationData, error)
	UpdateProductCategory(shopID string, authUsername string, catId string, guid []string) error

	LastActivity(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) (common.LastActivity, mongopagination.PaginationData, error)
	LastActivityStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) (common.LastActivity, error)
	GetModuleName() string
}

type InventoryService struct {
	invRepo        repositories.IInventoryRepository
	invMqRepo      repositories.IInventoryMQRepository
	syncCacheRepo  mastersync.IMasterSyncCacheRepository
	contextTimeout time.Duration
}

func NewInventoryService(inventoryRepo repositories.IInventoryRepository, inventoryMqRepo repositories.IInventoryMQRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) InventoryService {

	contextTimeout := time.Duration(15) * time.Second

	return InventoryService{
		invRepo:        inventoryRepo,
		invMqRepo:      inventoryMqRepo,
		syncCacheRepo:  syncCacheRepo,
		contextTimeout: contextTimeout,
	}
}

func (svc InventoryService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc InventoryService) IsExistsGuid(shopID string, guidFixed string) (bool, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.invRepo.FindByGuid(ctx, shopID, guidFixed)

	if err != nil {
		return false, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return false, nil
	}

	return true, nil
}

func (svc InventoryService) CreateInBatch(shopID string, authUsername string, inventories []models.Inventory) (models.InventoryBulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	createDataList := []models.InventoryDoc{}
	duplicateDataList := []models.Inventory{}

	payloadInventoryList, payloadDuplicateInventoryList := filterDuplicateInventory(inventories)

	itemCodeGuidList := []string{}
	for _, inventory := range payloadInventoryList {
		itemCodeGuidList = append(itemCodeGuidList, inventory.ItemGuid)
	}

	findItemGuid, err := svc.invRepo.FindByItemCodeGuid(ctx, shopID, itemCodeGuidList)

	if err != nil {
		return models.InventoryBulkImport{}, err
	}

	duplicateDataList, createDataList = preparePayloadDataInventory(shopID, authUsername, findItemGuid, payloadInventoryList)

	updateSuccessDataList, updateFailDataList := updateOnDuplicateInventory(ctx, shopID, authUsername, duplicateDataList, svc.invRepo)

	if len(createDataList) > 0 {
		err = svc.invRepo.CreateInBatch(ctx, createDataList)

		if err != nil {
			return models.InventoryBulkImport{}, err
		}
	}
	createDataKey := []string{}

	for _, inv := range createDataList {
		createDataKey = append(createDataKey, inv.ItemGuid)

		// reply kafka
		if svc.invMqRepo != nil {
			err = svc.invMqRepo.Create(ctx, inv.InventoryData)

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

func updateOnDuplicateInventory(ctx context.Context, shopID string, authUsername string, duplicateDataList []models.Inventory, repo repositories.IInventoryRepository) ([]models.InventoryDoc, []models.Inventory) {
	updateSuccessDataList := []models.InventoryDoc{}
	updateFailDataList := []models.Inventory{}
	for _, inv := range duplicateDataList {
		findDoc, err := repo.FindByItemGuid(ctx, shopID, inv.ItemGuid)

		if err != nil || findDoc.ID == primitive.NilObjectID {
			updateFailDataList = append(updateFailDataList, inv)
			continue
		}

		findDoc.Inventory = inv

		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()
		findDoc.LastUpdatedAt = time.Now()

		err = repo.Update(ctx, shopID, findDoc.GuidFixed, findDoc)

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

func (svc InventoryService) SaveInventory(shopID string, authUsername string, inventory models.Inventory) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, _ := svc.invRepo.FindByItemCode(ctx, shopID, inventory.ItemCode)

	if len(findDoc.GuidFixed) < 1 {
		newGuid := utils.NewGUID()
		_, err := svc.CreateWithGuid(ctx, shopID, authUsername, newGuid, inventory)
		return err
	} else {
		return svc.UpdateInventory(shopID, findDoc, authUsername, inventory)
	}

}

func (svc InventoryService) CreateWithGuid(ctx context.Context, shopID string, authUsername string, guidFixed string, inventory models.Inventory) (string, error) {

	if inventory.Barcodes != nil {
		reqBarcodes := []string{}
		for _, barcode := range *inventory.Barcodes {
			reqBarcodes = append(reqBarcodes, barcode.Barcode)
		}

		findDocBarcodes, err := svc.invRepo.FindByBarcodes(ctx, shopID, reqBarcodes)

		if err != nil {
			return "", err
		}

		if len(findDocBarcodes) > 0 {
			tempBarcode := *findDocBarcodes[0].Barcodes
			return "", fmt.Errorf("barcode '%s' is exists", tempBarcode[0].Barcode)
		}
	}

	findDoc, _ := svc.invRepo.FindByItemCode(ctx, shopID, inventory.ItemCode)

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("item code is exists")
	}

	newGuid := guidFixed

	invDoc := models.InventoryDoc{}

	invDoc.GuidFixed = newGuid
	invDoc.ShopID = shopID
	invDoc.Inventory = inventory

	invDoc.CreatedBy = authUsername
	invDoc.CreatedAt = time.Now()
	invDoc.LastUpdatedAt = time.Now()

	mongoIdx, err := svc.invRepo.Create(ctx, invDoc)

	if err != nil {
		return "", err
	}

	if svc.invMqRepo != nil {
		err = svc.invMqRepo.Create(ctx, invDoc.InventoryData)

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

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	newGuid := utils.NewGUID()
	mongoIdx, err := svc.CreateWithGuid(ctx, shopID, authUsername, newGuid, inventory)
	return mongoIdx, newGuid, err
}

func (svc InventoryService) UpdateInventoryByGuidfixed(shopID string, guid string, authUsername string, inventory models.Inventory) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.invRepo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	return svc.UpdateInventory(shopID, findDoc, authUsername, inventory)
}

func (svc InventoryService) UpdateInventoryByItemCode(shopID string, itemCode string, authUsername string, inventory models.Inventory) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.invRepo.FindByItemCode(ctx, shopID, itemCode)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	return svc.UpdateInventory(shopID, findDoc, authUsername, inventory)
}

func (svc InventoryService) UpdateInventory(shopID string, findDoc models.InventoryDoc, authUsername string, inventory models.Inventory) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	tempItemCode := findDoc.ItemCode

	if inventory.Barcodes != nil && len(*inventory.Barcodes) > 0 {
		reqBarcodes := []string{}
		idxReqBarcode := map[string]struct{}{}

		for _, barcode := range *inventory.Barcodes {
			reqBarcodes = append(reqBarcodes, barcode.Barcode)
			idxReqBarcode[barcode.Barcode] = struct{}{}
		}

		findDocBarcodes, err := svc.invRepo.FindByBarcodes(ctx, shopID, reqBarcodes)

		if err != nil {
			return err
		}

		for _, doc := range findDocBarcodes {
			for _, barcode := range *doc.Barcodes {

				if doc.GuidFixed != findDoc.GuidFixed {
					if _, ok := idxReqBarcode[barcode.Barcode]; ok {
						return fmt.Errorf("barcode %s is exists", barcode.Barcode)
					}
				}
			}
		}
	}

	findDoc.Inventory = inventory

	if len(tempItemCode) > 0 {
		findDoc.ItemCode = tempItemCode
	}

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()
	findDoc.LastUpdatedAt = time.Now()

	err := svc.invRepo.Update(ctx, shopID, findDoc.GuidFixed, findDoc)

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

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	err := svc.invRepo.Delete(ctx, shopID, guid, username)

	if err != nil {
		return err
	}

	docIndentity := common.Identity{
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

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	start := time.Now()

	idx, err := primitive.ObjectIDFromHex(id)
	findDoc, err := svc.invRepo.FindByID(ctx, idx)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.InventoryInfo{}, err
	}

	elapsed := time.Since(start)
	fmt.Printf("mongo :: pure id :: %s\n", elapsed)

	return findDoc.InventoryInfo, nil
}

func (svc InventoryService) InfoInventoryBarcode(shopID string, barcode string) (models.InventoryInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.invRepo.FindByItemBarcode(ctx, shopID, barcode)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.InventoryInfo{}, err
	}

	// if findDoc.ID == primitive.NilObjectID {
	// 	return models.InventoryInfo{}, nil
	// }

	return findDoc.InventoryInfo, nil
}

func (svc InventoryService) InfoInventory(shopID string, guid string) (models.InventoryInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.invRepo.FindByGuid(ctx, shopID, guid)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.InventoryInfo{}, err
	}

	// if findDoc.ID == primitive.NilObjectID {
	// 	return models.InventoryInfo{}, nil
	// }

	return findDoc.InventoryInfo, nil
}

func (svc InventoryService) InfoInventoryItemCode(shopID string, itemCode string) (models.InventoryInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.invRepo.FindByItemCode(ctx, shopID, itemCode)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.InventoryInfo{}, err
	}

	// if findDoc.ID == primitive.NilObjectID {
	// 	return models.InventoryInfo{}, nil
	// }

	return findDoc.InventoryInfo, nil
}

func (svc InventoryService) SearchInventory(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.InventoryInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList, pagination, err := svc.invRepo.FindPage(ctx, shopID, filters, pageable)

	if err != nil {
		return []models.InventoryInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc InventoryService) UpdateProductCategory(shopID string, authUsername string, catId string, guids []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	for _, guid := range guids {

		findDoc, err := svc.invRepo.FindByItemGuid(ctx, shopID, guid)

		if err != nil {
			return err
		}

		if findDoc.ID == primitive.NilObjectID {
			return errors.New("document not found")
		}

		findDoc.CategoryGuid = catId
		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()

		err = svc.invRepo.Update(ctx, shopID, findDoc.GuidFixed, findDoc)

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

func (svc InventoryService) LastActivity(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) (common.LastActivity, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.InventoryDeleteActivity
	var pagination1 mongopagination.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.invRepo.FindDeletedPage(ctx, shopID, lastUpdatedDate, map[string]interface{}{}, pageable)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.InventoryActivity
	var pagination2 mongopagination.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.invRepo.FindCreatedOrUpdatedPage(ctx, shopID, lastUpdatedDate, map[string]interface{}{}, pageable)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return common.LastActivity{}, pagination1, err1
	}

	if err2 != nil {
		return common.LastActivity{}, pagination2, err2
	}

	lastActivity := common.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	pagination := pagination1

	if pagination.Total < pagination2.Total {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc InventoryService) LastActivityStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) (common.LastActivity, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.InventoryDeleteActivity
	var err1 error

	go func() {
		deleteDocList, err1 = svc.invRepo.FindDeletedStep(ctx, shopID, lastUpdatedDate, map[string]interface{}{}, pageableStep)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.InventoryActivity

	var err2 error

	go func() {
		createAndUpdateDocList, err2 = svc.invRepo.FindCreatedOrUpdatedStep(ctx, shopID, lastUpdatedDate, map[string]interface{}{}, pageableStep)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return common.LastActivity{}, err1
	}

	if err2 != nil {
		return common.LastActivity{}, err2
	}

	lastActivity := common.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	return lastActivity, nil
}

func (svc InventoryService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc InventoryService) GetModuleName() string {
	return "inventory"
}
