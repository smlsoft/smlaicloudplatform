package inventory

import (
	"encoding/json"
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryService interface {
	CreateInventory(shopId string, authUsername string, inventory models.Inventory) (string, error)
	UpdateInventory(guid string, shopId string, authUsername string, inventory models.Inventory) error
	DeleteInventory(guid string, shopId string) error
	InfoInventory(guid string, shopId string) (models.Inventory, error)
	SearchInventory(shopId string, q string, page int, limit int) ([]models.Inventory, paginate.PaginationData, error)
}

type InventoryService struct {
	invRepo   IInventoryRepository
	invMqRepo IInventoryMQRepository
}

func NewInventoryService(inventoryRepo IInventoryRepository, inventoryMqRepo IInventoryMQRepository) IInventoryService {
	return &InventoryService{
		invRepo:   inventoryRepo,
		invMqRepo: inventoryMqRepo,
	}
}

func (svc *InventoryService) CreateInventory(shopId string, authUsername string, inventory models.Inventory) (string, error) {

	newGuid := utils.NewGUID()

	inventory.GuidFixed = newGuid
	inventory.ShopId = shopId
	inventory.Deleted = false
	inventory.CreatedBy = authUsername
	inventory.CreatedAt = time.Now()

	_, err := svc.invRepo.Create(inventory)

	if err != nil {
		return "", err
	}

	invStr, err := json.Marshal(inventory)
	if err != nil {
		return "", err
	}

	invRequest := &models.InventoryRequest{}
	err = json.Unmarshal([]byte(invStr), invRequest)

	if err != nil {
		return "", err
	}

	err = svc.invMqRepo.Create(*invRequest)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc *InventoryService) UpdateInventory(guid string, shopId string, authUsername string, inventory models.Inventory) error {

	findDoc, err := svc.invRepo.FindByGuid(guid, shopId)

	if err != nil {
		return err
	}

	if findDoc.Id == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ItemSku = inventory.ItemSku
	findDoc.CategoryGuid = inventory.CategoryGuid
	findDoc.Price = inventory.Price
	findDoc.Recommended = inventory.Recommended
	findDoc.Activated = inventory.Activated

	findDoc.Name1 = inventory.Name1
	findDoc.Name2 = inventory.Name2
	findDoc.Name3 = inventory.Name3
	findDoc.Name4 = inventory.Name4
	findDoc.Name5 = inventory.Name5

	findDoc.Description1 = inventory.Description1
	findDoc.Description2 = inventory.Description2
	findDoc.Description3 = inventory.Description3
	findDoc.Description4 = inventory.Description4
	findDoc.Description5 = inventory.Description5

	inventory.UpdatedBy = authUsername
	inventory.UpdatedAt = time.Now()

	err = svc.invRepo.Update(guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc *InventoryService) DeleteInventory(guid string, shopId string) error {

	err := svc.invRepo.Delete(guid, shopId)

	if err != nil {
		return err
	}
	return nil
}

func (svc *InventoryService) InfoInventory(guid string, shopId string) (models.Inventory, error) {
	findDoc, err := svc.invRepo.FindByGuid(guid, shopId)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.Inventory{}, err
	}

	if findDoc.Id == primitive.NilObjectID {
		return models.Inventory{}, errors.New("document not found")
	}

	return findDoc, nil
}

func (svc *InventoryService) SearchInventory(shopId string, q string, page int, limit int) ([]models.Inventory, paginate.PaginationData, error) {
	docList, pagination, err := svc.invRepo.FindPage(shopId, q, page, limit)

	if err != nil {
		return []models.Inventory{}, pagination, err
	}

	return docList, pagination, nil
}
