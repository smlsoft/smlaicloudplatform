package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/ordertype/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IOrderTypeRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.OrderTypeDoc) (string, error)
	CreateInBatch(docList []models.OrderTypeDoc) error
	Update(shopID string, guid string, doc models.OrderTypeDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.OrderTypeInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.OrderTypeDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.OrderTypeItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.OrderTypeDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.OrderTypeInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.OrderTypeInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.OrderTypeDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.OrderTypeActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.OrderTypeDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.OrderTypeActivity, error)
}

type OrderTypeRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.OrderTypeDoc]
	repositories.SearchRepository[models.OrderTypeInfo]
	repositories.GuidRepository[models.OrderTypeItemGuid]
	repositories.ActivityRepository[models.OrderTypeActivity, models.OrderTypeDeleteActivity]
}

func NewOrderTypeRepository(pst microservice.IPersisterMongo) *OrderTypeRepository {

	insRepo := &OrderTypeRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.OrderTypeDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.OrderTypeInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.OrderTypeItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.OrderTypeActivity, models.OrderTypeDeleteActivity](pst)

	return insRepo
}
