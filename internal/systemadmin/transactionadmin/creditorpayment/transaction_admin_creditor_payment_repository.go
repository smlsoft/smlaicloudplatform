package creditorpayment

import (
	"context"
	creditorPaymentModels "smlaicloudplatform/internal/transaction/pay/models"
	"smlaicloudplatform/pkg/microservice"
)

type ICreditorPaymentTransactionAdminRepository interface {
	FindCreditorPaymentDocByShopID(ctx context.Context, shopID string) ([]creditorPaymentModels.PayDoc, error)
}

type CreditorPaymentTransactionAdminRepository struct {
	pst microservice.IPersisterMongo
}

func NewCreditorPaymentTransactionAdminRepository(pst microservice.IPersisterMongo) ICreditorPaymentTransactionAdminRepository {
	return &CreditorPaymentTransactionAdminRepository{
		pst: pst,
	}
}

func (r CreditorPaymentTransactionAdminRepository) FindCreditorPaymentDocByShopID(ctx context.Context, shopID string) ([]creditorPaymentModels.PayDoc, error) {

	docs := []creditorPaymentModels.PayDoc{}

	err := r.pst.Find(ctx, &creditorPaymentModels.PayDoc{}, map[string]interface{}{"shopid": shopID}, &docs)
	if err != nil {
		return nil, err
	}

	return docs, nil
}
