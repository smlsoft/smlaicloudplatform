package shift

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/pos/shift/models"
	"smlaicloudplatform/pkg/microservice"
)

func MigrationDatabase(ms *microservice.Microservice, cfg config.IConfig) {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.ShiftPG{},
	)
}
