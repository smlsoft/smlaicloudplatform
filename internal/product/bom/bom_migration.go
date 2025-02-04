package bom

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/product/bom/models"
	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg config.IConfig) {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.ProductBarcodeBOMViewPG{},
	)
}
