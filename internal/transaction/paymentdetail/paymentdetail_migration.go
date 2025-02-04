package paymentdetail

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/transaction/paymentdetail/models"
	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg config.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.TransactionPaymentDetail{},
	)
	return nil
}
