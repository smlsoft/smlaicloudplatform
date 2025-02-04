package debtortransaction

import (
	"smlaicloudplatform/internal/transaction/models"
	"smlaicloudplatform/internal/transaction/transactionconsumer/repositories"
	"smlaicloudplatform/pkg/microservice"
)

type IDebtorTransactionPGRepository interface {
	Get(shopID string, docNo string) (*models.DebtorTransactionPG, error)
	Create(doc models.DebtorTransactionPG) error
	Update(shopID string, docNo string, doc models.DebtorTransactionPG) error
	Delete(shopID string, docNo string, doc models.DebtorTransactionPG) error
}

type DebtorTransactionPGRepository struct {
	pst microservice.IPersister
	repositories.ITransactionConsumerRepository[models.DebtorTransactionPG]
}

func NewDebtorTransactionPGRepository(pst microservice.IPersister) IDebtorTransactionPGRepository {

	repo := &DebtorTransactionPGRepository{
		pst: pst,
	}

	repo.ITransactionConsumerRepository = repositories.NewTransactionConsumerRepository[models.DebtorTransactionPG](pst)
	return repo
}
