package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/category/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type ICategoryRepository interface {
	Count(shopID string) (int, error)
	Create(category models.CategoryDoc) (string, error)
	CreateInBatch(docList []models.CategoryDoc) error
	Update(shopID string, guid string, category models.CategoryDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters map[string]interface{}) (models.CategoryDoc, error)
	FindByGuid(shopID string, guid string) (models.CategoryDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.CategoryItemGuid, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.CategoryInfo, mongopagination.PaginationData, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryActivity, mongopagination.PaginationData, error)
}

type CategoryRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.CategoryDoc]
	repositories.SearchRepository[models.CategoryInfo]
	repositories.GuidRepository[models.CategoryItemGuid]
	repositories.ActivityRepository[models.CategoryActivity, models.CategoryDeleteActivity]
}

func NewCategoryRepository(pst microservice.IPersisterMongo) CategoryRepository {
	insRepo := CategoryRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.CategoryDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.CategoryInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.CategoryItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.CategoryActivity, models.CategoryDeleteActivity](pst)

	return insRepo

}
