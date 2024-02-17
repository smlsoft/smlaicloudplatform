package saleinvoice

import (
	"context"
	saleInvoiceRepository "smlcloudplatform/internal/transaction/saleinvoice/repositories"
	"smlcloudplatform/pkg/microservice"
	msModels "smlcloudplatform/pkg/microservice/models"
	"time"
)

type ISaleInvoiceTransactionAdminService interface {
	ReSyncSaleInvoiceDoc(shopID string) error
	ReSyncSaleInvoiceDeleteDoc(shopID string) error
}

type SaleInvoiceTransactionAdminService struct {
	mongoRepo       ISaleInvoiceTransactionAdminRepository
	kafkaRepo       saleInvoiceRepository.ISaleInvoiceMessageQueueRepository
	timeoutDuration time.Duration
}

func NewSaleInvoiceTransactionAdminService(pst microservice.IPersisterMongo, kfProducer microservice.IProducer) ISaleInvoiceTransactionAdminService {

	mongoRepo := NewSaleInvoiceTransactionAdminRepository(pst)
	kafkaRepo := saleInvoiceRepository.NewSaleInvoiceMessageQueueRepository(kfProducer)

	return &SaleInvoiceTransactionAdminService{
		mongoRepo:       mongoRepo,
		kafkaRepo:       kafkaRepo,
		timeoutDuration: time.Duration(30) * time.Second,
	}
}

func (s *SaleInvoiceTransactionAdminService) ReSyncSaleInvoiceDoc(shopID string) error {

	pageRequest := msModels.Pageable{
		Limit: 20,
		Page:  1,
		Sorts: []msModels.KeyInt{
			{
				Key:   "guidfixed",
				Value: -1,
			},
		},
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), s.timeoutDuration)
		defer cancel()

		docs, pages, err := s.mongoRepo.FindPage(ctx, shopID, nil, pageRequest)
		if err != nil {
			return err
		}

		// barcodes, pages, err := svc.mongoRepo.FindPage(shopID, nil, pageRequest)
		// 	if err != nil {
		// 		return err
		// 	}

		err = s.kafkaRepo.CreateInBatch(docs)
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

func (s *SaleInvoiceTransactionAdminService) ReSyncSaleInvoiceDeleteDoc(shopID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), s.timeoutDuration)
	defer cancel()

	docs, err := s.mongoRepo.FindSaleInvoiceDeleteByShopID(ctx, shopID)

	if err != nil {
		return err
	}

	err = s.kafkaRepo.DeleteInBatch(docs)
	if err != nil {
		return err
	}
	return nil
}
