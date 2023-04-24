package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/productsection/sectiondepartment/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type ISectionDepartmentRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.SectionDepartmentDoc) (string, error)
	CreateInBatch(docList []models.SectionDepartmentDoc) error
	Update(shopID string, guid string, doc models.SectionDepartmentDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SectionDepartmentInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SectionDepartmentDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.SectionDepartmentItemGuid, error)
	FindOneFilter(shopID string, filters map[string]interface{}) (models.SectionDepartmentDoc, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.SectionDepartmentDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.SectionDepartmentInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.SectionDepartmentInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionDepartmentDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionDepartmentActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SectionDepartmentDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SectionDepartmentActivity, error)
}

type SectionDepartmentRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SectionDepartmentDoc]
	repositories.SearchRepository[models.SectionDepartmentInfo]
	repositories.GuidRepository[models.SectionDepartmentItemGuid]
	repositories.ActivityRepository[models.SectionDepartmentActivity, models.SectionDepartmentDeleteActivity]
}

func NewSectionDepartmentRepository(pst microservice.IPersisterMongo) *SectionDepartmentRepository {

	insRepo := &SectionDepartmentRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SectionDepartmentDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SectionDepartmentInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SectionDepartmentItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.SectionDepartmentActivity, models.SectionDepartmentDeleteActivity](pst)

	return insRepo
}
