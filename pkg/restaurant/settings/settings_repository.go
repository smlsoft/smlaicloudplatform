package settings

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/settings/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IRestaurantSettingsRepository interface {
	Count(shopID string) (int, error)
	Create(category models.RestaurantSettingsDoc) (string, error)
	CreateInBatch(docList []models.RestaurantSettingsDoc) error
	Update(shopID string, guid string, category models.RestaurantSettingsDoc) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.RestaurantSettingsInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.RestaurantSettingsDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.RestaurantSettingsItemGuid, error)
	FindOne(shopID string, filters interface{}) (models.RestaurantSettingsDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.RestaurantSettingsInfo, mongopagination.PaginationData, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.RestaurantSettingsDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.RestaurantSettingsActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.RestaurantSettingsDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.RestaurantSettingsActivity, error)
}

type RestaurantSettingsRepository struct {
	pst microservice.IPersisterMongo

	repositories.CrudRepository[models.RestaurantSettingsDoc]
	repositories.SearchRepository[models.RestaurantSettingsInfo]
	repositories.GuidRepository[models.RestaurantSettingsItemGuid]
	repositories.ActivityRepository[models.RestaurantSettingsActivity, models.RestaurantSettingsDeleteActivity]
}

func NewRestaurantSettingsRepository(pst microservice.IPersisterMongo) RestaurantSettingsRepository {
	tempRepo := RestaurantSettingsRepository{
		pst: pst,
	}

	tempRepo.CrudRepository = repositories.NewCrudRepository[models.RestaurantSettingsDoc](pst)
	tempRepo.SearchRepository = repositories.NewSearchRepository[models.RestaurantSettingsInfo](pst)
	tempRepo.GuidRepository = repositories.NewGuidRepository[models.RestaurantSettingsItemGuid](pst)
	tempRepo.ActivityRepository = repositories.NewActivityRepository[models.RestaurantSettingsActivity, models.RestaurantSettingsDeleteActivity](pst)

	return tempRepo
}
