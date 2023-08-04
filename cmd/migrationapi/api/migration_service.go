package api

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/logger"
	"smlcloudplatform/pkg/shop"
	shopModel "smlcloudplatform/pkg/shop/models"
	chartOfAccountModels "smlcloudplatform/pkg/vfgl/chartofaccount/models"
	journalModels "smlcloudplatform/pkg/vfgl/journal/models"
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
