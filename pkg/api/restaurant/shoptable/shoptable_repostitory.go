package shoptable

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IShopTableRepository interface {
	Count(shopID string) (int, error)
	Create(category restaurant.ShopTableDoc) (string, error)
	CreateInBatch(inventories []restaurant.ShopTableDoc) error
	Update(guid string, category restaurant.ShopTableDoc) error
	Delete(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]restaurant.ShopTableInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (restaurant.ShopTableDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]restaurant.ShopTableItemGuid, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.ShopTableDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.ShopTableActivity, mongopagination.PaginationData, error)
}

type ShopTableRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[restaurant.ShopTableDoc]
	repositories.SearchRepository[restaurant.ShopTableInfo]
	repositories.GuidRepository[restaurant.ShopTableItemGuid]
	repositories.ActivityRepository[restaurant.ShopTableActivity, restaurant.ShopTableDeleteActivity]
}

func NewShopTableRepository(pst microservice.IPersisterMongo) ShopTableRepository {
	return ShopTableRepository{
		pst: pst,
	}
}
