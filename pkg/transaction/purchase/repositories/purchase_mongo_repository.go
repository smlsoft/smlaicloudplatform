package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/purchase/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IPurchaseRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.PurchaseDoc) (string, error)
	CreateInBatch(docList []models.PurchaseDoc) error
	Update(shopID string, guid string, doc models.PurchaseDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.PurchaseInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.PurchaseDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.PurchaseItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.PurchaseDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.PurchaseInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.PurchaseInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PurchaseDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PurchaseActivity, error)
}

type PurchaseRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.PurchaseDoc]
	repositories.SearchRepository[models.PurchaseInfo]
	repositories.GuidRepository[models.PurchaseItemGuid]
	repositories.ActivityRepository[models.PurchaseActivity, models.PurchaseDeleteActivity]
}

func NewPurchaseRepository(pst microservice.IPersisterMongo) *PurchaseRepository {

	insRepo := &PurchaseRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.PurchaseDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.PurchaseInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.PurchaseItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.PurchaseActivity, models.PurchaseDeleteActivity](pst)

	return insRepo
}
