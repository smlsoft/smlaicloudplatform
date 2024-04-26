package main

import (
	"fmt"
	"os"
	migrationAPI "smlcloudplatform/cmd/migrationapi/api"
	"smlcloudplatform/docs"
	"smlcloudplatform/internal/apikeyservice"
	"smlcloudplatform/internal/authentication"
	"smlcloudplatform/internal/channel/salechannel"
	"smlcloudplatform/internal/channel/transportchannel"
	"smlcloudplatform/internal/config"
	"smlcloudplatform/internal/debtaccount/creditor"
	"smlcloudplatform/internal/debtaccount/creditorgroup"
	"smlcloudplatform/internal/debtaccount/customer"
	"smlcloudplatform/internal/debtaccount/customergroup"
	"smlcloudplatform/internal/debtaccount/debtor"
	"smlcloudplatform/internal/debtaccount/debtorgroup"
	"smlcloudplatform/internal/dimension"
	"smlcloudplatform/internal/documentwarehouse/documentimage"
	"smlcloudplatform/internal/filestatus"
	"smlcloudplatform/internal/images"
	"smlcloudplatform/internal/masterexpense"
	"smlcloudplatform/internal/masterincome"
	"smlcloudplatform/internal/mastersync"
	"smlcloudplatform/internal/member"
	"smlcloudplatform/internal/notify"
	"smlcloudplatform/internal/ocr"
	order_device "smlcloudplatform/internal/order/device"
	order_setting "smlcloudplatform/internal/order/setting"
	"smlcloudplatform/internal/organization/branch"
	"smlcloudplatform/internal/organization/businesstype"
	"smlcloudplatform/internal/organization/department"
	"smlcloudplatform/internal/payment/bankmaster"
	"smlcloudplatform/internal/payment/bookbank"
	"smlcloudplatform/internal/payment/qrpayment"
	"smlcloudplatform/internal/paymentmaster"
	pos_media "smlcloudplatform/internal/pos/media"
	pos_setting "smlcloudplatform/internal/pos/setting"
	"smlcloudplatform/internal/pos/shift"
	"smlcloudplatform/internal/pos/temp"
	"smlcloudplatform/internal/product/color"
	"smlcloudplatform/internal/product/eorder"
	"smlcloudplatform/internal/product/option"
	"smlcloudplatform/internal/product/optionpattern"
	"smlcloudplatform/internal/product/ordertype"
	"smlcloudplatform/internal/product/productbarcode"
	"smlcloudplatform/internal/product/productcategory"
	"smlcloudplatform/internal/product/productgroup"
	"smlcloudplatform/internal/product/producttype"
	"smlcloudplatform/internal/product/promotion"
	"smlcloudplatform/internal/product/unit"
	"smlcloudplatform/internal/productimport"
	"smlcloudplatform/internal/productsection/sectionbranch"
	"smlcloudplatform/internal/productsection/sectionbusinesstype"
	"smlcloudplatform/internal/productsection/sectiondepartment"
	"smlcloudplatform/internal/restaurant/device"
	"smlcloudplatform/internal/restaurant/kitchen"
	"smlcloudplatform/internal/restaurant/printer"
	"smlcloudplatform/internal/restaurant/settings"
	"smlcloudplatform/internal/restaurant/staff"
	"smlcloudplatform/internal/restaurant/table"
	"smlcloudplatform/internal/restaurant/zone"
	"smlcloudplatform/internal/shop"
	"smlcloudplatform/internal/shop/employee"
	"smlcloudplatform/internal/shopdesign/zonedesign"
	"smlcloudplatform/internal/slipimage"
	"smlcloudplatform/internal/smsreceive/smstransaction"
	"smlcloudplatform/internal/stockbalanceimport"
	"smlcloudplatform/internal/stockprocess"
	"smlcloudplatform/internal/systemadmin"
	"smlcloudplatform/internal/task"
	"smlcloudplatform/internal/transaction/documentformate"
	"smlcloudplatform/internal/transaction/paid"
	"smlcloudplatform/internal/transaction/pay"
	"smlcloudplatform/internal/transaction/payment"
	"smlcloudplatform/internal/transaction/paymentdetail"
	"smlcloudplatform/internal/transaction/purchase"
	"smlcloudplatform/internal/transaction/purchaseorder"
	"smlcloudplatform/internal/transaction/purchasereturn"
	"smlcloudplatform/internal/transaction/saleinvoice"
	"smlcloudplatform/internal/transaction/saleinvoicebomprice"
	"smlcloudplatform/internal/transaction/saleinvoicereturn"
	"smlcloudplatform/internal/transaction/smltransaction"
	"smlcloudplatform/internal/transaction/stockadjustment"
	"smlcloudplatform/internal/transaction/stockbalance"
	"smlcloudplatform/internal/transaction/stockbalancedetail"
	"smlcloudplatform/internal/transaction/stockpickupproduct"
	"smlcloudplatform/internal/transaction/stockreceiveproduct"
	"smlcloudplatform/internal/transaction/stockreturnproduct"
	"smlcloudplatform/internal/transaction/stocktransfer"
	"smlcloudplatform/internal/vfgl/accountgroup"
	"smlcloudplatform/internal/vfgl/accountperiodmaster"
	"smlcloudplatform/internal/vfgl/chartofaccount"
	"smlcloudplatform/internal/vfgl/journal"
	"smlcloudplatform/internal/vfgl/journalbook"
	"smlcloudplatform/internal/vfgl/journalreport"
	"smlcloudplatform/internal/warehouse"
	"smlcloudplatform/pkg/microservice"
	"time"

	"smlcloudplatform/internal/transaction/transactionconsumer"

	paid_consumer "smlcloudplatform/internal/transaction/transactionconsumer/paid"
	pay_consumer "smlcloudplatform/internal/transaction/transactionconsumer/pay"

	purchase_consumer "smlcloudplatform/internal/transaction/transactionconsumer/purchase"
	purchasereturn_consumer "smlcloudplatform/internal/transaction/transactionconsumer/purchasereturn"

	saleinvoice_consumer "smlcloudplatform/internal/transaction/transactionconsumer/saleinvoice"
	saleinvoicereutrn_consumer "smlcloudplatform/internal/transaction/transactionconsumer/saleinvoicereturn"

	stockadjustment_consumer "smlcloudplatform/internal/transaction/transactionconsumer/stockadjustment"
	stockpickupproduct_consumer "smlcloudplatform/internal/transaction/transactionconsumer/stockpickupproduct"
	stockreceiveproduct_consumer "smlcloudplatform/internal/transaction/transactionconsumer/stockreceiveproduct"
	stockreturnproduct_consumer "smlcloudplatform/internal/transaction/transactionconsumer/stockreturnproduct"

	stocktranferproduct_consumer "smlcloudplatform/internal/transaction/transactionconsumer/stocktransfer"

	creditorpayment_consumer "smlcloudplatform/internal/transaction/transactionconsumer/creditorpayment"
	debtorpayment_consumer "smlcloudplatform/internal/transaction/transactionconsumer/debtorpayment"

	"github.com/joho/godotenv"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	env := os.Getenv("MODE")
	if env == "" {
		os.Setenv("MODE", "development")
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() //
}

// @title           SML Cloud Platform API
// @version         1.0
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @securityDefinitions.apikey  AccessToken
// @in                          header
// @name                        Authorization

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes http https
func main() {

	devApiMode := os.Getenv("DEV_API_MODE")
	host := os.Getenv("HOST_API")
	if host != "" {
		fmt.Printf("Host: %v\n", host)
		docs.SwaggerInfo.Host = host
	}

	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	ms.HttpUsePrometheus()

	if devApiMode == "" || devApiMode == "2" {

		ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

		cacher := ms.Cacher(cfg.CacherConfig())
		authService := microservice.NewAuthService(cacher, 24*3*time.Hour, 24*30*time.Hour)
		publicPath := []string{
			"/migrationtools/",
			"/swagger/*",

			"/tokenlogin",

			"/login",
			"/login/phone-number",
			"/register",
			"/refresh",
			"/register-phonenumber",
			"/register/exists-username",
			"/register/exists-phonenumber",
			"/send-phonenumber-otp",

			"/employee/login",

			"/images*",
			"/productimage/*",

			"/healthz",
			"/ws",
			"/metrics",
			"/e-order/product",
			"/e-order/category",
			"/e-order/product-barcode",
			"/e-order/shop-info",
			"/e-order/shop-info/v1.1",
			"/e-order/restaurant/zone",
			"/e-order/restaurant/kitchen",
			"/e-order/restaurant/table",
			"/e-order/sale-invoice/last-pos-docno",
			"/e-order/notify",
			"/line-notify",
		}

		exceptShopPath := []string{
			"/shop",
			"/profile",
			"/list-shop",
			"/select-shop",
			"/create-shop",
			"/favorite-shop",
		}
		ms.HttpMiddleware(authService.MWFuncWithRedisMixShop(cacher, exceptShopPath, publicPath...))
		ms.RegisterLivenessProbeEndpoint("/healthz")
		ms.HttpUseCors()
		ms.HttpPreRemoveTrailingSlash()
		// ms.Echo().GET("/healthz", func(c echo.Context) error {
		// 	return c.String(http.StatusOK, "ok")
		// })
		azureFileBlob := microservice.NewPersisterAzureBlob()
		imagePersister := microservice.NewPersisterImage(azureFileBlob)

		httpServices := []HttpRegister{

			apikeyservice.NewApiKeyServiceHttp(ms, cfg),
			authentication.NewAuthenticationHttp(ms, cfg),
			apikeyservice.NewApiKeyServiceHttp(ms, cfg),
			shop.NewShopHttp(ms, cfg),

			shop.NewShopMemberHttp(ms, cfg),
			employee.NewEmployeeHttp(ms, cfg), member.NewMemberHttp(ms, cfg),

			option.NewOptionHttp(ms, cfg),
			unit.NewUnitHttp(ms, cfg),
			optionpattern.NewOptionPatternHttp(ms, cfg),
			color.NewColorHttp(ms, cfg),

			//product
			productcategory.NewProductCategoryHttp(ms, cfg),
			productbarcode.NewProductBarcodeHttp(ms, cfg),
			// product.NewProductHttp(ms, cfg),
			productgroup.NewProductGroupHttp(ms, cfg),
			producttype.NewProductTypeHttp(ms, cfg),

			images.NewImagesHttp(ms, cfg, imagePersister),

			// restaurant
			zone.NewZoneHttp(ms, cfg),
			table.NewTableHttp(ms, cfg),
			printer.NewPrinterHttp(ms, cfg),
			kitchen.NewKitchenHttp(ms, cfg),
			zonedesign.NewZoneDesignHttp(ms, cfg),
			settings.NewRestaurantSettingsHttp(ms, cfg),
			device.NewDeviceHttp(ms, cfg),
			staff.NewStaffHttp(ms, cfg),

			chartofaccount.NewChartOfAccountHttp(ms, cfg),
			journal.NewJournalHttp(ms, cfg),
			journal.NewJournalWs(ms, cfg),
			journalreport.NewJournalReportHttp(ms, cfg),
			accountgroup.NewAccountGroupHttp(ms, cfg),
			journalbook.NewJournalBookHttp(ms, cfg),

			documentimage.NewDocumentImageHttp(ms, cfg),
			mastersync.NewMasterSyncHttp(ms, cfg),
			smstransaction.NewSmsTransactionHttp(ms, cfg),
			paymentmaster.NewPaymentMasterHttp(ms, cfg),
			warehouse.NewWarehouseHttp(ms, cfg),

			accountperiodmaster.NewAccountPeriodMasterHttp(ms, cfg),

			bankmaster.NewBankMasterHttp(ms, cfg),
			bookbank.NewBookBankHttp(ms, cfg),
			qrpayment.NewQrPaymentHttp(ms, cfg),

			task.NewTaskHttp(ms, cfg),
			smltransaction.NewSMLTransactionHttp(ms, cfg),

			// debt account
			creditor.NewCreditorHttp(ms, cfg),
			creditorgroup.NewCreditorGroupHttp(ms, cfg),
			debtor.NewDebtorHttp(ms, cfg),
			debtorgroup.NewDebtorGroupHttp(ms, cfg),

			customer.NewCustomerHttp(ms, cfg),
			customergroup.NewCustomerGroupHttp(ms, cfg),

			branch.NewBranchHttp(ms, cfg),
			department.NewDepartmentHttp(ms, cfg),
			businesstype.NewBusinessTypeHttp(ms, cfg),

			//transaction
			purchase.NewPurchaseHttp(ms, cfg),
			purchasereturn.NewPurchaseReturnHttp(ms, cfg),
			saleinvoice.NewSaleInvoiceHttp(ms, cfg),
			saleinvoicereturn.NewSaleInvoiceReturnHttp(ms, cfg),
			stocktransfer.NewStockTransferHttp(ms, cfg),
			stockreceiveproduct.NewStockReceiveProductHttp(ms, cfg),
			stockreturnproduct.NewStockReturnProductHttp(ms, cfg),
			stockpickupproduct.NewStockPickupProductHttp(ms, cfg),
			stockadjustment.NewStockAdjustmentHttp(ms, cfg),
			paid.NewPaidHttp(ms, cfg),
			pay.NewPayHttp(ms, cfg),
			stockbalance.NewStockBalanceHttp(ms, cfg),
			stockbalancedetail.NewStockBalanceDetailHttp(ms, cfg),
			purchaseorder.NewPurchaseOrderHttp(ms, cfg),

			//product section
			sectionbranch.NewSectionBranchHttp(ms, cfg),
			sectiondepartment.NewSectionDepartmentHttp(ms, cfg),
			sectionbusinesstype.NewSectionBusinessTypeHttp(ms, cfg),

			//channel
			salechannel.NewSaleChannelHttp(ms, cfg),
			transportchannel.NewTransportChannelHttp(ms, cfg),

			// e-order
			eorder.NewEOrderHttp(ms, cfg),

			// promiotions
			promotion.NewPromotionHttp(ms, cfg),

			ordertype.NewOrderTypeHttp(ms, cfg),

			pos_setting.NewSettingHttp(ms, cfg),
			pos_media.NewMediaHttp(ms, cfg),
			shift.NewShiftHttp(ms, cfg),

			order_setting.NewSettingHttp(ms, cfg),
			order_device.NewDeviceHttp(ms, cfg),

			documentformate.NewDocumentFormateHttp(ms, cfg),
			ocr.NewOcrHttp(ms, cfg),

			notify.NewNotifyHttp(ms, cfg),
			slipimage.NewSlipImageHttp(ms, cfg),

			// import

			stockbalanceimport.NewStockBalanceImportHttp(ms, cfg),
			productimport.NewProductImportHttp(ms, cfg),

			dimension.NewDimensionHttp(ms, cfg),

			// master
			masterexpense.NewMasterExpenseHttp(ms, cfg),
			masterincome.NewMasterIncomeHttp(ms, cfg),

			// system admin
			systemadmin.NewSystemAdmin(ms, cfg),

			temp.NewPOSTempHttp(ms, cfg),
			filestatus.NewFileStatusHttp(ms, cfg),

			// member
			member.NewMemberHttp(ms, cfg),

			saleinvoicebomprice.NewSaleInvoicePriceHttp(ms, cfg),
		}

		serviceStartHttp(ms, httpServices...)

	}

	// Migration
	if devApiMode == "3" {
		// migration db only
		journal.MigrationJournalTable(ms, cfg)
		chartofaccount.MigrationChartOfAccountTable(ms, cfg)
		productbarcode.MigrationDatabase(ms, cfg)
		// transactionconsumer.MigrationDatabase(ms, cfg)
		// payment migration
		payment.MigrationDatabase(ms, cfg)
		paymentdetail.MigrationDatabase(ms, cfg)
		pay_consumer.MigrationDatabase(ms, cfg)
		paid_consumer.MigrationDatabase(ms, cfg)
		transactionconsumer.MigrationDatabase(ms, cfg)
		purchase_consumer.MigrationDatabase(ms, cfg)
		purchasereturn_consumer.MigrationDatabase(ms, cfg)
		saleinvoice_consumer.MigrationDatabase(ms, cfg)
		saleinvoicereutrn_consumer.MigrationDatabase(ms, cfg)
		stockreceiveproduct_consumer.MigrationDatabase(ms, cfg)
		stockpickupproduct_consumer.MigrationDatabase(ms, cfg)
		stockreturnproduct_consumer.MigrationDatabase(ms, cfg)
		stockadjustment_consumer.MigrationDatabase(ms, cfg)
		stocktranferproduct_consumer.MigrationDatabase(ms, cfg)
		creditorpayment_consumer.MigrationDatabase(ms, cfg)
		debtorpayment_consumer.MigrationDatabase(ms, cfg)
		warehouse.MigrationDatabase(ms, cfg)

		// debt account
		creditor.MigrationDatabase(ms, cfg)
		debtor.MigrationDatabase(ms, cfg)

		return
	}

	if devApiMode == "1" || devApiMode == "2" {

		ms.RegisterLivenessProbeEndpoint("/healthz")

		consumerGroupName := os.Getenv("CONSUMER_GROUP_NAME")
		if consumerGroupName == "" {
			consumerGroupName = "03"
		}

		// journal.StartJournalComsumeCreated(ms, cfg, consumerGroupName)
		// journal.StartJournalComsumeUpdated(ms, cfg, consumerGroupName)
		// journal.StartJournalComsumeDeleted(ms, cfg, consumerGroupName)
		// journal.StartJournalComsumeBlukCreated(ms, cfg, consumerGroupName)
		ms.RegisterConsumer(journal.InitJournalTransactionConsumer(ms, cfg))

		chartofaccount.StartChartOfAccountConsumerCreated(ms, cfg, consumerGroupName)
		chartofaccount.StartChartOfAccountConsumerUpdated(ms, cfg, consumerGroupName)
		chartofaccount.StartChartOfAccountConsumerDeleted(ms, cfg, consumerGroupName)
		chartofaccount.StartChartOfAccountConsumerBlukCreated(ms, cfg, consumerGroupName)

		// Transaction
		ms.RegisterConsumer(pay_consumer.InitPayTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(paid_consumer.InitPaidTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(stockprocess.NewStockProcessConsumer(ms, cfg))
		ms.RegisterConsumer(purchase_consumer.InitPurchaseTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(purchasereturn_consumer.InitPurchaseReturnTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(saleinvoice_consumer.InitSaleInvoiceTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(saleinvoicereutrn_consumer.InitSaleInvoiceReturnTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(stockreceiveproduct_consumer.InitStockReceiveTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(stockpickupproduct_consumer.InitStockReceiveTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(stockreturnproduct_consumer.InitStockReturnTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(stockadjustment_consumer.InitStockAdjustmentTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(stocktranferproduct_consumer.InitStockTransferTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(creditorpayment_consumer.InitCreditorPaymentTransactionConsumer(ms, cfg))
		ms.RegisterConsumer(debtorpayment_consumer.InitDebtorPaymentTransactionConsumer(ms, cfg))

		// Warehouse
		ms.RegisterConsumer(warehouse.InitWarehouseConsumer(ms, cfg))

		// Debt Account
		ms.RegisterConsumer(creditor.InitCreditorConsumer(ms, cfg))
		ms.RegisterConsumer(debtor.InitDebtorConsumer(ms, cfg))

		consumerServices := []ConsumerRegister{
			task.NewTaskConsumer(ms, cfg),
			productbarcode.NewProductBarcodeConsumer(ms, cfg),
		}

		serviceStartConsumer(ms, consumerServices...)
	}

	ms.RegisterHttp(migrationAPI.NewMigrationAPI(ms, cfg))
	ms.Start()
}

type HttpRegister interface {
	RegisterHttp()
}

func serviceStartHttp(ms *microservice.Microservice, services ...HttpRegister) {
	for _, service := range services {
		ms.RegisterHttp(service)
	}
}

type ConsumerRegister interface {
	RegisterConsumer()
}

func serviceStartConsumer(ms *microservice.Microservice, services ...ConsumerRegister) {
	for _, service := range services {
		service.RegisterConsumer()
	}
}
