package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/payment/bankmaster/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IBankMasterRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.BankMasterDoc) (string, error)
	CreateInBatch(docList []models.BankMasterDoc) error
	Update(shopID string, guid string, doc models.BankMasterDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.BankMasterInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.BankMasterDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.BankMasterItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.BankMasterDoc, error)

	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.BankMasterInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.BankMasterDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.BankMasterActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.BankMasterDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.BankMasterActivity, error)
}

type BankMasterRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.BankMasterDoc]
	repositories.SearchRepository[models.BankMasterInfo]
	repositories.GuidRepository[models.BankMasterItemGuid]
	repositories.ActivityRepository[models.BankMasterActivity, models.BankMasterDeleteActivity]
}

func NewBankMasterRepository(pst microservice.IPersisterMongo) *BankMasterRepository {

	insRepo := &BankMasterRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.BankMasterDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.BankMasterInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.BankMasterItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.BankMasterActivity, models.BankMasterDeleteActivity](pst)

	return insRepo
}
