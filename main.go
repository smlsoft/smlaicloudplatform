package main

import (
	"fmt"
	"os"
	migrationAPI "smlaicloudplatform/cmd/migrationapi/api"
	"smlaicloudplatform/docs"
	"smlaicloudplatform/internal/apikeyservice"
	"smlaicloudplatform/internal/authentication"
	"smlaicloudplatform/internal/channel/salechannel"
	"smlaicloudplatform/internal/channel/transportchannel"
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/debtaccount/creditor"
	"smlaicloudplatform/internal/debtaccount/creditorgroup"
	"smlaicloudplatform/internal/debtaccount/customer"
	"smlaicloudplatform/internal/debtaccount/customergroup"
	"smlaicloudplatform/internal/debtaccount/debtor"
	"smlaicloudplatform/internal/debtaccount/debtorgroup"
	"smlaicloudplatform/internal/dimension"
	"smlaicloudplatform/internal/documentwarehouse/documentimage"
	"smlaicloudplatform/internal/filestatus"
	"smlaicloudplatform/internal/images"
	"smlaicloudplatform/internal/masterexpense"
	"smlaicloudplatform/internal/masterincome"
	"smlaicloudplatform/internal/mastersync"
	"smlaicloudplatform/internal/media"
	"smlaicloudplatform/internal/member"
	"smlaicloudplatform/internal/notify"
	"smlaicloudplatform/internal/ocr"
	order_device "smlaicloudplatform/internal/order/device"
	order_setting "smlaicloudplatform/internal/order/setting"
	"smlaicloudplatform/internal/organization/branch"
	"smlaicloudplatform/internal/organization/businesstype"
	"smlaicloudplatform/internal/organization/department"
	"smlaicloudplatform/internal/payment/bankmaster"
	"smlaicloudplatform/internal/payment/bookbank"
	"smlaicloudplatform/internal/payment/qrpayment"
	"smlaicloudplatform/internal/paymentmaster"
	pos_media "smlaicloudplatform/internal/pos/media"
	pos_setting "smlaicloudplatform/internal/pos/setting"
	"smlaicloudplatform/internal/pos/shift"
	"smlaicloudplatform/internal/pos/temp"
	"smlaicloudplatform/internal/product/bom"
	"smlaicloudplatform/internal/product/color"
	"smlaicloudplatform/internal/product/eorder"
	"smlaicloudplatform/internal/product/option"
	"smlaicloudplatform/internal/product/optionpattern"
	"smlaicloudplatform/internal/product/ordertype"
	"smlaicloudplatform/internal/product/productbarcode"
	"smlaicloudplatform/internal/product/productcategory"
	"smlaicloudplatform/internal/product/productgroup"
	"smlaicloudplatform/internal/product/producttype"
	"smlaicloudplatform/internal/product/promotion"
	"smlaicloudplatform/internal/product/unit"
	"smlaicloudplatform/internal/productimport"
	"smlaicloudplatform/internal/productsection/sectionbranch"
	"smlaicloudplatform/internal/productsection/sectionbusinesstype"
	"smlaicloudplatform/internal/productsection/sectiondepartment"
	"smlaicloudplatform/internal/restaurant/device"
	"smlaicloudplatform/internal/restaurant/kitchen"
	"smlaicloudplatform/internal/restaurant/printer"
	"smlaicloudplatform/internal/restaurant/settings"
	"smlaicloudplatform/internal/restaurant/staff"
	"smlaicloudplatform/internal/restaurant/table"
	"smlaicloudplatform/internal/restaurant/zone"
	"smlaicloudplatform/internal/shop"
	"smlaicloudplatform/internal/shop/employee"
	"smlaicloudplatform/internal/shopdesign/zonedesign"
	"smlaicloudplatform/internal/slipimage"
	"smlaicloudplatform/internal/smsreceive/smstransaction"
	"smlaicloudplatform/internal/stockbalanceimport"
	"smlaicloudplatform/internal/stockprocess"
	"smlaicloudplatform/internal/systemadmin"
	"smlaicloudplatform/internal/task"
	"smlaicloudplatform/internal/transaction/documentformate"
	"smlaicloudplatform/internal/transaction/paid"
	"smlaicloudplatform/internal/transaction/pay"
	"smlaicloudplatform/internal/transaction/payment"
	"smlaicloudplatform/internal/transaction/paymentdetail"
	"smlaicloudplatform/internal/transaction/purchase"
	"smlaicloudplatform/internal/transaction/purchaseorder"
	"smlaicloudplatform/internal/transaction/purchasereturn"
	"smlaicloudplatform/internal/transaction/saleinvoice"
	"smlaicloudplatform/internal/transaction/saleinvoicebomprice"
	"smlaicloudplatform/internal/transaction/saleinvoicereturn"
	"smlaicloudplatform/internal/transaction/smltransaction"
	"smlaicloudplatform/internal/transaction/stockadjustment"
	"smlaicloudplatform/internal/transaction/stockbalance"
	"smlaicloudplatform/internal/transaction/stockbalancedetail"
	"smlaicloudplatform/internal/transaction/stockpickupproduct"
	"smlaicloudplatform/internal/transaction/stockreceiveproduct"
	"smlaicloudplatform/internal/transaction/stockreturnproduct"
	"smlaicloudplatform/internal/transaction/stocktransfer"
	"smlaicloudplatform/internal/vfgl/accountgroup"
	"smlaicloudplatform/internal/vfgl/accountperiodmaster"
	"smlaicloudplatform/internal/vfgl/chartofaccount"
	"smlaicloudplatform/internal/vfgl/journal"
	"smlaicloudplatform/internal/vfgl/journalbook"
	"smlaicloudplatform/internal/vfgl/journalreport"
	"smlaicloudplatform/internal/warehouse"
	"smlaicloudplatform/pkg/microservice"
	"time"

	"smlaicloudplatform/internal/transaction/transactionconsumer"

	paid_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/paid"
	pay_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/pay"

	purchase_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/purchase"
	purchasereturn_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/purchasereturn"

	saleinvoice_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/saleinvoice"
	saleinvoicereutrn_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/saleinvoicereturn"

	stockadjustment_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/stockadjustment"
	stockpickupproduct_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/stockpickupproduct"
	stockreceiveproduct_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/stockreceiveproduct"
	stockreturnproduct_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/stockreturnproduct"

	stocktranferproduct_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/stocktransfer"

	creditorpayment_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/creditorpayment"
	debtorpayment_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/debtorpayment"

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
			"/poslogin",
			"/login/email",
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

			// BOM
			bom.NewBOMHttp(ms, cfg),
			saleinvoicebomprice.NewSaleInvoiceBomPriceHttp(ms, cfg),
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
		shift.MigrationDatabase(ms, cfg)

		// BOM
		bom.MigrationDatabase(ms, cfg)
		saleinvoicebomprice.MigrationDatabase(ms, cfg)

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

		// Shift
		ms.RegisterConsumer(shift.InitShiftConsumer(ms, cfg))

		// BOM
		ms.RegisterConsumer(bom.InitBOMConsumer(ms, cfg))
		ms.RegisterConsumer(saleinvoicebomprice.InitSaleInvoiceBomPriceConsumer(ms, cfg))

		consumerServices := []ConsumerRegister{
			task.NewTaskConsumer(ms, cfg),
			productbarcode.NewProductBarcodeConsumer(ms, cfg),
		}

		serviceStartConsumer(ms, consumerServices...)
	}

	ms.RegisterHttp(migrationAPI.NewMigrationAPI(ms, cfg))
	ms.RegisterHttp(media.InitMediaUploadHttp(ms, cfg))
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
