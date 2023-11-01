package table

import (
	"context"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/table/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ITableRepository interface {
	Count(ctx context.Context, shopID string) (int, error)
	Create(ctx context.Context, category models.TableDoc) (string, error)
	CreateInBatch(ctx context.Context, docList []models.TableDoc) error
	Update(ctx context.Context, shopID string, guid string, category models.TableDoc) error
	DeleteByGuidfixed(ctx context.Context, shopID string, guid string, username string) error
	Delete(ctx context.Context, shopID string, authUsername string, filters map[string]interface{}) error
	FindPageFilter(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.TableInfo, mongopagination.PaginationData, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.TableDoc, error)
	FindInItemGuid(ctx context.Context, shopID string, columnName string, itemGuidList []string) ([]models.TableItemGuid, error)
	FindByDocIndentityGuid(ctx context.Context, shopID string, columnName string, filters interface{}) (models.TableDoc, error)
	FindStep(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.TableInfo, int, error)
	SaveXOrder(ctx context.Context, shopID string, guid string, xorder uint) error

	FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.TableDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.TableActivity, mongopagination.PaginationData, error)
	FindDeletedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.TableDeleteActivity, error)
	FindCreatedOrUpdatedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.TableActivity, error)
}

type TableRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.TableDoc]
	repositories.SearchRepository[models.TableInfo]
	repositories.GuidRepository[models.TableItemGuid]
	repositories.ActivityRepository[models.TableActivity, models.TableDeleteActivity]
}

func NewTableRepository(pst microservice.IPersisterMongo) *TableRepository {
	insRepo := TableRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.TableDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.TableInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.TableItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.TableActivity, models.TableDeleteActivity](pst)

	return &insRepo
}

func (repo TableRepository) SaveXOrder(ctx context.Context, shopID string, guid string, xorder uint) error {

	filters := bson.M{
		"shopid":    shopID,
		"guidfixed": guid,
	}

	return repo.pst.Update(ctx, models.TableDoc{}, filters, bson.M{"$set": bson.M{"xorder": xorder}})
}
