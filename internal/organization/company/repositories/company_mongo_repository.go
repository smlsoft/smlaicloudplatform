package repositories

import (
	"context"
	"smlaicloudplatform/internal/organization/company/models"
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/pkg/microservice"
	micromodels "smlaicloudplatform/pkg/microservice/models"
	"time"

	"github.com/smlsoft/mongopagination"
)

type ICompanyRepository interface {
	Count(ctx context.Context, shopID string) (int, error)
	Create(ctx context.Context, doc models.CompanyDoc) (string, error)
	CreateInBatch(ctx context.Context, docList []models.CompanyDoc) error
	Update(ctx context.Context, shopID string, guid string, doc models.CompanyDoc) error
	DeleteByGuidfixed(ctx context.Context, shopID string, guid string, username string) error
	Delete(ctx context.Context, shopID string, username string, filters map[string]interface{}) error
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.CompanyInfo, mongopagination.PaginationData, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.CompanyDoc, error)

	FindInItemGuid(ctx context.Context, shopID string, columnName string, itemGuidList []string) ([]models.CompanyItemGuid, error)
	FindInItemGuids(ctx context.Context, shopID string, columnName string, itemGuidList []interface{}) ([]models.CompanyItemGuid, error)
	FindByDocIndentityGuid(ctx context.Context, shopID string, indentityField string, indentityValue interface{}) (models.CompanyDoc, error)
	FindPageFilter(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.CompanyInfo, mongopagination.PaginationData, error)
	FindStep(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.CompanyInfo, int, error)
	FindOneFilter(ctx context.Context, shopID string, filters map[string]interface{}) (models.CompanyDoc, error)

	FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CompanyDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CompanyActivity, mongopagination.PaginationData, error)
	FindDeletedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CompanyDeleteActivity, error)
	FindCreatedOrUpdatedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CompanyActivity, error)
}

type CompanyRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.CompanyDoc]
	repositories.SearchRepository[models.CompanyInfo]
	repositories.GuidRepository[models.CompanyItemGuid]
	repositories.ActivityRepository[models.CompanyActivity, models.CompanyDeleteActivity]
}

func NewCompanyRepository(pst microservice.IPersisterMongo) *CompanyRepository {

	insRepo := &CompanyRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.CompanyDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.CompanyInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.CompanyItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.CompanyActivity, models.CompanyDeleteActivity](pst)

	return insRepo
}
