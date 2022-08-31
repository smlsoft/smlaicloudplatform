package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/paymentmaster/models"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IPaymentMasterRepository interface {
	Count(shopID string) (int, error)
	Create(category models.PaymentMasterDoc) (string, error)
	CreateInBatch(docList []models.PaymentMasterDoc) error
	Update(shopID string, guid string, category models.PaymentMasterDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Find(shopID string, colNameSearch []string, q string) ([]models.PaymentMasterInfo, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.PaymentMasterInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.PaymentMasterDoc, error)
}

type PaymentMasterRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.PaymentMasterDoc]
	repositories.SearchRepository[models.PaymentMasterInfo]
	repositories.GuidRepository[models.PaymentMasterItemGuid]
}

func NewPaymentMasterRepository(pst microservice.IPersisterMongo) PaymentMasterRepository {

	insRepo := PaymentMasterRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.PaymentMasterDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.PaymentMasterInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.PaymentMasterItemGuid](pst)

	return insRepo
}
