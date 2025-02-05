package repositories

import (
	"context"
	"smlaicloudplatform/internal/organization/language/models"
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/pkg/microservice"
	micromodels "smlaicloudplatform/pkg/microservice/models"
	"time"

	"github.com/smlsoft/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ILanguageRepository interface {
	Count(ctx context.Context, shopID string) (int, error)
	Create(ctx context.Context, doc models.LanguageDoc) (string, error)
	CreateInBatch(ctx context.Context, docList []models.LanguageDoc) error
	Update(ctx context.Context, shopID string, guid string, doc models.LanguageDoc) error
	DeleteByGuidfixed(ctx context.Context, shopID string, guid string, username string) error
	Delete(ctx context.Context, shopID string, username string, filters map[string]interface{}) error
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.LanguageInfo, mongopagination.PaginationData, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.LanguageDoc, error)

	FindInItemGuid(ctx context.Context, shopID string, columnName string, itemGuidList []string) ([]models.LanguageItemGuid, error)
	FindByDocIndentityGuid(ctx context.Context, shopID string, indentityField string, indentityValue interface{}) (models.LanguageDoc, error)
	FindPageFilter(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.LanguageInfo, mongopagination.PaginationData, error)
	FindStep(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.LanguageInfo, int, error)

	FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.LanguageDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.LanguageActivity, mongopagination.PaginationData, error)
	FindDeletedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.LanguageDeleteActivity, error)
	FindCreatedOrUpdatedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.LanguageActivity, error)

	FindOneByCode(ctx context.Context, shopID, code string) (models.LanguageDoc, error)
}

type LanguageRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.LanguageDoc]
	repositories.SearchRepository[models.LanguageInfo]
	repositories.GuidRepository[models.LanguageItemGuid]
	repositories.ActivityRepository[models.LanguageActivity, models.LanguageDeleteActivity]
}

func NewLanguageRepository(pst microservice.IPersisterMongo) *LanguageRepository {

	insRepo := &LanguageRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.LanguageDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.LanguageInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.LanguageItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.LanguageActivity, models.LanguageDeleteActivity](pst)

	return insRepo
}

func (repo LanguageRepository) FindOneByCode(ctx context.Context, shopID string, code string) (models.LanguageDoc, error) {
	doc := models.LanguageDoc{}
	err := repo.pst.FindOne(ctx,
		models.LanguageDoc{},
		bson.M{
			"shopid":    shopID,
			"deletedat": bson.M{"$exists": false},
			"code":      code,
		}, &doc)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
