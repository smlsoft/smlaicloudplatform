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
	Update(shopID string, guid string, category restaurant.ShopTableDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
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
	insRepo := ShopTableRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[restaurant.ShopTableDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[restaurant.ShopTableInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[restaurant.ShopTableItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[restaurant.ShopTableActivity, restaurant.ShopTableDeleteActivity](pst)

	return insRepo
}
