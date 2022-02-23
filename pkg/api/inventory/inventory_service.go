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
	CreateInventory(merchantId string, authUsername string, inventory models.Inventory) (string, error)
	UpdateInventory(guid string, merchantId string, authUsername string, inventory models.Inventory) error
	DeleteInventory(guid string, merchantId string) error
	InfoInventory(guid string, merchantId string) (models.Inventory, error)
	SearchInventory(merchantId string, q string, page int, limit int) ([]models.Inventory, paginate.PaginationData, error)
}

type InventoryService struct {
	invRepo IInventoryRepository
}

func NewInventoryService(inventoryRepo IInventoryRepository) IInventoryService {
	return &InventoryService{
		invRepo: inventoryRepo,
	}
}

func (svc *InventoryService) CreateInventory(merchantId string, authUsername string, inventory models.Inventory) (string, error) {

	newGuid := utils.NewGUID()

	inventory.GuidFixed = newGuid
	inventory.MerchantId = merchantId
	inventory.WaitType = 0
	inventory.Deleted = false
	inventory.CreatedBy = authUsername
	inventory.CreatedAt = time.Now()

	_, err := svc.invRepo.Create(inventory)

	if err != nil {
		return "", err
	}

	return newGuid, nil
}

func (svc *InventoryService) UpdateInventory(guid string, merchantId string, authUsername string, inventory models.Inventory) error {

	findDoc, err := svc.invRepo.FindByGuid(guid, merchantId)

	if err != nil {
		return err
	}

	if findDoc.Id == primitive.NilObjectID {
		return errors.New("guid invalid")
	}

	findDoc.ItemSku = inventory.ItemSku
	findDoc.CategoryGuid = inventory.CategoryGuid
	findDoc.LineNumber = inventory.LineNumber
	findDoc.Price = inventory.Price
	findDoc.Recommended = inventory.Recommended
	findDoc.HaveImage = inventory.HaveImage
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

func (svc *InventoryService) DeleteInventory(guid string, merchantId string) error {

	err := svc.invRepo.Delete(guid, merchantId)

	if err != nil {
		return err
	}
	return nil
}

func (svc *InventoryService) InfoInventory(guid string, merchantId string) (models.Inventory, error) {
	findDoc, err := svc.invRepo.FindByGuid(guid, merchantId)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.Inventory{}, err
	}

	if findDoc.Id == primitive.NilObjectID {
		return models.Inventory{}, errors.New("guid invalid")
	}

	return findDoc, nil
}

func (svc *InventoryService) SearchInventory(merchantId string, q string, page int, limit int) ([]models.Inventory, paginate.PaginationData, error) {
	docList, pagination, err := svc.invRepo.FindPage(merchantId, q, page, limit)

	if err != nil {
		return []models.Inventory{}, pagination, err
	}

	return docList, pagination, nil
}
