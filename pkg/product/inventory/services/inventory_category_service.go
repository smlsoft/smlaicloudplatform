package services

import (
	categoryRepo "smlcloudplatform/pkg/product/category/repositories"
	"smlcloudplatform/pkg/product/inventory/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryCategoryService interface {
	UpdateInventoryCategoryBulk(authUsername string, shopID string, catId string, guids []string) error
}

type InventoryCategoryService struct {
	inventoryRepository repositories.IInventoryRepository
	categoryRepository  categoryRepo.ICategoryRepository
	invMqRepo           repositories.IInventoryMQRepository
}

func NewInventorycategoryService(inventoryRepository repositories.InventoryRepository, categoryRepository categoryRepo.CategoryRepository, inventoryMqRepo repositories.IInventoryMQRepository) *InventoryCategoryService {
	return &InventoryCategoryService{
		inventoryRepository: inventoryRepository,
		categoryRepository:  categoryRepository,
		invMqRepo:           inventoryMqRepo,
	}
}

func (ics *InventoryCategoryService) UpdateInventoryCategoryBulk(shopID string, authUsername string, catId string, guids []string) error {

	// find category
	findCategory, err := ics.categoryRepository.FindByGuid(shopID, catId)
	if err != nil || findCategory.ID == primitive.NilObjectID {
		return err
	}

	itemsList, err := ics.inventoryRepository.FindByItemGuidList(shopID, guids)
	for _, findDoc := range itemsList {

		findDoc.CategoryGuid = catId
		findDoc.Category = &findCategory.Category
		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()

		err = ics.inventoryRepository.Update(shopID, findDoc.GuidFixed, findDoc)

		if err != nil {
			return err
		}

		err = ics.invMqRepo.Update(findDoc.InventoryData)

		if err != nil {
			return err
		}
	}

	return nil
}
