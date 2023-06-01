package shopzone

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/shopzone/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IShopZoneRepository interface {
	Count(shopID string) (int, error)
	Create(category models.ShopZoneDoc) (string, error)
	CreateInBatch(docList []models.ShopZoneDoc) error
	Update(shopID string, guid string, category models.ShopZoneDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ShopZoneInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ShopZoneDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ShopZoneItemGuid, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ShopZoneDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ShopZoneActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ShopZoneDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ShopZoneActivity, error)
}

type ShopZoneRepository struct {
	pst microservice.IPersisterMongo

	repositories.CrudRepository[models.ShopZoneDoc]
	repositories.SearchRepository[models.ShopZoneInfo]
	repositories.GuidRepository[models.ShopZoneItemGuid]
	repositories.ActivityRepository[models.ShopZoneActivity, models.ShopZoneDeleteActivity]
}

func NewShopZoneRepository(pst microservice.IPersisterMongo) ShopZoneRepository {
	tempRepo := ShopZoneRepository{
		pst: pst,
	}

	tempRepo.CrudRepository = repositories.NewCrudRepository[models.ShopZoneDoc](pst)
	tempRepo.SearchRepository = repositories.NewSearchRepository[models.ShopZoneInfo](pst)
	tempRepo.GuidRepository = repositories.NewGuidRepository[models.ShopZoneItemGuid](pst)
	tempRepo.ActivityRepository = repositories.NewActivityRepository[models.ShopZoneActivity, models.ShopZoneDeleteActivity](pst)

	return tempRepo
}
