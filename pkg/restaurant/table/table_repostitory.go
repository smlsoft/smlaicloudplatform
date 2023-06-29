package table

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/table/models"
	"time"

	"github.com/userplant/mongopagination"
)

type ITableRepository interface {
	Create(category models.TableDoc) (string, error)
	CreateInBatch(docList []models.TableDoc) error
	Update(shopID string, guid string, category models.TableDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.TableInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.TableDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.TableItemGuid, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.TableDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.TableActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.TableDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.TableActivity, error)
}

type TableRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.TableDoc]
	repositories.SearchRepository[models.TableInfo]
	repositories.GuidRepository[models.TableItemGuid]
	repositories.ActivityRepository[models.TableActivity, models.TableDeleteActivity]
}

func NewTableRepository(pst microservice.IPersisterMongo) *TableRepository {
	insRepo := TableRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.TableDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.TableInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.TableItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.TableActivity, models.TableDeleteActivity](pst)

	return &insRepo
}
