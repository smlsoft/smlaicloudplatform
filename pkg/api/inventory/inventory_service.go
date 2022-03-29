package inventory

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryService interface {
	IsExistsGuid(shopID string, guidFixed string) (bool, error)
	CreateIndex(doc models.InventoryIndex) error
	CreateWithGuid(shopID string, authUsername string, guidFixed string, inventory models.Inventory) (string, error)
	CreateInventory(shopID string, authUsername string, inventory models.Inventory) (string, error)
	UpdateInventory(guid string, shopID string, authUsername string, inventory models.Inventory) error
	DeleteInventory(guid string, shopID string, username string) error
	InfoInventory(guid string, shopID string) (models.InventoryInfo, error)
	SearchInventory(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error)
}

type InventoryService struct {
	invRepo   IInventoryRepository
	invPgRepo IInventoryPGRepository
	invMqRepo IInventoryMQRepository
}

func NewInventoryService(inventoryRepo IInventoryRepository, inventoryPgRepo IInventoryPGRepository, inventoryMqRepo IInventoryMQRepository) InventoryService {
	return InventoryService{
		invRepo:   inventoryRepo,
		invPgRepo: inventoryPgRepo,
		invMqRepo: inventoryMqRepo,
	}
}

// Find guid in postgresql index
func (svc InventoryService) IsExistsGuid(shopID string, guidFixed string) (bool, error) {

	count, err := svc.invPgRepo.Count(shopID, guidFixed)

	if err != nil {
		return false, err
	}

	if count == 0 {
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

	err = svc.invMqRepo.Create(invDoc.InventoryData)

	if err != nil {
		return "", err
	}

	return mongoIdx, nil
}

func (svc InventoryService) CreateInventory(shopID string, authUsername string, inventory models.Inventory) (string, error) {

	newGuid := utils.NewGUID()

	invDoc := models.InventoryDoc{}

	invDoc.GuidFixed = newGuid
	invDoc.ShopID = shopID
	invDoc.Inventory = inventory

	invDoc.CreatedBy = authUsername
	invDoc.CreatedAt = time.Now()

	_, err := svc.invRepo.Create(invDoc)

	if err != nil {
		return "", err
	}

	err = svc.invMqRepo.Create(invDoc.InventoryData)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc InventoryService) UpdateInventory(guid string, shopID string, authUsername string, inventory models.Inventory) error {

	findDoc, err := svc.invRepo.FindByGuid(guid, shopID)

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

func (svc InventoryService) DeleteInventory(guid string, shopID string, username string) error {

	err := svc.invRepo.Delete(guid, shopID, username)

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

func (svc InventoryService) InfoInventory(guid string, shopID string) (models.InventoryInfo, error) {
	findDoc, err := svc.invRepo.FindByGuid(guid, shopID)

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
