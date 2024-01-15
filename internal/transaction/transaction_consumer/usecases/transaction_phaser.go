package usecases

import (
	"smlcloudplatform/internal/transaction/models"
	purchaseModels "smlcloudplatform/internal/transaction/purchase/models"
)

type ITransactionPhaser interface {
	PhasePurchaseDoc(doc *purchaseModels.PurchaseDoc) (*models.StockTransaction, error)
}

func NewTransactionPhaser() ITransactionPhaser {
	return &TransactionPhaser{}
}

type TransactionPhaser struct{}

func (p *TransactionPhaser) PhasePurchaseDoc(doc *purchaseModels.PurchaseDoc) (*models.StockTransaction, error) {
	return nil, nil
}
