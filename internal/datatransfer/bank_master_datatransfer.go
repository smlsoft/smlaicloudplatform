package datatransfer

import (
	"context"
	"smlcloudplatform/internal/payment/bankmaster/models"
	bankRepository "smlcloudplatform/internal/payment/bankmaster/repositories"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/userplant/mongopagination"
)

type BankMasterDataTransfer struct {
	transferConnection IDataTransferConnection
}

type IBankMasterDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.BankMasterDoc, mongopagination.PaginationData, error)
}

type BankMasterDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.BankMasterDoc]
}

func NewBankMasterDataTransferRepository(mongodbPersister microservice.IPersisterMongo) IBankMasterDataTransferRepository {
	repo := &BankMasterDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.BankMasterDoc](mongodbPersister)
	return repo
}

func (repo BankMasterDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.BankMasterDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewBankMasterDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &BankMasterDataTransfer{
		transferConnection: transferConnection,
	}
}

func (dt *BankMasterDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewBankMasterDataTransferRepository(dt.transferConnection.GetSourceConnection())
	targetRepository := bankRepository.NewBankMasterRepository(dt.transferConnection.GetTargetConnection())

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
