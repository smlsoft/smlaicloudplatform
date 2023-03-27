package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/shop/employee/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IEmployeeRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.EmployeeDoc) (string, error)
	CreateInBatch(docList []models.EmployeeDoc) error
	Update(shopID string, guid string, doc models.EmployeeDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.EmployeeInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.EmployeeDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.EmployeeItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.EmployeeDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.EmployeeInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.EmployeeInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.EmployeeDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.EmployeeActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.EmployeeDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.EmployeeActivity, error)
}

type EmployeeRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.EmployeeDoc]
	repositories.SearchRepository[models.EmployeeInfo]
	repositories.GuidRepository[models.EmployeeItemGuid]
	repositories.ActivityRepository[models.EmployeeActivity, models.EmployeeDeleteActivity]
}

func NewEmployeeRepository(pst microservice.IPersisterMongo) *EmployeeRepository {

	insRepo := &EmployeeRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.EmployeeDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.EmployeeInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.EmployeeItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.EmployeeActivity, models.EmployeeDeleteActivity](pst)

	return insRepo
}
