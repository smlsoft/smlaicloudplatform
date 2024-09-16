package datatransfer

import (
	"context"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/internal/slipimage/models"
	slipImageRepository "smlcloudplatform/internal/slipimage/repositories"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/userplant/mongopagination"
)

type SlipImageDataTransfer struct {
	transferConnection IDataTransferConnection
}

type ISlipImageDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SlipImageDoc, mongopagination.PaginationData, error)
}

type SlipImageDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.SlipImageDoc]
}

func NewSlipImageDataTransferRepository(mongodbPersister microservice.IPersisterMongo) ISlipImageDataTransferRepository {
	repo := &SlipImageDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.SlipImageDoc](mongodbPersister)
	return repo
}

func (repo SlipImageDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SlipImageDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewSlipImageDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &SlipImageDataTransfer{
		transferConnection: transferConnection,
	}
}

func (pdt *SlipImageDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewSlipImageDataTransferRepository(pdt.transferConnection.GetSourceConnection())
	targetRepository := slipImageRepository.NewSlipImageMongoRepository(pdt.transferConnection.GetTargetConnection())

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
