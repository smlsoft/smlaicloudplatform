package repositories

import (
	"context"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/notify/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type INotifyRepository interface {
	Count(ctx context.Context, shopID string) (int, error)
	Create(ctx context.Context, doc models.NotifyDoc) (string, error)
	CreateInBatch(ctx context.Context, docList []models.NotifyDoc) error
	Update(ctx context.Context, shopID string, guid string, doc models.NotifyDoc) error
	DeleteByGuidfixed(sctx context.Context, hopID string, guid string, username string) error
	Delete(ctx context.Context, shopID string, username string, filters map[string]interface{}) error
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.NotifyInfo, mongopagination.PaginationData, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.NotifyDoc, error)
	FindByGuids(ctx context.Context, shopID string, guids []string) ([]models.NotifyDoc, error)

	FindByDocIndentityGuid(ctx context.Context, shopID string, indentityField string, indentityValue interface{}) (models.NotifyDoc, error)
	FindPageFilter(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.NotifyInfo, mongopagination.PaginationData, error)
	FindStep(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.NotifyInfo, int, error)

	FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.NotifyDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(sctx context.Context, hopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.NotifyActivity, mongopagination.PaginationData, error)
	FindDeletedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.NotifyDeleteActivity, error)
	FindCreatedOrUpdatedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.NotifyActivity, error)
}

type NotifyRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.NotifyDoc]
	repositories.SearchRepository[models.NotifyInfo]

	repositories.ActivityRepository[models.NotifyActivity, models.NotifyDeleteActivity]
}

func NewNotifyRepository(pst microservice.IPersisterMongo) *NotifyRepository {

	insRepo := &NotifyRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.NotifyDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.NotifyInfo](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.NotifyActivity, models.NotifyDeleteActivity](pst)

	return insRepo
}
