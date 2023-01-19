package shoptable

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/shoptable/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IShopTableRepository interface {
	Count(shopID string) (int, error)
	Create(category models.ShopTableDoc) (string, error)
	CreateInBatch(docList []models.ShopTableDoc) error
	Update(shopID string, guid string, category models.ShopTableDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ShopTableInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ShopTableDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ShopTableItemGuid, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.ShopTableDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.ShopTableActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.ShopTableDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.ShopTableActivity, error)
}

type ShopTableRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ShopTableDoc]
	repositories.SearchRepository[models.ShopTableInfo]
	repositories.GuidRepository[models.ShopTableItemGuid]
	repositories.ActivityRepository[models.ShopTableActivity, models.ShopTableDeleteActivity]
}

func NewShopTableRepository(pst microservice.IPersisterMongo) ShopTableRepository {
	insRepo := ShopTableRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ShopTableDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ShopTableInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ShopTableItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.ShopTableActivity, models.ShopTableDeleteActivity](pst)

	return insRepo
}
