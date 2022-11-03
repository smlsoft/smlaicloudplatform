package kitchen

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/kitchen/models"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IKitchenRepository interface {
	Count(shopID string) (int, error)
	Create(category models.KitchenDoc) (string, error)
	CreateInBatch(docList []models.KitchenDoc) error
	Update(shopID string, guid string, category models.KitchenDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.KitchenInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.KitchenDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.KitchenItemGuid, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.KitchenDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.KitchenActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.KitchenDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.KitchenActivity, error)
}

type KitchenRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.KitchenDoc]
	repositories.SearchRepository[models.KitchenInfo]
	repositories.GuidRepository[models.KitchenItemGuid]
	repositories.ActivityRepository[models.KitchenActivity, models.KitchenDeleteActivity]
}

func NewKitchenRepository(pst microservice.IPersisterMongo) KitchenRepository {

	insRepo := KitchenRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.KitchenDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.KitchenInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.KitchenItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.KitchenActivity, models.KitchenDeleteActivity](pst)

	return insRepo
}
