package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/debtaccount/debtorgroup/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IDebtorGroupRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.DebtorGroupDoc) (string, error)
	CreateInBatch(docList []models.DebtorGroupDoc) error
	Update(shopID string, guid string, doc models.DebtorGroupDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.DebtorGroupInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.DebtorGroupDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.DebtorGroupDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.DebtorGroupItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.DebtorGroupDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.DebtorGroupInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.DebtorGroupInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DebtorGroupDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DebtorGroupActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DebtorGroupDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DebtorGroupActivity, error)
}

type DebtorGroupRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DebtorGroupDoc]
	repositories.SearchRepository[models.DebtorGroupInfo]
	repositories.GuidRepository[models.DebtorGroupItemGuid]
	repositories.ActivityRepository[models.DebtorGroupActivity, models.DebtorGroupDeleteActivity]
}

func NewDebtorGroupRepository(pst microservice.IPersisterMongo) *DebtorGroupRepository {

	insRepo := &DebtorGroupRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DebtorGroupDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DebtorGroupInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.DebtorGroupItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.DebtorGroupActivity, models.DebtorGroupDeleteActivity](pst)

	return insRepo
}
