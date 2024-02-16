package saleinvoice

import (
	"context"
	saleInvoiceModels "smlcloudplatform/internal/transaction/saleinvoice/models"
	"smlcloudplatform/pkg/microservice"

	"go.mongodb.org/mongo-driver/bson"
)

type ISaleInvoiceTransactionAdminRepository interface {
	FindSaleInvoiceByShopID(ctx context.Context, shopID string) ([]saleInvoiceModels.SaleInvoiceDoc, error)
	FindSaleInvoiceDeleteByShopID(ctx context.Context, shopID string) ([]saleInvoiceModels.SaleInvoiceDoc, error)
}

type SaleInvoiceTransactionAdminRepository struct {
	pst microservice.IPersisterMongo
}

func NewSaleInvoiceTransactionAdminRepository(pst microservice.IPersisterMongo) ISaleInvoiceTransactionAdminRepository {
	return &SaleInvoiceTransactionAdminRepository{
		pst: pst,
	}
}

func (r SaleInvoiceTransactionAdminRepository) FindSaleInvoiceByShopID(ctx context.Context, shopID string) ([]saleInvoiceModels.SaleInvoiceDoc, error) {
	docList := []saleInvoiceModels.SaleInvoiceDoc{}

	err := r.pst.Find(ctx, &saleInvoiceModels.SaleInvoiceDoc{},
		bson.M{
			"shopid":    shopID,
			"deletedat": bson.M{"$exists": false},
		},
		&docList)
	if err != nil {
		return nil, err
	}

	return docList, nil
}

func (r SaleInvoiceTransactionAdminRepository) FindSaleInvoiceDeleteByShopID(ctx context.Context, shopID string) ([]saleInvoiceModels.SaleInvoiceDoc, error) {
	docList := []saleInvoiceModels.SaleInvoiceDoc{}

	err := r.pst.Find(ctx, &saleInvoiceModels.SaleInvoiceDoc{},
		bson.M{
			"shopid":    shopID,
			"deletedat": bson.M{"$exists": true},
		},
		&docList)
	if err != nil {
		return nil, err
	}

	return docList, nil
}
