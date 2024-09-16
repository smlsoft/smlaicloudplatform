package datatransfer

import (
	"context"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/internal/transaction/saleinvoice/models"
	saleInvoiceRepository "smlcloudplatform/internal/transaction/saleinvoice/repositories"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/userplant/mongopagination"
)

type SaleInvoiceDataTransfer struct {
	transferConnection IDataTransferConnection
}

type ISaleInvoiceDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SaleInvoiceDoc, mongopagination.PaginationData, error)
}

type SaleInvoiceDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.SaleInvoiceDoc]
}

func NewSaleInvoiceDataTransferRepository(mongodbPersister microservice.IPersisterMongo) ISaleInvoiceDataTransferRepository {
	repo := &SaleInvoiceDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.SaleInvoiceDoc](mongodbPersister)
	return repo
}

func (repo SaleInvoiceDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SaleInvoiceDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewSaleInvoiceDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &SaleInvoiceDataTransfer{
		transferConnection: transferConnection,
	}
}

func (pdt *SaleInvoiceDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewSaleInvoiceDataTransferRepository(pdt.transferConnection.GetSourceConnection())
	targetRepository := saleInvoiceRepository.NewSaleInvoiceRepository(pdt.transferConnection.GetTargetConnection())

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
