package inventory

import (
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
)

type ICategoryService interface {
	CreateCategory(merchantId string, authUsername string, category models.Category) (string, error)
	UpdateCategory(guid string, merchantId string, authUsername string, category models.Category) error
	DeleteCategory(guid string, merchantId string) error
	InfoCategory(guid string, merchantId string) (models.Category, error)
	SearchCategory(merchantId string, q string, page int, limit int) ([]models.Category, paginate.PaginationData, error)
}

type CategoryService struct {
	repo ICategoryRepository
}

func NewCategoryService(categoryRepository ICategoryRepository) ICategoryService {
	return &CategoryService{
		repo: categoryRepository,
	}
}

func (svc *CategoryService) CreateCategory(merchantId string, authUsername string, category models.Category) (string, error) {
	countCategory, err := svc.repo.Count(merchantId)

	if err != nil {
		return "", err
	}

	newGuidFixed := utils.NewGUID()
	category.MerchantId = merchantId
	category.GuidFixed = newGuidFixed
	category.LineNumber = int(countCategory) + 1
	category.CreatedBy = authUsername
	category.CreatedAt = time.Now()
	category.Deleted = false

	_, err = svc.repo.Create(category)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc *CategoryService) UpdateCategory(guid string, merchantId string, authUsername string, category models.Category) error {

	findDoc, err := svc.repo.FindByGuid(guid, merchantId)

	if err != nil {
		return err
	}

	findDoc.Name1 = category.Name1
	findDoc.HaveImage = category.HaveImage
	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc *CategoryService) DeleteCategory(guid string, merchantId string) error {
	err := svc.repo.Delete(guid, merchantId)

	if err != nil {
		return err
	}
	return nil
}

func (svc *CategoryService) InfoCategory(guid string, merchantId string) (models.Category, error) {

	findDoc, err := svc.repo.FindByGuid(guid, merchantId)

	if err != nil {
		return models.Category{}, err
	}

	return findDoc, nil

}

func (svc *CategoryService) SearchCategory(merchantId string, q string, page int, limit int) ([]models.Category, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(merchantId, q, page, limit)

	if err != nil {
		return []models.Category{}, pagination, err
	}

	return docList, pagination, nil
}
