package kitchen

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IKitchenRepository interface {
	Count(shopID string) (int, error)
	Create(category restaurant.KitchenDoc) (string, error)
	CreateInBatch(inventories []restaurant.KitchenDoc) error
	Update(shopID string, guid string, category restaurant.KitchenDoc) error
	Delete(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]restaurant.KitchenInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (restaurant.KitchenDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]restaurant.KitchenItemGuid, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.KitchenDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.KitchenActivity, mongopagination.PaginationData, error)
}

type KitchenRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[restaurant.KitchenDoc]
	repositories.SearchRepository[restaurant.KitchenInfo]
	repositories.GuidRepository[restaurant.KitchenItemGuid]
	repositories.ActivityRepository[restaurant.KitchenActivity, restaurant.KitchenDeleteActivity]
}

func NewKitchenRepository(pst microservice.IPersisterMongo) KitchenRepository {

	insRepo := KitchenRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[restaurant.KitchenDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[restaurant.KitchenInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[restaurant.KitchenItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[restaurant.KitchenActivity, restaurant.KitchenDeleteActivity](pst)

	return insRepo
}
