package datatransfer

import (
	"context"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/internal/warehouse/models"
	warehouseRepository "smlcloudplatform/internal/warehouse/repositories"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/userplant/mongopagination"
)

type ProductWarehouseDataTransfer struct {
	transferConnection IDataTransferConnection
}

type IProductWarehouseDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.WarehouseDoc, mongopagination.PaginationData, error)
}

type ProductWarehouseDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.WarehouseDoc]
}

func NewProductWarehouseDataTransferRepository(mongodbPersister microservice.IPersisterMongo) IProductWarehouseDataTransferRepository {
	repo := &ProductWarehouseDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.WarehouseDoc](mongodbPersister)
	return repo
}

func (repo ProductWarehouseDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.WarehouseDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewProductWarehouseDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &ProductWarehouseDataTransfer{
		transferConnection: transferConnection,
	}
}

func (pdt *ProductWarehouseDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewProductWarehouseDataTransferRepository(pdt.transferConnection.GetSourceConnection())
	targetRepository := warehouseRepository.NewWarehouseRepository(pdt.transferConnection.GetTargetConnection())

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
