package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/debtaccount/creditor/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type ICreditorRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.CreditorDoc) (string, error)
	CreateInBatch(docList []models.CreditorDoc) error
	Update(shopID string, guid string, doc models.CreditorDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.CreditorInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.CreditorDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.CreditorItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.CreditorDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.CreditorInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.CreditorInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CreditorDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CreditorActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CreditorDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CreditorActivity, error)
}

type CreditorRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.CreditorDoc]
	repositories.SearchRepository[models.CreditorInfo]
	repositories.GuidRepository[models.CreditorItemGuid]
	repositories.ActivityRepository[models.CreditorActivity, models.CreditorDeleteActivity]
}

func NewCreditorRepository(pst microservice.IPersisterMongo) *CreditorRepository {

	insRepo := &CreditorRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.CreditorDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.CreditorInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.CreditorItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.CreditorActivity, models.CreditorDeleteActivity](pst)

	return insRepo
}
