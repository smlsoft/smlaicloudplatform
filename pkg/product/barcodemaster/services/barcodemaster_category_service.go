package services

import (
	"smlcloudplatform/pkg/product/barcodemaster/repositories"
	categoryRepo "smlcloudplatform/pkg/product/category/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IBarcodeMasterCategoryService interface {
	UpdateBarcodeMasterCategoryBulk(authUsername string, shopID string, catId string, guids []string) error
}

type BarcodeMasterCategoryService struct {
	barcodemasterRepository repositories.IBarcodeMasterRepository
	categoryRepository      categoryRepo.ICategoryRepository
	invMqRepo               repositories.IBarcodeMasterMQRepository
}

func NewBarcodeMastercategoryService(barcodemasterRepository repositories.BarcodeMasterRepository, categoryRepository categoryRepo.CategoryRepository, barcodemasterMqRepo repositories.IBarcodeMasterMQRepository) *BarcodeMasterCategoryService {
	return &BarcodeMasterCategoryService{
		barcodemasterRepository: barcodemasterRepository,
		categoryRepository:      categoryRepository,
		invMqRepo:               barcodemasterMqRepo,
	}
}

func (ics *BarcodeMasterCategoryService) UpdateBarcodeMasterCategoryBulk(shopID string, authUsername string, catId string, guids []string) error {

	// find category
	findCategory, err := ics.categoryRepository.FindByGuid(shopID, catId)
	if err != nil || findCategory.ID == primitive.NilObjectID {
		return err
	}

	itemsList, err := ics.barcodemasterRepository.FindByItemGuidList(shopID, guids)
	for _, findDoc := range itemsList {

		findDoc.CategoryGuid = catId
		findDoc.Category = &findCategory.Category
		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()

		err = ics.barcodemasterRepository.Update(shopID, findDoc.GuidFixed, findDoc)

		if err != nil {
			return err
		}

		err = ics.invMqRepo.Update(findDoc.BarcodeMasterData)

		if err != nil {
			return err
		}
	}

	return nil
}
