package inventory

import (
	"smlcloudplatform/pkg/api/category"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryCategoryService interface {
	UpdateInventoryCategoryBulk(authUsername string, shopId string, catId string, guids []string) error
}

type InventoryCategoryService struct {
	inventoryRepository IInventoryRepository
	categoryRepository  category.ICategoryRepository
	invMqRepo           IInventoryMQRepository
}

func NewInventorycategoryService(inventoryRepository InventoryRepository, categoryRepository category.CategoryRepository, inventoryMqRepo IInventoryMQRepository) *InventoryCategoryService {
	return &InventoryCategoryService{
		inventoryRepository: inventoryRepository,
		categoryRepository:  categoryRepository,
		invMqRepo:           inventoryMqRepo,
	}
}

func (ics *InventoryCategoryService) UpdateInventoryCategoryBulk(shopId string, authUsername string, catId string, guids []string) error {

	// find category
	findCategory, err := ics.categoryRepository.FindByCategoryGuid(shopId, catId)
	if err != nil || findCategory.ID == primitive.NilObjectID {
		return err
	}

	itemsList, err := ics.inventoryRepository.FindByItemGuidList(shopId, guids)
	for _, findDoc := range itemsList {

		findDoc.CategoryGuid = catId
		findDoc.Category = findCategory.Category
		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()

		err = ics.inventoryRepository.Update(findDoc.GuidFixed, findDoc)

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
