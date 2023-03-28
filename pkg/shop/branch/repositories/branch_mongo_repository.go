package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/shop/branch/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IBranchRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.BranchDoc) (string, error)
	CreateInBatch(docList []models.BranchDoc) error
	Update(shopID string, guid string, doc models.BranchDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.BranchInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.BranchDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.BranchItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.BranchDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.BranchInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.BranchInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.BranchDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.BranchActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.BranchDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.BranchActivity, error)
}

type BranchRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.BranchDoc]
	repositories.SearchRepository[models.BranchInfo]
	repositories.GuidRepository[models.BranchItemGuid]
	repositories.ActivityRepository[models.BranchActivity, models.BranchDeleteActivity]
}

func NewBranchRepository(pst microservice.IPersisterMongo) *BranchRepository {

	insRepo := &BranchRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.BranchDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.BranchInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.BranchItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.BranchActivity, models.BranchDeleteActivity](pst)

	return insRepo
}
