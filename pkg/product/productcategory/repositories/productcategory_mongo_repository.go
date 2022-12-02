package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/productcategory/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IProductCategoryRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.ProductCategoryDoc) (string, error)
	CreateInBatch(docList []models.ProductCategoryDoc) error
	Update(shopID string, guid string, doc models.ProductCategoryDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.ProductCategoryInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ProductCategoryDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ProductCategoryItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.ProductCategoryDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.ProductCategoryInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.ProductCategoryInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.ProductCategoryDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.ProductCategoryActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.ProductCategoryDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.ProductCategoryActivity, error)
}

type ProductCategoryRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ProductCategoryDoc]
	repositories.SearchRepository[models.ProductCategoryInfo]
	repositories.GuidRepository[models.ProductCategoryItemGuid]
	repositories.ActivityRepository[models.ProductCategoryActivity, models.ProductCategoryDeleteActivity]
}

func NewProductCategoryRepository(pst microservice.IPersisterMongo) *ProductCategoryRepository {

	insRepo := &ProductCategoryRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ProductCategoryDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ProductCategoryInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ProductCategoryItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.ProductCategoryActivity, models.ProductCategoryDeleteActivity](pst)

	return insRepo
}
