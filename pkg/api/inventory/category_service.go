package inventory

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
	InfoCategory(guid string, shopID string) (models.Category, error)
	SearchCategory(shopID string, q string, page int, limit int) ([]models.Category, paginate.PaginationData, error)
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
	category.ShopID = shopID
	category.GuidFixed = newGuidFixed
	category.CreatedBy = authUsername
	category.CreatedAt = time.Now()
	category.Deleted = false

	_, err := svc.repo.Create(category)

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

	findDoc.Name1 = category.Name1
	findDoc.Name2 = category.Name2
	findDoc.Name3 = category.Name3
	findDoc.Name4 = category.Name4
	findDoc.Name5 = category.Name5
	findDoc.Image = category.Image
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

func (svc CategoryService) InfoCategory(guid string, shopID string) (models.Category, error) {

	findDoc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return models.Category{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.Category{}, errors.New("document not found")
	}

	return findDoc, nil

}

func (svc CategoryService) SearchCategory(shopID string, q string, page int, limit int) ([]models.Category, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.Category{}, pagination, err
	}

	return docList, pagination, nil
}
