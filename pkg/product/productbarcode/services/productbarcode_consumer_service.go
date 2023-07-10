package services

import (
	"context"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/product/productbarcode/repositories"
	"smlcloudplatform/pkg/product/productbarcode/usecases"

	msModels "smlcloudplatform/internal/microservice/models"
	commonModels "smlcloudplatform/pkg/models"
)

type IProductBarcodeConsumeService interface {
	UpSert(shopID string, barcode string, doc models.ProductBarcodeDoc) (*models.ProductBarcodePg, error)
	Delete(ctx context.Context, shopID string, barcode string) error
	ReSync(shopID string) error
}

type ProductBarcodeConsumeService struct {
	productPgRepo         repositories.IProductBarcodePGRepository
	productMongoRepo      repositories.IProductBarcodeRepository
	productClickhouseRepo repositories.IProductBarcodeClickhouseRepository
	phaser                usecases.IProductBarcodePhaser
}

func NewProductBarcodeConsumerService(
	pst microservice.IPersister,
	mongoDBPersister microservice.IPersisterMongo,
	clickhousePersister microservice.IPersisterClickHouse,
	phaser usecases.IProductBarcodePhaser,
) IProductBarcodeConsumeService {

	productPgRepo := repositories.NewProductBarcodePGRepository(pst)
	productMongoRepo := repositories.NewProductBarcodeRepository(mongoDBPersister, nil)
	productClickhouseRepo := repositories.NewProductBarcodeClickhouseRepository(clickhousePersister)

	return &ProductBarcodeConsumeService{
		productPgRepo:         productPgRepo,
		productMongoRepo:      productMongoRepo,
		phaser:                phaser,
		productClickhouseRepo: productClickhouseRepo,
	}
}

func (svc ProductBarcodeConsumeService) UpSert(shopID string, barcode string, doc models.ProductBarcodeDoc) (*models.ProductBarcodePg, error) {

	pgDoc, err := svc.phaser.PhaseProductBarcodeDoc(&doc)
	if err != nil {
		return nil, err
	}

	findbarcodePG, err := svc.productPgRepo.Get(shopID, pgDoc.Barcode)
	if err != nil {
		return nil, err
	}

	if findbarcodePG != nil {
		err = svc.productPgRepo.Update(shopID, pgDoc.Barcode, pgDoc)
	} else {
		err = svc.productPgRepo.Create(pgDoc)
	}

	return nil, nil
}

func (svc ProductBarcodeConsumeService) Delete(ctx context.Context, shopID string, barcode string) error {

	err := svc.productPgRepo.Delete(shopID, barcode)
	if err != nil {
		return err
	}
	return nil
}
func (svc ProductBarcodeConsumeService) ReSync(shopID string) error {

	// resync 100

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
		barcodes, pages, err := svc.productMongoRepo.FindPage(context.Background(), shopID, nil, pageRequest)
		if err != nil {
			return err
		}

		for _, barcode := range barcodes {

			doc := models.ProductBarcodeDoc{
				ProductBarcodeData: models.ProductBarcodeData{
					ShopIdentity: commonModels.ShopIdentity{
						ShopID: shopID,
					},
					ProductBarcodeInfo: barcode,
				},
			}

			svc.UpSert(shopID, barcode.Barcode, doc)
		}

		if pages.TotalPage > int64(pageRequest.Page) {
			pageRequest.Page++
		} else {
			break
		}
	}

	return nil
}
