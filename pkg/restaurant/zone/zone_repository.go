package zone

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/zone/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IZoneRepository interface {
	Count(shopID string) (int, error)
	Create(category models.ZoneDoc) (string, error)
	CreateInBatch(docList []models.ZoneDoc) error
	Update(shopID string, guid string, category models.ZoneDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ZoneInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ZoneDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ZoneItemGuid, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ZoneDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ZoneActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ZoneDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ZoneActivity, error)
}

type ZoneRepository struct {
	pst microservice.IPersisterMongo

	repositories.CrudRepository[models.ZoneDoc]
	repositories.SearchRepository[models.ZoneInfo]
	repositories.GuidRepository[models.ZoneItemGuid]
	repositories.ActivityRepository[models.ZoneActivity, models.ZoneDeleteActivity]
}

func NewZoneRepository(pst microservice.IPersisterMongo) ZoneRepository {
	tempRepo := ZoneRepository{
		pst: pst,
	}

	tempRepo.CrudRepository = repositories.NewCrudRepository[models.ZoneDoc](pst)
	tempRepo.SearchRepository = repositories.NewSearchRepository[models.ZoneInfo](pst)
	tempRepo.GuidRepository = repositories.NewGuidRepository[models.ZoneItemGuid](pst)
	tempRepo.ActivityRepository = repositories.NewActivityRepository[models.ZoneActivity, models.ZoneDeleteActivity](pst)

	return tempRepo
}
