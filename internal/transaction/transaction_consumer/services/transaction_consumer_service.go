package services

import (
	"smlcloudplatform/internal/transaction/models"
	purchaseModels "smlcloudplatform/internal/transaction/purchase/models"
	"smlcloudplatform/internal/transaction/transaction_consumer/repositories"
	"smlcloudplatform/internal/transaction/transaction_consumer/usecases"
	"smlcloudplatform/pkg/microservice"
)

type ITransactionConsumerService interface {
	UpSert(shopID string, docNo string, doc models.StockTransaction) error
	Delete(shopID string, docNo string) error
	// UpSertPurchaseTrx(shopID string, docNo string, doc purchaseModels.PurchaseDoc) error
}

func NewTransactionConsumerService(pst microservice.IPersister) ITransactionConsumerService {

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

func (svc *TransactionConsumerService) Delete(shopID string, docNo string) error {
	err := svc.transactionPGRepo.Delete(shopID, docNo)
	if err != nil {
		return err
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
