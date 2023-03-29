package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/debtaccount/creditorgroup/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type ICreditorGroupRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.CreditorGroupDoc) (string, error)
	CreateInBatch(docList []models.CreditorGroupDoc) error
	Update(shopID string, guid string, doc models.CreditorGroupDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.CreditorGroupInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.CreditorGroupDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.CreditorGroupDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.CreditorGroupItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.CreditorGroupDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.CreditorGroupInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.CreditorGroupInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CreditorGroupDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CreditorGroupActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CreditorGroupDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CreditorGroupActivity, error)
}

type CreditorGroupRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.CreditorGroupDoc]
	repositories.SearchRepository[models.CreditorGroupInfo]
	repositories.GuidRepository[models.CreditorGroupItemGuid]
	repositories.ActivityRepository[models.CreditorGroupActivity, models.CreditorGroupDeleteActivity]
}

func NewCreditorGroupRepository(pst microservice.IPersisterMongo) *CreditorGroupRepository {

	insRepo := &CreditorGroupRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.CreditorGroupDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.CreditorGroupInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.CreditorGroupItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.CreditorGroupActivity, models.CreditorGroupDeleteActivity](pst)

	return insRepo
}
