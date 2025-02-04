package transactionconsumer

import (
	pkgConfig "smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/transaction/models"
	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg pkgConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.StockTransaction{},
		models.StockTransactionDetail{},
		models.CreditorTransactionPG{},
		models.DebtorTransactionPG{},

		// models.PurchaseTransactionPG{},
		// models.PurchaseTransactionDetailPG{},
		// models.PurchaseTransactionPG{},
		// models.PurchaseTransactionDetailPG{},

		// models.SaleInvoiceTransactionPG{},
		// models.SaleInvoiceTransactionDetailPG{},
	)
	return nil
}
