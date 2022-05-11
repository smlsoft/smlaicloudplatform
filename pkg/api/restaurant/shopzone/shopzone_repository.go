package shopzone

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IShopZoneRepository interface {
	Count(shopID string) (int, error)
	Create(category restaurant.ShopZoneDoc) (string, error)
	CreateInBatch(inventories []restaurant.ShopZoneDoc) error
	Update(guid string, category restaurant.ShopZoneDoc) error
	Delete(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]restaurant.ShopZoneInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (restaurant.ShopZoneDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]restaurant.ShopZoneItemGuid, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.ShopZoneDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.ShopZoneActivity, mongopagination.PaginationData, error)
}

type ShopZoneRepository struct {
	pst          microservice.IPersisterMongo
	crudRepo     repositories.CrudRepository[restaurant.ShopZoneDoc]
	searchRepo   repositories.SearchRepository[restaurant.ShopZoneInfo]
	guidRepo     repositories.GuidRepository[restaurant.ShopZoneItemGuid]
	activityRepo repositories.ActivityRepository[restaurant.ShopZoneActivity, restaurant.ShopZoneDeleteActivity]
}

func NewShopZoneRepository(pst microservice.IPersisterMongo) ShopZoneRepository {
	crudRepo := repositories.NewCrudRepository[restaurant.ShopZoneDoc](pst)
	searchRepo := repositories.NewSearchRepository[restaurant.ShopZoneInfo](pst)
	guidRepo := repositories.NewGuidRepository[restaurant.ShopZoneItemGuid](pst)
	activityRepo := repositories.NewActivityRepository[restaurant.ShopZoneActivity, restaurant.ShopZoneDeleteActivity](pst)

	return ShopZoneRepository{
		pst:          pst,
		crudRepo:     crudRepo,
		searchRepo:   searchRepo,
		guidRepo:     guidRepo,
		activityRepo: activityRepo,
	}
}
