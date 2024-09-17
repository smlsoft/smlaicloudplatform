package datatransfer

import (
	"context"
	"smlcloudplatform/internal/debtaccount/creditorgroup/models"
	creditorGroupRepositories "smlcloudplatform/internal/debtaccount/creditorgroup/repositories"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/pkg/microservice"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/userplant/mongopagination"
)

type CreditorGroupDataTransfer struct {
	transferConnection IDataTransferConnection
}

type ICreditorGroupDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable msModels.Pageable) ([]models.CreditorGroupDoc, mongopagination.PaginationData, error)
}

type CreditorGroupDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.CreditorGroupDoc]
}

func NewCreditorGroupDataTransferRepository(mongodbPersister microservice.IPersisterMongo) ICreditorGroupDataTransferRepository {
	repo := &CreditorGroupDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.CreditorGroupDoc](mongodbPersister)
	return repo
}

func (repo CreditorGroupDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable msModels.Pageable) ([]models.CreditorGroupDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewCreditorGroupDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {

	return &CreditorGroupDataTransfer{
		transferConnection: transferConnection,
	}
}

func (pdt *CreditorGroupDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewCreditorGroupDataTransferRepository(pdt.transferConnection.GetSourceConnection())
	targetRepository := creditorGroupRepositories.NewCreditorGroupRepository(pdt.transferConnection.GetTargetConnection())

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
