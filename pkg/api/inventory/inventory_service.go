package inventory

import (
	"errors"
	"fmt"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"sync"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryService interface {
	IsExistsGuid(shopID string, guidFixed string) (bool, error)
	CreateInBatch(shopID string, authUsername string, inventories []models.Inventory) error
	CreateIndex(doc models.InventoryIndex) error
	CreateWithGuid(shopID string, authUsername string, guidFixed string, inventory models.Inventory) (string, error)
	CreateInventory(shopID string, authUsername string, inventory models.Inventory) (string, string, error)
	UpdateInventory(shopID string, guid string, authUsername string, inventory models.Inventory) error
	DeleteInventory(shopID string, guid string, username string) error
	InfoInventory(shopID string, guid string) (models.InventoryInfo, error)
	InfoMongoInventory(id string) (models.InventoryInfo, error)
	InfoIndexInventory(shopID string, guid string) (models.InventoryInfo, error)
	SearchInventory(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error)
	LastActivityInventory(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, paginate.PaginationData, error)
	UpdateProductCategory(shopId string, authUsername string, catId string, guid []models.DocIdentity) error
}

type InventoryService struct {
	invRepo   IInventoryRepository
	invPgRepo IInventoryIndexPGRepository
	invMqRepo IInventoryMQRepository
}

func NewInventoryService(inventoryRepo IInventoryRepository, inventoryPgRepo IInventoryIndexPGRepository, inventoryMqRepo IInventoryMQRepository) InventoryService {
	return InventoryService{
		invRepo:   inventoryRepo,
		invPgRepo: inventoryPgRepo,
		invMqRepo: inventoryMqRepo,
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

func (svc InventoryService) CreateIndex(doc models.InventoryIndex) error {

	err := svc.invPgRepo.Create(doc)
	if err != nil {
		return err
	}

	return nil

}

func (svc InventoryService) CreateInBatch(shopID string, authUsername string, inventories []models.Inventory) error {

	tempInvDataList := []models.InventoryDoc{}

	for _, inventory := range inventories {
		newGuid := utils.NewGUID()

		invDoc := models.InventoryDoc{}

		invDoc.GuidFixed = newGuid
		invDoc.ShopID = shopID
		invDoc.Inventory = inventory

		invDoc.CreatedBy = authUsername
		invDoc.CreatedAt = time.Now()

		tempInvDataList = append(tempInvDataList, invDoc)
	}

	err := svc.invRepo.CreateInBatch(tempInvDataList)

	if err != nil {
		return err
	}

	return nil

}

func (svc InventoryService) CreateWithGuid(shopID string, authUsername string, guidFixed string, inventory models.Inventory) (string, error) {

	newGuid := guidFixed

	invDoc := models.InventoryDoc{}

	invDoc.GuidFixed = newGuid
	invDoc.ShopID = shopID
	invDoc.Inventory = inventory

	invDoc.CreatedBy = authUsername
	invDoc.CreatedAt = time.Now()

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

	err = svc.invRepo.Update(guid, findDoc)

	if err != nil {
		return err
	}

	err = svc.invMqRepo.Update(findDoc.InventoryData)

	if err != nil {
		return err
	}

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
	start := time.Now()
	findDoc, err := svc.invRepo.FindByGuid(shopID, guid)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.InventoryInfo{}, err
	}

	// if findDoc.ID == primitive.NilObjectID {
	// 	return models.InventoryInfo{}, nil
	// }

	elapsed := time.Since(start)
	fmt.Printf("mongo :: shopID,guidFixed :: %s\n", elapsed)

	return findDoc.InventoryInfo, nil
}

func (svc InventoryService) InfoIndexInventory(shopID string, guid string) (models.InventoryInfo, error) {
	start := time.Now()
	invIndex, err := svc.invPgRepo.FindByGuid(shopID, guid)

	idx, err := primitive.ObjectIDFromHex(invIndex.ID)

	findDoc, err := svc.invRepo.FindByID(idx)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.InventoryInfo{}, err
	}
	elapsed := time.Since(start)
	fmt.Printf("mongo,pg :: shopID,guidFixed :: %s\n", elapsed)
	return findDoc.InventoryInfo, nil
}

func (svc InventoryService) SearchInventory(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.invRepo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.InventoryInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc InventoryService) LastActivityInventory(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, paginate.PaginationData, error) {

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

	if pagination.TotalPage < pagination2.TotalPage {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc InventoryService) UpdateProductCategory(shopId string, authUsername string, catId string, guids []models.DocIdentity) error {

	for _, guid := range guids {

		findDoc, err := svc.invRepo.FindByGuid(shopId, guid.GuidFixed)

		if err != nil {
			return err
		}

		if findDoc.ID == primitive.NilObjectID {
			return errors.New("document not found")
		}

		findDoc.CategoryGuid = catId
		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()

		err = svc.invRepo.Update(guid.GuidFixed, findDoc)

		if err != nil {
			return err
		}

		err = svc.invMqRepo.Update(findDoc.InventoryData)

		if err != nil {
			return err
		}
	}
	return nil
}
