package warehouse

import (
	pkgConfig "smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/warehouse/models"
	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg pkgConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.WarehousePG{},
	)
	return nil
}
