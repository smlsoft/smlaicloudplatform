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
	CreateInventory(shopID string, authUsername string, inventory models.Inventory) (string, error)
	UpdateInventory(guid string, shopID string, authUsername string, inventory models.Inventory) error
	DeleteInventory(guid string, shopID string) error
	InfoInventory(guid string, shopID string) (models.InventoryInfo, error)
	SearchInventory(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error)
}

type InventoryService struct {
	invRepo   IInventoryRepository
	invMqRepo IInventoryMQRepository
}

func NewInventoryService(inventoryRepo IInventoryRepository, inventoryMqRepo IInventoryMQRepository) InventoryService {
	return InventoryService{
		invRepo:   inventoryRepo,
		invMqRepo: inventoryMqRepo,
	}
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

	err = svc.invMqRepo.Create(invDoc.InventoryInfo)

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

	return nil
}

func (svc InventoryService) DeleteInventory(guid string, shopID string) error {

	err := svc.invRepo.Delete(guid, shopID)

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

	if findDoc.ID == primitive.NilObjectID {
		return models.InventoryInfo{}, errors.New("document not found")
	}

	return findDoc.InventoryInfo, nil
}

func (svc InventoryService) SearchInventory(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.invRepo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.InventoryInfo{}, pagination, err
	}

	return docList, pagination, nil
}
