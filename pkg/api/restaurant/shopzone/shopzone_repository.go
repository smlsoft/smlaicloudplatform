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
	Update(shopID string, guid string, category restaurant.ShopZoneDoc) error
	Delete(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]restaurant.ShopZoneInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (restaurant.ShopZoneDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]restaurant.ShopZoneItemGuid, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.ShopZoneDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.ShopZoneActivity, mongopagination.PaginationData, error)
}

type ShopZoneRepository struct {
	pst microservice.IPersisterMongo

	repositories.CrudRepository[restaurant.ShopZoneDoc]
	repositories.SearchRepository[restaurant.ShopZoneInfo]
	repositories.GuidRepository[restaurant.ShopZoneItemGuid]
	repositories.ActivityRepository[restaurant.ShopZoneActivity, restaurant.ShopZoneDeleteActivity]
}

func NewShopZoneRepository(pst microservice.IPersisterMongo) ShopZoneRepository {
	tempRepo := ShopZoneRepository{
		pst: pst,
	}

	tempRepo.CrudRepository = repositories.NewCrudRepository[restaurant.ShopZoneDoc](pst)
	tempRepo.SearchRepository = repositories.NewSearchRepository[restaurant.ShopZoneInfo](pst)
	tempRepo.GuidRepository = repositories.NewGuidRepository[restaurant.ShopZoneItemGuid](pst)
	tempRepo.ActivityRepository = repositories.NewActivityRepository[restaurant.ShopZoneActivity, restaurant.ShopZoneDeleteActivity](pst)

	return tempRepo
}
