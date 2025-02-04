package systemadmin

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/systemadmin/accountadmin"
	"smlaicloudplatform/internal/systemadmin/chartofaccountadmin"
	"smlaicloudplatform/internal/systemadmin/creditoradmin"
	"smlaicloudplatform/internal/systemadmin/datamigration"
	"smlaicloudplatform/internal/systemadmin/debtoradmin"
	journal "smlaicloudplatform/internal/systemadmin/journaladmin"
	"smlaicloudplatform/internal/systemadmin/productadmin"
	"smlaicloudplatform/internal/systemadmin/servicetools"
	"smlaicloudplatform/internal/systemadmin/shopadmin"
	"smlaicloudplatform/internal/systemadmin/transactionadmin"
	"smlaicloudplatform/pkg/microservice"
)

const SYSTEM_ADMIN_ROUTE_PREFIX = "/systemadm"

type ISystemAdmin interface {
	RegisterHttp()
}

type SystemAdmin struct {
	ms                      *microservice.Microservice
	migrationHttp           datamigration.IMigrationAPI
	serviceTools            servicetools.IServiceTools
	shopAdminHttp           shopadmin.IShopAdminHttp
	accountAdminHttp        accountadmin.IAccountAdminHttp
	chartOfAccountAdminHttp chartofaccountadmin.IChartOfAccountAdminHttp
	productAdminHttp        productadmin.IProductAdminHttp
	creditorAdminHttp       creditoradmin.ICreditorAdminHttp
	debtorAdminHttp         debtoradmin.IDebtorAdminHttp
	transactionAdminHttp    transactionadmin.ITransactionAdminHttp
	journalAdminHttp        journal.IJournalTransactionAdminHttp
}

func NewSystemAdmin(ms *microservice.Microservice, cfg config.IConfig) ISystemAdmin {

	migrationHttp := datamigration.NewMigrationAPI(ms, cfg)
	servicetools := servicetools.NewServiceTools(ms.Logger, cfg, ms.MongoPersister(cfg.MongoPersisterConfig()))
	shopAdminHttp := shopadmin.NewShopAdminHttp(ms, cfg)
	accountAdminHttp := accountadmin.NewAccountAdminHttp(ms, cfg)
	productAdminHttp := productadmin.NewProductAdminHttp(ms, cfg)
	creditorAdminHttp := creditoradmin.NewCreditorAdminHttp(ms, cfg)
	debtorAdminHttp := debtoradmin.NewDebtorAdminHttp(ms, cfg)
	transactionAdminHttp := transactionadmin.NewTransactionAdminHttp(ms, cfg)
	chartOfAccountAdminHttp := chartofaccountadmin.NewChartOfAccountAdminHttp(ms, cfg)
	journalAdminHttp := journal.NewJournalTransactionAdminHttp(ms, cfg)

	return &SystemAdmin{
		ms:                      ms,
		migrationHttp:           migrationHttp,
		serviceTools:            servicetools,
		shopAdminHttp:           shopAdminHttp,
		accountAdminHttp:        accountAdminHttp,
		productAdminHttp:        productAdminHttp,
		creditorAdminHttp:       creditorAdminHttp,
		debtorAdminHttp:         debtorAdminHttp,
		transactionAdminHttp:    transactionAdminHttp,
		chartOfAccountAdminHttp: chartOfAccountAdminHttp,
		journalAdminHttp:        journalAdminHttp,
	}
}

func (s *SystemAdmin) RegisterHttp() {
	s.migrationHttp.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
	s.serviceTools.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
	s.shopAdminHttp.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
	s.accountAdminHttp.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
	s.productAdminHttp.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
	s.transactionAdminHttp.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
	s.chartOfAccountAdminHttp.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
	s.journalAdminHttp.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
	s.creditorAdminHttp.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
	s.debtorAdminHttp.RegisterHttp(s.ms, SYSTEM_ADMIN_ROUTE_PREFIX)
}
