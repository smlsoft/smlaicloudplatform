package creditor

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/debtaccount/creditor/models"
	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg config.IConfig) {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.CreditorPG{},
	)
}
