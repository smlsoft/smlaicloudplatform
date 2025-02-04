package saleinvoicebomprice

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/transaction/saleinvoicebomprice/models"
	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg config.IConfig) {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.SaleInvoiceBomPricePg{},
	)
}
