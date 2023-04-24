package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/productsection/sectionbusinesstype/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type ISectionBusinessTypeRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.SectionBusinessTypeDoc) (string, error)
	CreateInBatch(docList []models.SectionBusinessTypeDoc) error
	Update(shopID string, guid string, doc models.SectionBusinessTypeDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SectionBusinessTypeInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SectionBusinessTypeDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.SectionBusinessTypeItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.SectionBusinessTypeDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.SectionBusinessTypeInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.SectionBusinessTypeInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBusinessTypeDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBusinessTypeActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SectionBusinessTypeDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SectionBusinessTypeActivity, error)
}

type SectionBusinessTypeRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SectionBusinessTypeDoc]
	repositories.SearchRepository[models.SectionBusinessTypeInfo]
	repositories.GuidRepository[models.SectionBusinessTypeItemGuid]
	repositories.ActivityRepository[models.SectionBusinessTypeActivity, models.SectionBusinessTypeDeleteActivity]
}

func NewSectionBusinessTypeRepository(pst microservice.IPersisterMongo) *SectionBusinessTypeRepository {

	insRepo := &SectionBusinessTypeRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SectionBusinessTypeDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SectionBusinessTypeInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SectionBusinessTypeItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.SectionBusinessTypeActivity, models.SectionBusinessTypeDeleteActivity](pst)

	return insRepo
}
