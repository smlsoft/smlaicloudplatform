package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IChartOfAccountRepository interface {
	Count(shopID string) (int, error)
	Create(category vfgl.ChartOfAccountDoc) (string, error)
	CreateInBatch(docList []vfgl.ChartOfAccountDoc) error
	Update(shopID string, guid string, doc vfgl.ChartOfAccountDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]vfgl.ChartOfAccountInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (vfgl.ChartOfAccountDoc, error)
}

type ChartOfAccountRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[vfgl.ChartOfAccountDoc]
	repositories.SearchRepository[vfgl.ChartOfAccountInfo]
	repositories.GuidRepository[vfgl.ChartOfAccountIndentityId]
}

func NewChartOfAccountRepository(pst microservice.IPersisterMongo) ChartOfAccountRepository {
	repo := ChartOfAccountRepository{
		pst: pst,
	}

	repo.CrudRepository = repositories.NewCrudRepository[vfgl.ChartOfAccountDoc](pst)
	repo.SearchRepository = repositories.NewSearchRepository[vfgl.ChartOfAccountInfo](pst)
	repo.GuidRepository = repositories.NewGuidRepository[vfgl.ChartOfAccountIndentityId](pst)
	return repo
}
