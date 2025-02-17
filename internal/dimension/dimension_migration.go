package dimension

import (
	msConfig "smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/dimension/models"

	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg msConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.DimensionPg{},
		models.DimensionItemPg{},
	)
	return nil
}
