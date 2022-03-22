package category

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICategoryService interface {
	CreateCategory(shopID string, authUsername string, category models.Category) (string, error)
	UpdateCategory(guid string, shopID string, authUsername string, category models.Category) error
	DeleteCategory(guid string, shopID string) error
	InfoCategory(guid string, shopID string) (models.CategoryInfo, error)
	SearchCategory(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error)
}

type CategoryService struct {
	repo ICategoryRepository
}

func NewCategoryService(categoryRepository ICategoryRepository) CategoryService {
	return CategoryService{
		repo: categoryRepository,
	}
}

func (svc CategoryService) CreateCategory(shopID string, authUsername string, category models.Category) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.CategoryDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Category = category

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc CategoryService) UpdateCategory(guid string, shopID string, authUsername string, category models.Category) error {

	findDoc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Category = category

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc CategoryService) DeleteCategory(guid string, shopID string) error {
	err := svc.repo.Delete(guid, shopID)

	if err != nil {
		return err
	}
	return nil
}

func (svc CategoryService) InfoCategory(guid string, shopID string) (models.CategoryInfo, error) {

	findDoc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return models.CategoryInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.CategoryInfo{}, errors.New("document not found")
	}

	return findDoc.CategoryInfo, nil

}

func (svc CategoryService) SearchCategory(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.CategoryInfo{}, pagination, err
	}

	return docList, pagination, nil
}
