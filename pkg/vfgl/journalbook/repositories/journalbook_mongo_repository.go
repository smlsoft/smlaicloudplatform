package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/journalbook/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IJournalBookMongoRepository interface {
	Count(shopID string) (int, error)
	Create(category models.JournalBookDoc) (string, error)
	CreateInBatch(docList []models.JournalBookDoc) error
	Update(shopID string, guid string, category models.JournalBookDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters map[string]interface{}) (models.JournalBookDoc, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.JournalBookInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.JournalBookDoc, error)
}

type JournalBookMongoRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.JournalBookDoc]
	repositories.SearchRepository[models.JournalBookInfo]
	repositories.GuidRepository[models.JournalBookIdentifier]
}

func NewJournalBookMongoRepository(pst microservice.IPersisterMongo) JournalBookMongoRepository {

	insRepo := JournalBookMongoRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.JournalBookDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.JournalBookInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.JournalBookIdentifier](pst)

	return insRepo
}
