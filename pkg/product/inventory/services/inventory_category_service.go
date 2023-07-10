package services

import (
	"context"
	"smlcloudplatform/pkg/product/inventory/repositories"
	categoryRepo "smlcloudplatform/pkg/product/productcategory/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryCategoryService interface {
	UpdateInventoryCategoryBulk(authUsername string, shopID string, catId string, guids []string) error
}

type InventoryCategoryService struct {
	inventoryRepository repositories.IInventoryRepository
	categoryRepository  categoryRepo.IProductCategoryRepository
	invMqRepo           repositories.IInventoryMQRepository
	contextTimeout      time.Duration
}

func NewInventorycategoryService(inventoryRepository repositories.InventoryRepository, categoryRepository categoryRepo.ProductCategoryRepository, inventoryMqRepo repositories.IInventoryMQRepository) *InventoryCategoryService {

	contextTimeout := time.Duration(15) * time.Second

	return &InventoryCategoryService{
		inventoryRepository: inventoryRepository,
		categoryRepository:  categoryRepository,
		invMqRepo:           inventoryMqRepo,
		contextTimeout:      contextTimeout,
	}
}

func (svc InventoryCategoryService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc *InventoryCategoryService) UpdateInventoryCategoryBulk(shopID string, authUsername string, catId string, guids []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	// find category
	findCategory, err := svc.categoryRepository.FindByGuid(ctx, shopID, catId)
	if err != nil || findCategory.ID == primitive.NilObjectID {
		return err
	}

	itemsList, err := svc.inventoryRepository.FindByItemGuidList(ctx, shopID, guids)
	for _, findDoc := range itemsList {

		findDoc.CategoryGuid = catId
		findDoc.Category = &findCategory.ProductCategory
		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()

		err = svc.inventoryRepository.Update(ctx, shopID, findDoc.GuidFixed, findDoc)

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
