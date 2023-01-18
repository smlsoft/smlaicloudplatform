package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/accountgroup/models"

	"github.com/userplant/mongopagination"
)

type IAccountGroupMongoRepository interface {
	Count(shopID string) (int, error)
	Create(category models.AccountGroupDoc) (string, error)
	CreateInBatch(docList []models.AccountGroupDoc) error
	Update(shopID string, guid string, category models.AccountGroupDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters interface{}) (models.AccountGroupDoc, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.AccountGroupInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.AccountGroupDoc, error)
}

type AccountGroupMongoRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.AccountGroupDoc]
	repositories.SearchRepository[models.AccountGroupInfo]
	repositories.GuidRepository[models.AccountGroupIdentifier]
}

func NewAccountGroupMongoRepository(pst microservice.IPersisterMongo) AccountGroupMongoRepository {

	insRepo := AccountGroupMongoRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.AccountGroupDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.AccountGroupInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.AccountGroupIdentifier](pst)

	return insRepo
}
