package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/bankmaster/models"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IBankMasterRepository interface {
	Count(shopID string) (int, error)
	Create(category models.BankMasterDoc) (string, error)
	CreateInBatch(inventories []models.BankMasterDoc) error
	Update(shopID string, guid string, category models.BankMasterDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.BankMasterInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.BankMasterDoc, error)
}

type BankMasterRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.BankMasterDoc]
	repositories.SearchRepository[models.BankMasterInfo]
	repositories.GuidRepository[models.BankMasterItemGuid]
}

func NewBankMasterRepository(pst microservice.IPersisterMongo) BankMasterRepository {

	insRepo := BankMasterRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.BankMasterDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.BankMasterInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.BankMasterItemGuid](pst)

	return insRepo
}
