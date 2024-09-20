package datatransfer

import (
	"context"
	"smlcloudplatform/internal/product/productbarcode/models"
	productbarcoderepository "smlcloudplatform/internal/product/productbarcode/repositories"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/pkg/microservice"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/smlsoft/mongopagination"
)

type ProductBarcodeDataTransfer struct {
	transferConnection IDataTransferConnection
}

type IProductBarcodeDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable msModels.Pageable) ([]models.ProductBarcodeDoc, mongopagination.PaginationData, error)
}

type ProductBarcodeDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.ProductBarcodeDoc]
}

func NewProductBarcodeDataTransferRepository(mongodbPersister microservice.IPersisterMongo) IProductBarcodeDataTransferRepository {

	repo := &ProductBarcodeDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.ProductBarcodeDoc](mongodbPersister)
	return repo
}

func (repo ProductBarcodeDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable msModels.Pageable) ([]models.ProductBarcodeDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewProductBarcodeDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &ProductBarcodeDataTransfer{
		transferConnection: transferConnection,
	}
}
func (pbd *ProductBarcodeDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceProductBarcodeRepository := NewProductBarcodeDataTransferRepository(pbd.transferConnection.GetSourceConnection())
	targetProductBarcodeRepository := productbarcoderepository.NewProductBarcodeRepository(pbd.transferConnection.GetTargetConnection(), nil)

	pageRequest := msModels.Pageable{
		Limit: 100,
		Page:  1,
		Sorts: []msModels.KeyInt{
			{
				Key:   "guidfixed",
				Value: -1,
			},
		},
	}

	for {
		barcodes, pages, err := sourceProductBarcodeRepository.FindPage(ctx, shopID, nil, pageRequest)
		if err != nil {
			return err
		}

		if len(barcodes) > 0 {
			err = targetProductBarcodeRepository.CreateInBatch(ctx, barcodes)
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
