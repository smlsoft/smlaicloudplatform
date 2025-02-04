package payment

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/transaction/payment/models"
	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg config.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.TransactionPayment{},
	)
	return nil
}
