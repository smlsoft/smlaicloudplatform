package services

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/transaction/models"
	purchaseModels "smlcloudplatform/pkg/transaction/purchase/models"
	"smlcloudplatform/pkg/transaction/transaction_consumer/repositories"
	"smlcloudplatform/pkg/transaction/transaction_consumer/usecases"
)

type ITransactionConsumerService interface {
	UpSert(shopID string, docNo string, doc models.StockTransaction) error
	UpSertPurchaseTrx(shopID string, docNo string, doc purchaseModels.PurchaseDoc) error
}

func NewPurchaseConsumerService(pst microservice.IPersister) ITransactionConsumerService {

	pgRepo := repositories.NewTransactionPGRepository(pst)
	return &TransactionConsumerService{
		transactionPGRepo: pgRepo,
	}
}

type TransactionConsumerService struct {
	phaser            usecases.ITransactionPhaser
	transactionPGRepo repositories.ITransactionPGRepository
}

func (svc *TransactionConsumerService) UpSert(shopID string, docNo string, doc models.StockTransaction) error {

	findTrx, err := svc.transactionPGRepo.Get(shopID, docNo)
	if err != nil {
		return err
	}

	if findTrx == nil {
		err = svc.transactionPGRepo.Create(doc)
		if err != nil {
			return err
		}
	} else {
		err = svc.transactionPGRepo.Update(shopID, docNo, doc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *TransactionConsumerService) UpSertPurchaseTrx(shopID string, docNo string, doc purchaseModels.PurchaseDoc) error {
	trx, err := svc.phaser.PhasePurchaseDoc(&doc)
	if err != nil {
		return err
	}

	err = svc.UpSert(trx.ShopID, trx.DocNo, *trx)
	return err
}
