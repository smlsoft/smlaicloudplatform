package main

import (
	"fmt"
	"os"
	"smlcloudplatform/docs"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/apikeyservice"
	"smlcloudplatform/pkg/authentication"
	"smlcloudplatform/pkg/channel/salechannel"
	"smlcloudplatform/pkg/channel/transportchannel"
	"smlcloudplatform/pkg/debtaccount/creditor"
	"smlcloudplatform/pkg/debtaccount/creditorgroup"
	"smlcloudplatform/pkg/debtaccount/customer"
	"smlcloudplatform/pkg/debtaccount/customergroup"
	"smlcloudplatform/pkg/debtaccount/debtor"
	"smlcloudplatform/pkg/debtaccount/debtorgroup"
	"smlcloudplatform/pkg/documentwarehouse/documentimage"
	"smlcloudplatform/pkg/images"
	"smlcloudplatform/pkg/mastersync"
	"smlcloudplatform/pkg/member"
	"smlcloudplatform/pkg/organization/branch"
	"smlcloudplatform/pkg/organization/businesstype"
	"smlcloudplatform/pkg/organization/department"
	"smlcloudplatform/pkg/payment/bankmaster"
	"smlcloudplatform/pkg/payment/bookbank"
	"smlcloudplatform/pkg/payment/qrpayment"
	"smlcloudplatform/pkg/paymentmaster"
	"smlcloudplatform/pkg/product/color"
	"smlcloudplatform/pkg/product/eorder"
	"smlcloudplatform/pkg/product/inventory"
	"smlcloudplatform/pkg/product/inventoryimport"
	"smlcloudplatform/pkg/product/inventorysearchconsumer"
	"smlcloudplatform/pkg/product/option"
	"smlcloudplatform/pkg/product/optionpattern"
	"smlcloudplatform/pkg/product/product"
	"smlcloudplatform/pkg/product/productbarcode"
	"smlcloudplatform/pkg/product/productcategory"
	"smlcloudplatform/pkg/product/productgroup"
	"smlcloudplatform/pkg/product/unit"
	"smlcloudplatform/pkg/productsection/sectionbranch"
	"smlcloudplatform/pkg/productsection/sectionbusinesstype"
	"smlcloudplatform/pkg/productsection/sectiondepartment"
	"smlcloudplatform/pkg/restaurant/device"
	"smlcloudplatform/pkg/restaurant/kitchen"
	"smlcloudplatform/pkg/restaurant/printer"
	"smlcloudplatform/pkg/restaurant/restaurantsettings"
	"smlcloudplatform/pkg/restaurant/shoptable"
	"smlcloudplatform/pkg/restaurant/shopzone"
	"smlcloudplatform/pkg/restaurant/staff"
	"smlcloudplatform/pkg/shop"
	"smlcloudplatform/pkg/shop/employee"
	"smlcloudplatform/pkg/shopdesign/zonedesign"
	"smlcloudplatform/pkg/smsreceive/smstransaction"
	"smlcloudplatform/pkg/sysinfo"
	"smlcloudplatform/pkg/task"
	"smlcloudplatform/pkg/transaction/paid"
	"smlcloudplatform/pkg/transaction/pay"
	"smlcloudplatform/pkg/transaction/purchase"
	"smlcloudplatform/pkg/transaction/purchasereturn"
	"smlcloudplatform/pkg/transaction/saleinvoice"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn"
	"smlcloudplatform/pkg/transaction/smltransaction"
	"smlcloudplatform/pkg/transaction/stockadjustment"
	"smlcloudplatform/pkg/transaction/stockpickupproduct"
	"smlcloudplatform/pkg/transaction/stockreceiveproduct"
	"smlcloudplatform/pkg/transaction/stockreturnproduct"
	"smlcloudplatform/pkg/transaction/stocktransfer"
	"smlcloudplatform/pkg/vfgl/accountgroup"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster"
	"smlcloudplatform/pkg/vfgl/chartofaccount"
	"smlcloudplatform/pkg/vfgl/journal"
	"smlcloudplatform/pkg/vfgl/journalbook"
	"smlcloudplatform/pkg/vfgl/journalreport"
	"smlcloudplatform/pkg/warehouse"

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

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	ms.HttpUsePrometheus()

	if devApiMode == "" || devApiMode == "2" {

		ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

		cacher := ms.Cacher(cfg.CacherConfig())
		authService := microservice.NewAuthService(cacher, 24*3)
		publicPath := []string{
			"/swagger",
			"/login",
			"/tokenlogin",
			"/register",

			"/employee/login",

			"/images*",
			"/productimage",

			"/healthz",
			"/ws",
			"/metrics",
			"/e-order/product",
			"/e-order/category",
			"/e-order/product-barcode",
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

		httpServices := []HttpRouteSetup{

			apikeyservice.NewApiKeyServiceHttp(ms, cfg),
			authentication.NewAuthenticationHttp(ms, cfg),
			apikeyservice.NewApiKeyServiceHttp(ms, cfg),
			shop.NewShopHttp(ms, cfg),

			shop.NewShopMemberHttp(ms, cfg),
			employee.NewEmployeeHttp(ms, cfg), member.NewMemberHttp(ms, cfg),

			inventory.NewInventoryHttp(ms, cfg),
			option.NewOptionHttp(ms, cfg),
			unit.NewUnitHttp(ms, cfg),
			optionpattern.NewOptionPatternHttp(ms, cfg),
			color.NewColorHttp(ms, cfg),

			//product
			productcategory.NewProductCategoryHttp(ms, cfg),
			productbarcode.NewProductBarcodeHttp(ms, cfg),
			product.NewProductHttp(ms, cfg),
			productgroup.NewProductGroupHttp(ms, cfg),

			inventoryimport.NewInventoryImportHttp(ms, cfg),
			inventoryimport.NewInventoryImporOptionMaintHttp(ms, cfg),
			inventoryimport.NewCategoryImportHttp(ms, cfg),

			images.NewImagesHttp(ms, cfg, imagePersister),

			// restaurant
			shopzone.NewShopZoneHttp(ms, cfg),
			shoptable.NewShopTableHttp(ms, cfg),
			printer.NewPrinterHttp(ms, cfg),
			kitchen.NewKitchenHttp(ms, cfg),
			zonedesign.NewZoneDesignHttp(ms, cfg),
			restaurantsettings.NewRestaurantSettingsHttp(ms, cfg),
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
			sysinfo.NewSysInfoHttp(ms, cfg),

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

			//product section
			sectionbranch.NewSectionBranchHttp(ms, cfg),
			sectiondepartment.NewSectionDepartmentHttp(ms, cfg),
			sectionbusinesstype.NewSectionBusinessTypeHttp(ms, cfg),

			//channel
			salechannel.NewSaleChannelHttp(ms, cfg),
			transportchannel.NewTransportChannelHttp(ms, cfg),

			// e-order
			eorder.NewEOrderHttp(ms, cfg),
		}

		startHttpServices(httpServices...)

	}

	if devApiMode == "1" || devApiMode == "2" {

		ms.RegisterLivenessProbeEndpoint("/healthz")

		consumerGroupName := os.Getenv("CONSUMER_GROUP_NAME")
		if consumerGroupName == "" {
			consumerGroupName = "03"
		}

		inventoryConsumer := inventorysearchconsumer.NewInventorySearchConsumer(ms, cfg)
		inventoryConsumer.Start()

		journal.MigrationJournalTable(ms, cfg)
		journal.StartJournalComsumeCreated(ms, cfg, consumerGroupName)
		journal.StartJournalComsumeUpdated(ms, cfg, consumerGroupName)
		journal.StartJournalComsumeDeleted(ms, cfg, consumerGroupName)
		journal.StartJournalComsumeBlukCreated(ms, cfg, consumerGroupName)

		chartofaccount.MigrationChartOfAccountTable(ms, cfg)
		chartofaccount.StartChartOfAccountConsumerCreated(ms, cfg, consumerGroupName)
		chartofaccount.StartChartOfAccountConsumerUpdated(ms, cfg, consumerGroupName)
		chartofaccount.StartChartOfAccountConsumerDeleted(ms, cfg, consumerGroupName)
		chartofaccount.StartChartOfAccountConsumerBlukCreated(ms, cfg, consumerGroupName)

		task.NewTaskConsumer(ms, cfg).RegisterConsumer()

	}

	ms.Start()
}

type HttpRouteSetup interface {
	RouteSetup()
}

func startHttpServices(services ...HttpRouteSetup) {
	for _, service := range services {
		service.RouteSetup()
	}
}
