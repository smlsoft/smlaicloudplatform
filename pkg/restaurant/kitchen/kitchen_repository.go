package kitchen

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/kitchen/models"

	"github.com/userplant/mongopagination"
)

type IKitchenRepository interface {
	Count(shopID string) (int, error)
	Create(category models.KitchenDoc) (string, error)
	CreateInBatch(docList []models.KitchenDoc) error
	Update(shopID string, guid string, category models.KitchenDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.KitchenInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.KitchenDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.KitchenItemGuid, error)

	// FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.KitchenDeleteActivity, mongopagination.PaginationData, error)
	// FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.KitchenActivity, mongopagination.PaginationData, error)
	// FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.KitchenDeleteActivity, error)
	// FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.KitchenActivity, error)
	// FindStep(shopID string, searchInFields []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.KitchenInfo, int, error)
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
