package usecases

import (
	"smlcloudplatform/pkg/transaction/models"
	purchaseModels "smlcloudplatform/pkg/transaction/purchase/models"
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
