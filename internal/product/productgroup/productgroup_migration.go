package productgroup

import (
	msConfig "smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/product/productgroup/models"

	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg msConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.ProductGroupPg{},
	)
	return nil
}
