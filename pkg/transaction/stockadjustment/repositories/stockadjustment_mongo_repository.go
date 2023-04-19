package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stockadjustment/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IStockAdjustmentRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.StockAdjustmentDoc) (string, error)
	CreateInBatch(docList []models.StockAdjustmentDoc) error
	Update(shopID string, guid string, doc models.StockAdjustmentDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.StockAdjustmentDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.StockAdjustmentItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.StockAdjustmentDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.StockAdjustmentInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockAdjustmentDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockAdjustmentActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockAdjustmentDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockAdjustmentActivity, error)
}

type StockAdjustmentRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.StockAdjustmentDoc]
	repositories.SearchRepository[models.StockAdjustmentInfo]
	repositories.GuidRepository[models.StockAdjustmentItemGuid]
	repositories.ActivityRepository[models.StockAdjustmentActivity, models.StockAdjustmentDeleteActivity]
}

func NewStockAdjustmentRepository(pst microservice.IPersisterMongo) *StockAdjustmentRepository {

	insRepo := &StockAdjustmentRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.StockAdjustmentDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.StockAdjustmentInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.StockAdjustmentItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.StockAdjustmentActivity, models.StockAdjustmentDeleteActivity](pst)

	return insRepo
}
