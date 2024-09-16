package datatransfer

import (
	"context"
	"smlcloudplatform/internal/payment/bookbank/models"
	bookbankRepository "smlcloudplatform/internal/payment/bookbank/repositories"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/userplant/mongopagination"
)

type BookBankDataTransfer struct {
	transferConnection IDataTransferConnection
}

type IBookBankDataTransferDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.BookBankDoc, mongopagination.PaginationData, error)
}

type BookBankDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.BookBankDoc]
}

func NewBookBankDataTransferRepository(mongodbPersister microservice.IPersisterMongo) IBookBankDataTransferDataTransferRepository {
	repo := &BookBankDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.BookBankDoc](mongodbPersister)
	return repo
}

func (repo BookBankDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.BookBankDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewBookBankDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &BookBankDataTransfer{
		transferConnection: transferConnection,
	}
}

func (dt *BookBankDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewBookBankDataTransferRepository(dt.transferConnection.GetSourceConnection())
	targetRepository := bookbankRepository.NewBookBankRepository(dt.transferConnection.GetTargetConnection())

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
