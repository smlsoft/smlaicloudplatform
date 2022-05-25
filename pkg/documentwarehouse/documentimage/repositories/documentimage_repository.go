package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IDocumentImageRepository interface {
	Create(doc models.DocumentImageDoc) (string, error)
	Update(shopID string, guid string, doc models.DocumentImageDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters map[string]interface{}) (models.DocumentImageDoc, error)
	FindByGuid(shopID string, guid string) (models.DocumentImageDoc, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)
}

type DocumentImageRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DocumentImageDoc]
	repositories.SearchRepository[models.DocumentImageInfo]
}

func NewDocumentImageRepository(pst microservice.IPersisterMongo) DocumentImageRepository {
	insRepo := DocumentImageRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DocumentImageDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DocumentImageInfo](pst)

	return insRepo
}
