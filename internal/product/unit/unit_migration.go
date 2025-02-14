package unit

import (
	msConfig "smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/product/unit/models"

	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg msConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.UnitPg{},
	)
	return nil
}
