package repositories

import (
	"context"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/chartofaccount/models"

	"github.com/userplant/mongopagination"
)

type IChartOfAccountRepository interface {
	Count(ctx context.Context, shopID string) (int, error)
	Create(category models.ChartOfAccountDoc) (string, error)
	CreateInBatch(ctx context.Context, docList []models.ChartOfAccountDoc) error
	Update(ctx context.Context, shopID string, guid string, doc models.ChartOfAccountDoc) error
	DeleteByGuidfixed(ctx context.Context, shopID string, guid string, username string) error
	FindOne(shopID string, filters interface{}) (models.ChartOfAccountDoc, error)
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ChartOfAccountInfo, mongopagination.PaginationData, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.ChartOfAccountDoc, error)
}

type ChartOfAccountRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ChartOfAccountDoc]
	repositories.SearchRepository[models.ChartOfAccountInfo]
	repositories.GuidRepository[models.ChartOfAccountIndentityId]
}

func NewChartOfAccountRepository(pst microservice.IPersisterMongo) ChartOfAccountRepository {
	repo := ChartOfAccountRepository{
		pst: pst,
	}

	repo.CrudRepository = repositories.NewCrudRepository[models.ChartOfAccountDoc](pst)
	repo.SearchRepository = repositories.NewSearchRepository[models.ChartOfAccountInfo](pst)
	repo.GuidRepository = repositories.NewGuidRepository[models.ChartOfAccountIndentityId](pst)
	return repo
}
