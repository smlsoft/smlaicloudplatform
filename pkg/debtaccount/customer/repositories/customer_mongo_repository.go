package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/debtaccount/customer/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type ICustomerRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.CustomerDoc) (string, error)
	CreateInBatch(docList []models.CustomerDoc) error
	Update(shopID string, guid string, doc models.CustomerDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.CustomerInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.CustomerDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.CustomerItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.CustomerDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.CustomerInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.CustomerInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CustomerDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CustomerActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CustomerDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CustomerActivity, error)
}

type CustomerRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.CustomerDoc]
	repositories.SearchRepository[models.CustomerInfo]
	repositories.GuidRepository[models.CustomerItemGuid]
	repositories.ActivityRepository[models.CustomerActivity, models.CustomerDeleteActivity]
}

func NewCustomerRepository(pst microservice.IPersisterMongo) *CustomerRepository {

	insRepo := &CustomerRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.CustomerDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.CustomerInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.CustomerItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.CustomerActivity, models.CustomerDeleteActivity](pst)

	return insRepo
}
