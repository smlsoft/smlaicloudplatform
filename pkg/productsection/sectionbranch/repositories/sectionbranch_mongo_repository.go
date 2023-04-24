package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/productsection/sectionbranch/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type ISectionBranchRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.SectionBranchDoc) (string, error)
	CreateInBatch(docList []models.SectionBranchDoc) error
	Update(shopID string, guid string, doc models.SectionBranchDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SectionBranchInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SectionBranchDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.SectionBranchItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.SectionBranchDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.SectionBranchInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.SectionBranchInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBranchDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBranchActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SectionBranchDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SectionBranchActivity, error)
}

type SectionBranchRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SectionBranchDoc]
	repositories.SearchRepository[models.SectionBranchInfo]
	repositories.GuidRepository[models.SectionBranchItemGuid]
	repositories.ActivityRepository[models.SectionBranchActivity, models.SectionBranchDeleteActivity]
}

func NewSectionBranchRepository(pst microservice.IPersisterMongo) *SectionBranchRepository {

	insRepo := &SectionBranchRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SectionBranchDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SectionBranchInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SectionBranchItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.SectionBranchActivity, models.SectionBranchDeleteActivity](pst)

	return insRepo
}
