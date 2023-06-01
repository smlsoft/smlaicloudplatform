package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/debtaccount/debtor/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IDebtorRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.DebtorDoc) (string, error)
	CreateInBatch(docList []models.DebtorDoc) error
	Update(shopID string, guid string, doc models.DebtorDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.DebtorInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.DebtorDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.DebtorItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.DebtorDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.DebtorInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.DebtorInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DebtorDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DebtorActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DebtorDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DebtorActivity, error)
}

type DebtorRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DebtorDoc]
	repositories.SearchRepository[models.DebtorInfo]
	repositories.GuidRepository[models.DebtorItemGuid]
	repositories.ActivityRepository[models.DebtorActivity, models.DebtorDeleteActivity]
}

func NewDebtorRepository(pst microservice.IPersisterMongo) *DebtorRepository {

	insRepo := &DebtorRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DebtorDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DebtorInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.DebtorItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.DebtorActivity, models.DebtorDeleteActivity](pst)

	return insRepo
}
