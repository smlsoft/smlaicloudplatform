package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/productgroup/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IProductGroupRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.ProductGroupDoc) (string, error)
	CreateInBatch(docList []models.ProductGroupDoc) error
	Update(shopID string, guid string, doc models.ProductGroupDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductGroupInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ProductGroupDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ProductGroupItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.ProductGroupDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductGroupInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.ProductGroupInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductGroupDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductGroupActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductGroupDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductGroupActivity, error)
}

type ProductGroupRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ProductGroupDoc]
	repositories.SearchRepository[models.ProductGroupInfo]
	repositories.GuidRepository[models.ProductGroupItemGuid]
	repositories.ActivityRepository[models.ProductGroupActivity, models.ProductGroupDeleteActivity]
}

func NewProductGroupRepository(pst microservice.IPersisterMongo) *ProductGroupRepository {

	insRepo := &ProductGroupRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ProductGroupDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ProductGroupInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ProductGroupItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.ProductGroupActivity, models.ProductGroupDeleteActivity](pst)

	return insRepo
}
