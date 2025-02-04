package datamigration

import (
	"smlaicloudplatform/internal/shop"
	shopModel "smlaicloudplatform/internal/shop/models"
	accountGroupRepositories "smlaicloudplatform/internal/vfgl/accountgroup/repositories"
	accountModel "smlaicloudplatform/internal/vfgl/chartofaccount/models"
	chartofaccountrepositories "smlaicloudplatform/internal/vfgl/chartofaccount/repositories"
	journalBookRepositories "smlaicloudplatform/internal/vfgl/journalbook/repositories"

	// chartOfAccountServices "smlaicloudplatform/internal/vfgl/chartofaccount/services"

	"smlaicloudplatform/internal/logger"
	journalModels "smlaicloudplatform/internal/vfgl/journal/models"
	"smlaicloudplatform/pkg/microservice"
)

type IMigrationService interface {
	ImportShop(shops []shopModel.ShopDoc) error
	ImportJournal(journals []journalModels.JournalDoc) error
	ImportChartOfAccount(chartOfAccounts []accountModel.ChartOfAccountDoc) error
	ResyncChartOfAccount(charts []accountModel.ChartOfAccountDoc) error
	InitCenterAccountGroup() error
	InitialChartOfAccountCenter() error
	InitJournalBookCenter() error
}

type MigrationService struct {
	logger          logger.ILogger
	mongoPersister  microservice.IPersisterMongo
	mqPersister     microservice.IProducer
	userShopRepo    shop.IShopUserRepository
	chartRepo       chartofaccountrepositories.IChartOfAccountRepository
	chartMQRepo     chartofaccountrepositories.IChartOfAccountMQRepository
	journalBookRepo journalBookRepositories.IJournalBookMongoRepository

	accountGroupRepo accountGroupRepositories.IAccountGroupMongoRepository
	// chartService   chartOfAccountServices.ChartOfAccountHttpService
}

func NewMigrationService(logger logger.ILogger, mongoPersister microservice.IPersisterMongo, mqPersister microservice.IProducer) *MigrationService {

	chartRepo := chartofaccountrepositories.NewChartOfAccountRepository(mongoPersister)
	chartMQRepo := chartofaccountrepositories.NewChartOfAccountMQRepository(mqPersister)

	userShopRepo := shop.NewShopUserRepository(mongoPersister)
	accountGroupRepo := accountGroupRepositories.NewAccountGroupMongoRepository(mongoPersister)
	journalBookRepo := journalBookRepositories.NewJournalBookMongoRepository(mongoPersister)
	// chartService := chartOfAccountServices.NewChartOfAccountHttpService(chartRepo, nil, chartMQRepo)
	return &MigrationService{
		logger:           logger,
		mongoPersister:   mongoPersister,
		mqPersister:      mqPersister,
		userShopRepo:     userShopRepo,
		chartRepo:        chartRepo,
		chartMQRepo:      chartMQRepo,
		accountGroupRepo: accountGroupRepo,
		journalBookRepo:  journalBookRepo,
	}
}
