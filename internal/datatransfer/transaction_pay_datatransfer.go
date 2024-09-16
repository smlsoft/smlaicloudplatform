package datatransfer

import (
	"context"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/internal/transaction/pay/models"
	transactionPayRepository "smlcloudplatform/internal/transaction/pay/repositories"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/userplant/mongopagination"
)

type TransactionPayDataTransfer struct {
	transferConnection IDataTransferConnection
}

type ITransactionPayDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.PayDoc, mongopagination.PaginationData, error)
}

type TransactionPayDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.PayDoc]
}

func NewTransactionPayDataTransferRepository(mongodbPersister microservice.IPersisterMongo) ITransactionPayDataTransferRepository {
	repo := &TransactionPayDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.PayDoc](mongodbPersister)
	return repo
}

func (repo TransactionPayDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.PayDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewTransactionPayDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &TransactionPayDataTransfer{
		transferConnection: transferConnection,
	}
}

func (pdt *TransactionPayDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewTransactionPayDataTransferRepository(pdt.transferConnection.GetSourceConnection())
	targetRepository := transactionPayRepository.NewPayRepository(pdt.transferConnection.GetTargetConnection())

	pageRequest := msModels.Pageable{
		Limit: 100,
		Page:  1,
	}

	for {

		docs, pages, err := sourceRepository.FindPage(ctx, shopID, nil, pageRequest)
		if err != nil {
			return err
		}

		if len(docs) > 0 {
			err = targetRepository.CreateInBatch(ctx, docs)
			if err != nil {
				return err
			}
		}

		if pages.TotalPage > int64(pageRequest.Page) {
			pageRequest.Page++
		} else {
			break
		}
	}

	return nil
}
