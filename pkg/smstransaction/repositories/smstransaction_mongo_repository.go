package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/smstransaction/models"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type ISmsTransactionRepository interface {
	Count(shopID string) (int, error)
	Create(category models.SmsTransactionDoc) (string, error)
	CreateInBatch(inventories []models.SmsTransactionDoc) error
	Update(shopID string, guid string, category models.SmsTransactionDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.SmsTransactionInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SmsTransactionDoc, error)
}

type SmsTransactionRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SmsTransactionDoc]
	repositories.SearchRepository[models.SmsTransactionInfo]
	repositories.GuidRepository[models.SmsTransactionItemGuid]
}

func NewSmsTransactionRepository(pst microservice.IPersisterMongo) SmsTransactionRepository {

	insRepo := SmsTransactionRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SmsTransactionDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SmsTransactionInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SmsTransactionItemGuid](pst)

	return insRepo
}
