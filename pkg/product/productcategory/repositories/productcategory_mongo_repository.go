package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/productcategory/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IProductCategoryRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.ProductCategoryDoc) (string, error)
	CreateInBatch(docList []models.ProductCategoryDoc) error
	Update(shopID string, guid string, doc models.ProductCategoryDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductCategoryInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ProductCategoryDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ProductCategoryItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.ProductCategoryDoc, error)

	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductCategoryInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.ProductCategoryDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.ProductCategoryActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.ProductCategoryDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.ProductCategoryActivity, error)
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
