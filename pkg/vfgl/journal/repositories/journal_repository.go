package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/journal/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IJournalRepository interface {
	Count(shopID string) (int, error)
	Create(category models.JournalDoc) (string, error)
	CreateInBatch(inventories []models.JournalDoc) error
	Update(shopID string, guid string, category models.JournalDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.JournalInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.JournalDoc, error)
}

type JournalRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.JournalDoc]
	repositories.SearchRepository[models.JournalInfo]
	repositories.GuidRepository[models.JournalItemGuid]
}

func NewJournalRepository(pst microservice.IPersisterMongo) JournalRepository {

	insRepo := JournalRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.JournalDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.JournalInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.JournalItemGuid](pst)

	return insRepo
}
