package api

import (
	"smlaicloudplatform/internal/logger"
	"smlaicloudplatform/internal/shop"
	shopModel "smlaicloudplatform/internal/shop/models"
	chartOfAccountModels "smlaicloudplatform/internal/vfgl/chartofaccount/models"
	journalModels "smlaicloudplatform/internal/vfgl/journal/models"
	"smlaicloudplatform/pkg/microservice"
)

type IMigrationService interface {
	ImportShop(shops []shopModel.ShopDoc) error
	ImportJournal(journals []journalModels.JournalDoc) error
	ImportChartOfAccount(chartOfAccounts []chartOfAccountModels.ChartOfAccountDoc) error
}

type MigrationService struct {
	logger         logger.ILogger
	mongoPersister microservice.IPersisterMongo
	mqPersister    microservice.IProducer
	userShopRepo   shop.IShopUserRepository
}

func NewMigrationService(logger logger.ILogger, mongoPersister microservice.IPersisterMongo, mqPersister microservice.IProducer) *MigrationService {

	userShopRepo := shop.NewShopUserRepository(mongoPersister)

	return &MigrationService{
		logger:         logger,
		mongoPersister: mongoPersister,
		mqPersister:    mqPersister,
		userShopRepo:   userShopRepo,
	}
}
