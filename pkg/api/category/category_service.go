package category

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"sync"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICategoryService interface {
	CreateCategory(shopID string, authUsername string, category models.Category) (string, error)
	UpdateCategory(guid string, shopID string, authUsername string, category models.Category) error
	DeleteCategory(guid string, shopID string, authUsername string) error
	InfoCategory(guid string, shopID string) (models.CategoryInfo, error)
	SearchCategory(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error)
	LastActivityCategory(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, paginate.PaginationData, error)
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

func (svc CategoryService) DeleteCategory(guid string, shopID string, authUsername string) error {
	err := svc.repo.Delete(guid, shopID, authUsername)

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

func (svc CategoryService) LastActivityCategory(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, paginate.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.CategoryDeleteActivity
	var pagination1 paginate.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.repo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.CategoryActivity
	var pagination2 paginate.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.repo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return models.LastActivity{}, pagination1, err1
	}

	if err2 != nil {
		return models.LastActivity{}, pagination2, err2
	}

	lastActivity := models.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	pagination := pagination1

	if pagination.TotalPage < pagination2.TotalPage {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}
