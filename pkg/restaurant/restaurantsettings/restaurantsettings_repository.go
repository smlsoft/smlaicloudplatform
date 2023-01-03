package restaurantsettings

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/restaurantsettings/models"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IRestaurantSettingsRepository interface {
	Count(shopID string) (int, error)
	Create(category models.RestaurantSettingsDoc) (string, error)
	CreateInBatch(docList []models.RestaurantSettingsDoc) error
	Update(shopID string, guid string, category models.RestaurantSettingsDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.RestaurantSettingsInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.RestaurantSettingsDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.RestaurantSettingsItemGuid, error)
	FindOne(shopID string, filters interface{}) (models.RestaurantSettingsDoc, error)
	FindPageFilterSort(shopID string, filters map[string]interface{}, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.RestaurantSettingsInfo, mongopagination.PaginationData, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.RestaurantSettingsDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.RestaurantSettingsActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.RestaurantSettingsDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.RestaurantSettingsActivity, error)
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
