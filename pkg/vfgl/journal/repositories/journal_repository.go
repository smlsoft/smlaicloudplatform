package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IJournalRepository interface {
	Count(shopID string) (int, error)
	Create(category vfgl.JournalDoc) (string, error)
	CreateInBatch(inventories []vfgl.JournalDoc) error
	Update(shopID string, guid string, category vfgl.JournalDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]vfgl.JournalInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (vfgl.JournalDoc, error)
}

type JournalRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[vfgl.JournalDoc]
	repositories.SearchRepository[vfgl.JournalInfo]
	repositories.GuidRepository[vfgl.JournalItemGuid]
}

func NewJournalRepository(pst microservice.IPersisterMongo) JournalRepository {

	insRepo := JournalRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[vfgl.JournalDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[vfgl.JournalInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[vfgl.JournalItemGuid](pst)

	return insRepo
}
