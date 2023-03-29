package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/warehouse/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IWarehouseRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.WarehouseDoc) (string, error)
	CreateInBatch(docList []models.WarehouseDoc) error
	Update(shopID string, guid string, doc models.WarehouseDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.WarehouseInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.WarehouseDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.WarehouseItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.WarehouseDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.WarehouseInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.WarehouseInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.WarehouseDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.WarehouseActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.WarehouseDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.WarehouseActivity, error)
}

type WarehouseRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.WarehouseDoc]
	repositories.SearchRepository[models.WarehouseInfo]
	repositories.GuidRepository[models.WarehouseItemGuid]
	repositories.ActivityRepository[models.WarehouseActivity, models.WarehouseDeleteActivity]
}

func NewWarehouseRepository(pst microservice.IPersisterMongo) *WarehouseRepository {

	insRepo := &WarehouseRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.WarehouseDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.WarehouseInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.WarehouseItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.WarehouseActivity, models.WarehouseDeleteActivity](pst)

	return insRepo
}
