package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/chartofaccount/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IChartOfAccountRepository interface {
	Count(shopID string) (int, error)
	Create(category models.ChartOfAccountDoc) (string, error)
	CreateInBatch(docList []models.ChartOfAccountDoc) error
	Update(shopID string, guid string, doc models.ChartOfAccountDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters interface{}) (models.ChartOfAccountDoc, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.ChartOfAccountInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ChartOfAccountDoc, error)
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
