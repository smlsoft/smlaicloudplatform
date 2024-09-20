package repositories

import (
	"context"
	"smlcloudplatform/internal/pos/shift/models"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	"time"

	"github.com/smlsoft/mongopagination"
)

type IShiftRepository interface {
	Count(ctx context.Context, shopID string) (int, error)
	Create(ctx context.Context, doc models.ShiftDoc) (string, error)
	CreateInBatch(ctx context.Context, docList []models.ShiftDoc) error
	Update(ctx context.Context, shopID string, guid string, doc models.ShiftDoc) error
	DeleteByGuidfixed(sctx context.Context, hopID string, guid string, username string) error
	Delete(ctx context.Context, shopID string, username string, filters map[string]interface{}) error
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ShiftInfo, mongopagination.PaginationData, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.ShiftDoc, error)
	FindByGuids(ctx context.Context, shopID string, guids []string) ([]models.ShiftDoc, error)

	FindInItemGuid(ctx context.Context, shopID string, columnName string, itemGuidList []string) ([]models.ShiftItemGuid, error)
	FindByDocIndentityGuid(ctx context.Context, shopID string, indentityField string, indentityValue interface{}) (models.ShiftDoc, error)
	FindPageFilter(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.ShiftInfo, mongopagination.PaginationData, error)
	FindStep(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.ShiftInfo, int, error)

	FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ShiftDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(sctx context.Context, hopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ShiftActivity, mongopagination.PaginationData, error)
	FindDeletedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ShiftDeleteActivity, error)
	FindCreatedOrUpdatedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ShiftActivity, error)
}

type ShiftRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ShiftDoc]
	repositories.SearchRepository[models.ShiftInfo]
	repositories.GuidRepository[models.ShiftItemGuid]
	repositories.ActivityRepository[models.ShiftActivity, models.ShiftDeleteActivity]
}

func NewShiftRepository(pst microservice.IPersisterMongo) *ShiftRepository {

	insRepo := &ShiftRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ShiftDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ShiftInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ShiftItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.ShiftActivity, models.ShiftDeleteActivity](pst)

	return insRepo
}
