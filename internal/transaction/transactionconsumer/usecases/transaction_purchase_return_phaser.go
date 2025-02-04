package usecases

import (
	"smlaicloudplatform/internal/transaction/models"
	purchaseReturnModels "smlaicloudplatform/internal/transaction/purchasereturn/models"
)

type PurchaseReturnTransactionPhaser struct{}

func (p PurchaseReturnTransactionPhaser) PhaseSingleDoc(msg string) (*models.PurchaseReturnTransactionPG, error) {
	return nil, nil
}

func (p PurchaseReturnTransactionPhaser) PhaseMultipleDoc(input string) (*[]models.PurchaseReturnTransactionPG, error) {
	return nil, nil
}

func (p PurchaseReturnTransactionPhaser) PhaseStockTransactionPurchaseReturnDoc(purchaseDoc purchaseReturnModels.PurchaseReturnDoc) (*models.PurchaseReturnTransactionPG, error) {
	return nil, nil
}
