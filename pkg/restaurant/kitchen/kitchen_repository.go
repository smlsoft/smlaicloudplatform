package kitchen

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/kitchen/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IKitchenRepository interface {
	Count(shopID string) (int, error)
	Create(category models.KitchenDoc) (string, error)
	CreateInBatch(docList []models.KitchenDoc) error
	Update(shopID string, guid string, category models.KitchenDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, authUsername string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.KitchenInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.KitchenDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.KitchenItemGuid, error)
	FindByDocIndentityGuid(shopID string, columnName string, filters interface{}) (models.KitchenDoc, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.KitchenInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.KitchenDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.KitchenActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.KitchenDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.KitchenActivity, error)
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
