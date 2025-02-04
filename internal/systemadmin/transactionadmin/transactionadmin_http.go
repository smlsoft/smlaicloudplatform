package transactionadmin

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/creditorpayment"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/debtorpayment"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/purchase"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/purchasereturn"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/saleinvoice"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/saleinvoicereturn"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/stockadjustment"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/stockbalance"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/stockpickupproduct"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/stockreceiveproduct"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/stockreturnproduct"
	"smlaicloudplatform/internal/systemadmin/transactionadmin/stocktransfer"
	"smlaicloudplatform/pkg/microservice"
)

type ITransactionAdminHttp interface {
	RegisterHttp(ms *microservice.Microservice, prefix string)
}

type TransactionAdminHttp struct {
	purchaseAdminHttp          purchase.IPurchaseTransactionAdminHttp
	purchaseReturnAdminHttp    purchasereturn.IPurchaseReturnTransactionAdminHttp
	saleInvoiceAdminHttp       saleinvoice.ISaleInvoiceTransactionAdminHttp
	saleInvoiceReturnAdminHttp saleinvoicereturn.ISaleInvoiceReturnTransactionAdminHttp
	stockAdjustAdminHttp       stockadjustment.IStockAdjustmentTransactionAdminHttp
	stockReceiveProductHttp    stockreceiveproduct.IStockReceiveTransactionAdminHttp
	stockBalanceAdminHttp      stockbalance.IStockBalanceTransactionAdminHttp
	stockPickupProductHttp     stockpickupproduct.IStockPickupTransactionAdminHttp
	stockReturnProductHttp     stockreturnproduct.IStockReturnProductTransactionAdminHttp
	stockTransferHttp          stocktransfer.IStockTransferTransactionAdminHttp
	creditorPaymentAdminHttp   creditorpayment.ICreditorPaymentTransactionAdminHttp
	debtorPaymentAdminHttp     debtorpayment.IDebtorPaymentTransactionAdminHttp
}

func NewTransactionAdminHttp(ms *microservice.Microservice, cfg config.IConfig) ITransactionAdminHttp {

	purchaseAdminHttp := purchase.NewPurchaseTransactionAdminHttp(ms, cfg)
	purchaseReturnAdminHttp := purchasereturn.NewPurchaseReturnTransactionAdminHttp(ms, cfg)
	saleInvoiceAdminHttp := saleinvoice.NewSaleInvoiceTransactionAdminHttp(ms, cfg)
	saleInvoiceReturnAdminHttp := saleinvoicereturn.NewSaleInvoiceReturnTransactionAdminHttp(ms, cfg)
	stockAdjustAdminHttp := stockadjustment.NewStockAdjustmentTransactionAdminHttp(ms, cfg)
	stockReceiveProductHttp := stockreceiveproduct.NewStockReceiveTransactionAdminHttp(ms, cfg)
	stockPickupProductHttp := stockpickupproduct.NewStockPickupTransactionAdminHttp(ms, cfg)
	stockReturnProductHttp := stockreturnproduct.NewStockReturnProductTransactionAdminHttp(ms, cfg)
	stockTransferHttp := stocktransfer.NewStockTransferTransactionAdminHttp(ms, cfg)
	creditorPaymentAdminHttp := creditorpayment.NewCreditorPaymentTransactionAdminHttp(ms, cfg)
	debtorPaymentAdminHttp := debtorpayment.NewDebtorPaymentTransactionAdminHttp(ms, cfg)
	stockBalanceAdminHttp := stockbalance.NewStockBalanceTransactionAdminHttp(ms, cfg)

	return &TransactionAdminHttp{
		purchaseAdminHttp:          purchaseAdminHttp,
		purchaseReturnAdminHttp:    purchaseReturnAdminHttp,
		saleInvoiceAdminHttp:       saleInvoiceAdminHttp,
		saleInvoiceReturnAdminHttp: saleInvoiceReturnAdminHttp,
		stockAdjustAdminHttp:       stockAdjustAdminHttp,
		stockReceiveProductHttp:    stockReceiveProductHttp,
		stockPickupProductHttp:     stockPickupProductHttp,
		stockReturnProductHttp:     stockReturnProductHttp,
		stockTransferHttp:          stockTransferHttp,
		creditorPaymentAdminHttp:   creditorPaymentAdminHttp,
		debtorPaymentAdminHttp:     debtorPaymentAdminHttp,
		stockBalanceAdminHttp:      stockBalanceAdminHttp,
	}
}

func (s *TransactionAdminHttp) RegisterHttp(ms *microservice.Microservice, prefix string) {
	s.purchaseAdminHttp.RegisterHttp(ms, prefix)
	s.purchaseReturnAdminHttp.RegisterHttp(ms, prefix)
	s.saleInvoiceAdminHttp.RegisterHttp(ms, prefix)
	s.saleInvoiceReturnAdminHttp.RegisterHttp(ms, prefix)
	s.stockAdjustAdminHttp.RegisterHttp(ms, prefix)
	s.stockReceiveProductHttp.RegisterHttp(ms, prefix)
	s.stockPickupProductHttp.RegisterHttp(ms, prefix)
	s.stockReturnProductHttp.RegisterHttp(ms, prefix)
	s.stockTransferHttp.RegisterHttp(ms, prefix)
	s.creditorPaymentAdminHttp.RegisterHttp(ms, prefix)
	s.debtorPaymentAdminHttp.RegisterHttp(ms, prefix)
	s.stockBalanceAdminHttp.RegisterHttp(ms, prefix)
}
