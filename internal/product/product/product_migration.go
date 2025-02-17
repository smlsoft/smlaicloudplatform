package products

import (
	msConfig "smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/product/product/models"

	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg msConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.ProductPg{},
		models.ProductDimensionPg{},
	)
	return nil
}
