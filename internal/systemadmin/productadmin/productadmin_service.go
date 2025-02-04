package productadmin

import (
	"context"
	productBarcodeRepositories "smlaicloudplatform/internal/product/productbarcode/repositories"
	stockProcessModels "smlaicloudplatform/internal/stockprocess/models"
	stockProcessRepository "smlaicloudplatform/internal/stockprocess/repositories"
	"smlaicloudplatform/pkg/microservice"
	msModels "smlaicloudplatform/pkg/microservice/models"
	"time"
)

type IProductAdminService interface {
	ReSyncProductBarcode(shopID string) error
	ReCalcStockBalance(shopID string) error
	DeleteProductBarcodeAll(shopID string, userName string) error
}

type ProductAdminService struct {
	kafkaRepo          productBarcodeRepositories.IProductBarcodeMessageQueueRepository
	mongoRepo          IProductAdminMongoRepository
	stockProcessMGRepo stockProcessRepository.IStockProcessMessageQueueRepository
	timeoutDuration    time.Duration
}

func NewProductAdminService(pst microservice.IPersisterMongo, kfProducer microservice.IProducer) IProductAdminService {
	kafkaRepo := productBarcodeRepositories.NewProductBarcodeMessageQueueRepository(kfProducer)
	mongoRepo := NewProductAdminMongoRepository(pst)

	stockProcessMGRepo := stockProcessRepository.NewStockProcessMessageQueueRepository(kfProducer)

	return &ProductAdminService{
		kafkaRepo:          kafkaRepo,
		mongoRepo:          mongoRepo,
		stockProcessMGRepo: stockProcessMGRepo,
		timeoutDuration:    time.Duration(30) * time.Second,
	}
}

func (svc ProductAdminService) ReSyncProductBarcode(shopID string) error {

	// find product barcode by shopid
	ctx, cancel := context.WithTimeout(context.Background(), svc.timeoutDuration)
	defer cancel()

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
		barcodes, pages, err := svc.mongoRepo.FindPage(ctx, shopID, nil, pageRequest)
		if err != nil {
			return err
		}

		// barcodeDocs := []models.ProductBarcodeDoc{}
		// for _, barcode := range barcodes {

		// 	barcodeDocs = append(barcodeDocs, barcode)
		// }

		// Send Message to MQ
		err = svc.kafkaRepo.CreateInBatch(barcodes)
		if err != nil {
			return err
		}

		if pages.TotalPage > int64(pageRequest.Page) {
			pageRequest.Page++
		} else {
			break
		}
	}

	return nil

}

func (svc ProductAdminService) ReCalcStockBalance(shopID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), svc.timeoutDuration)
	defer cancel()

	// find productbarocode by shopid
	barcodes, err := svc.mongoRepo.FindProductBarcodeByShopId(ctx, shopID)
	if err != nil {
		return err
	}

	var requestStockProcessLists []stockProcessModels.StockProcessRequest

	if len(barcodes) > 0 {
		for _, item := range barcodes {

			if item.Barcode != "" {
				requestStockProcessLists = append(requestStockProcessLists, stockProcessModels.StockProcessRequest{
					ShopID:  shopID,
					Barcode: item.Barcode,
				})
			}
		}
	}

	err = svc.stockProcessMGRepo.CreateInBatch(requestStockProcessLists)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductAdminService) DeleteProductBarcodeAll(shopID string, userName string) error {

	ctx, cancel := context.WithTimeout(context.Background(), svc.timeoutDuration)
	defer cancel()

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
		barcodes, pages, err := svc.mongoRepo.FindPage(ctx, shopID, nil, pageRequest)
		if err != nil {
			return err
		}

		ids := []string{}
		for _, barcode := range barcodes {

			ids = append(ids, barcode.ID.Hex())
		}

		err = svc.mongoRepo.DeleteProductBarcodeByShopId(ctx, shopID, userName, ids)
		if err != nil {
			return err
		}

		// Send Message to MQ
		err = svc.kafkaRepo.DeleteInBatch(barcodes)
		if err != nil {
			return err
		}

		if pages.TotalPage > int64(pageRequest.Page) {
			pageRequest.Page++
		} else {
			break
		}
	}

	return nil
}
