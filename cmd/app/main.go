package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	"smlaicloudplatform/internal/member"
	"smlaicloudplatform/internal/notify"
	"smlaicloudplatform/internal/ocr"
	"smlaicloudplatform/internal/organization/branch"
	"smlaicloudplatform/internal/organization/businesstype"
	"smlaicloudplatform/internal/organization/department"
	"smlaicloudplatform/internal/payment/bankmaster"
	"smlaicloudplatform/internal/payment/bookbank"
	"smlaicloudplatform/internal/payment/qrpayment"
	"smlaicloudplatform/internal/paymentmaster"
	"smlaicloudplatform/internal/product/bom"
	"smlaicloudplatform/internal/product/color"
	"smlaicloudplatform/internal/product/eorder"
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
	"smlaicloudplatform/internal/report/reportquerym"
	"smlaicloudplatform/internal/slipimage"
	"smlaicloudplatform/internal/stockbalanceimport"
	"smlaicloudplatform/pkg/microservice"
	"time"

	// "smlaicloudplatform/internal/report/reportquery"
	"smlaicloudplatform/internal/restaurant/device"
	"smlaicloudplatform/internal/restaurant/kitchen"
	"smlaicloudplatform/internal/restaurant/notifier"
	"smlaicloudplatform/internal/restaurant/notifierdevice"
	"smlaicloudplatform/internal/restaurant/printer"
	"smlaicloudplatform/internal/restaurant/settings"
	"smlaicloudplatform/internal/restaurant/staff"
	"smlaicloudplatform/internal/restaurant/table"
	"smlaicloudplatform/internal/restaurant/zone"
	"smlaicloudplatform/internal/shop"

	// "smlaicloudplatform/internal/shop/branch"
	order_device "smlaicloudplatform/internal/order/device"
	order_setting "smlaicloudplatform/internal/order/setting"
	"smlaicloudplatform/internal/pos/media"
	pos_setting "smlaicloudplatform/internal/pos/setting"
	"smlaicloudplatform/internal/pos/shift"
	"smlaicloudplatform/internal/pos/temp"
	"smlaicloudplatform/internal/shop/employee"
	"smlaicloudplatform/internal/shopdesign/zonedesign"
	"smlaicloudplatform/internal/smsreceive/smspatterns"
	"smlaicloudplatform/internal/smsreceive/smspaymentsettings"
	"smlaicloudplatform/internal/smsreceive/smstransaction"
	"smlaicloudplatform/internal/storefront"
	"smlaicloudplatform/internal/task"
	"smlaicloudplatform/internal/tools"
	"smlaicloudplatform/internal/transaction/documentformate"
	"smlaicloudplatform/internal/transaction/paid"
	"smlaicloudplatform/internal/transaction/pay"
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

	purchase_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/purchase"
	saleinvoice_consumer "smlaicloudplatform/internal/transaction/transactionconsumer/saleinvoice"

	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
)

func main() {

	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	cacher := ms.Cacher(cfg.CacherConfig())
	// jwtService := microservice.NewJwtService(cacher, cfg.JwtSecretKey(), 24*3)
	authService := microservice.NewAuthService(cacher, 24*3*time.Hour, 24*30*time.Hour)

	publicPath := []string{
		"/slip/*",
		"/swagger",
		"/login",
		"/poslogin",
		"/login/phone-number",
		"/register",
		"/refresh",
		"/register-phonenumber",
		"/register/exists-username",
		"/register/exists-phonenumber",
		"/send-phonenumber-otp",

		"/employee/login",

		"/images*",
		"/productimage",

		"/healthz",
		"/ws",
		"/metrics",

		"/e-order/product",
		"/e-order/category",
		"/e-order/product-barcode",
		"/e-order/shop-info",
		"/e-order/shop-info/v1.1",
		"/e-order/sale-invoice/last-pos-docno",

		"/e-order/restaurant/zone",
		"/e-order/restaurant/kitchen",
		"/e-order/restaurant/table",
		"/e-order/notify",
		"/line-notify",

		"/restaurant/notifier-device/ref-confirm",
	}

	exceptShopPath := []string{
		"/shop",
		"/list-shop",
		"/select-shop",
		"/create-shop",
	}

	ms.HttpPreRemoveTrailingSlash()
	// ms.HttpUsePrometheus()
	// ms.HttpUseJaeger()

	// ms.HttpMiddleware(authService.MWFuncWithRedis(cacher, publicPath...))
	ms.HttpMiddleware(authService.MWFuncWithRedisMixShop(cacher, exceptShopPath, publicPath...))

	ms.RegisterLivenessProbeEndpoint("/healthz")

	httpServices := []HttpRegister{
		authentication.NewAuthenticationHttp(ms, cfg),
		shop.NewShopHttp(ms, cfg),
		shop.NewShopMemberHttp(ms, cfg),
		member.NewMemberHttp(ms, cfg),
		employee.NewEmployeeHttp(ms, cfg),

		//restaurants
		zone.NewZoneHttp(ms, cfg),
		table.NewTableHttp(ms, cfg),
		printer.NewPrinterHttp(ms, cfg),
		kitchen.NewKitchenHttp(ms, cfg),
		settings.NewRestaurantSettingsHttp(ms, cfg),
		device.NewDeviceHttp(ms, cfg),
		staff.NewStaffHttp(ms, cfg),
		notifier.NewNotifierHttp(ms, cfg),
		notifierdevice.NewNotifierDeviceHttp(ms, cfg),

		//Journal
		journal.NewJournalHttp(ms, cfg),
		journal.NewJournalWs(ms, cfg),
		accountgroup.NewAccountGroupHttp(ms, cfg),
		journalbook.NewJournalBookHttp(ms, cfg),
		zonedesign.NewZoneDesignHttp(ms, cfg),

		mastersync.NewMasterSyncHttp(ms, cfg),

		documentimage.NewDocumentImageHttp(ms, cfg),
		chartofaccount.NewChartOfAccountHttp(ms, cfg),
		//new

		paymentmaster.NewPaymentMasterHttp(ms, cfg),
		apikeyservice.NewApiKeyServiceHttp(ms, cfg),

		smstransaction.NewSmsTransactionHttp(ms, cfg),
		smspatterns.NewSmsPatternsHttp(ms, cfg),
		smspaymentsettings.NewSmsPaymentSettingsHttp(ms, cfg),

		warehouse.NewWarehouseHttp(ms, cfg),
		storefront.NewStorefrontHttp(ms, cfg),

		unit.NewUnitHttp(ms, cfg),
		journalreport.NewJournalReportHttp(ms, cfg),

		optionpattern.NewOptionPatternHttp(ms, cfg),
		color.NewColorHttp(ms, cfg),
		productcategory.NewProductCategoryHttp(ms, cfg),
		productbarcode.NewProductBarcodeHttp(ms, cfg),
		productgroup.NewProductGroupHttp(ms, cfg),

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

		department.NewDepartmentHttp(ms, cfg),
		businesstype.NewBusinessTypeHttp(ms, cfg),
		branch.NewBranchHttp(ms, cfg),

		//transaction
		purchase.NewPurchaseHttp(ms, cfg),
		purchasereturn.NewPurchaseReturnHttp(ms, cfg),
		saleinvoice.NewSaleInvoiceHttp(ms, cfg),
		saleinvoicereturn.NewSaleInvoiceReturnHttp(ms, cfg),
		stockreceiveproduct.NewStockReceiveProductHttp(ms, cfg),
		stockreturnproduct.NewStockReturnProductHttp(ms, cfg),
		stockpickupproduct.NewStockPickupProductHttp(ms, cfg),
		stockadjustment.NewStockAdjustmentHttp(ms, cfg),
		stocktransfer.NewStockTransferHttp(ms, cfg),
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

		paid.NewPaidHttp(ms, cfg),
		pay.NewPayHttp(ms, cfg),

		promotion.NewPromotionHttp(ms, cfg),

		eorder.NewEOrderHttp(ms, cfg),

		ordertype.NewOrderTypeHttp(ms, cfg),

		// POS
		pos_setting.NewSettingHttp(ms, cfg),
		shift.NewShiftHttp(ms, cfg),
		order_setting.NewSettingHttp(ms, cfg),
		order_device.NewDeviceHttp(ms, cfg),

		documentformate.NewDocumentFormateHttp(ms, cfg),

		//reportquery.NewReportQueryHttp(ms, cfg),
		reportquerym.NewReportQueryHttp(ms, cfg),
		producttype.NewProductTypeHttp(ms, cfg),

		media.NewMediaHttp(ms, cfg),

		ocr.NewOcrHttp(ms, cfg),
		stockbalanceimport.NewStockBalanceImportHttp(ms, cfg),
		productbarcode.NewProductBarcodeHttp(ms, cfg),
		notify.NewNotifyHttp(ms, cfg),
		slipimage.NewSlipImageHttp(ms, cfg),
		productimport.NewProductImportHttp(ms, cfg),

		dimension.NewDimensionHttp(ms, cfg),

		// master
		masterexpense.NewMasterExpenseHttp(ms, cfg),
		masterincome.NewMasterIncomeHttp(ms, cfg),

		temp.NewPOSTempHttp(ms, cfg),
		filestatus.NewFileStatusHttp(ms, cfg),

		member.NewMemberHttp(ms, cfg),

		bom.NewBOMHttp(ms, cfg),
		saleinvoicebomprice.NewSaleInvoiceBomPriceHttp(ms, cfg),
	}

	azureFileBlob := microservice.NewPersisterAzureBlob()
	imagePersister := microservice.NewPersisterImage(azureFileBlob)

	ms.RegisterHttp(images.NewImagesHttp(ms, cfg, imagePersister))

	serviceStartHttp(ms, httpServices...)

	// journal.MigrationJournalTable(ms, cfg)

	// warehouse.MigrationDatabase(ms, cfg)
	// creditor.MigrationDatabase(ms, cfg)
	// ms.RegisterConsumer(creditor.InitCreditorConsumer(ms, cfg))

	debtor.MigrationDatabase(ms, cfg)
	ms.RegisterConsumer(debtor.InitDebtorConsumer(ms, cfg))
	// warehouse.MigrationDatabase(ms, cfg)

	// inventory.StartInventoryAsync(ms, cfg)
	// inventory.StartInventoryComsumeCreated(ms, cfg)

	// transactionconsumer.MigrationDatabase(ms, cfg)

	// payment.MigrationDatabase(ms, cfg)
	// paymentdetail.MigrationDatabase(ms, cfg)

	// purchase_consumer.MigrationDatabase(ms, cfg)
	// saleinvoice_consumer.MigrationDatabase(ms, cfg)

	// productbarcode.MigrationDatabase(ms, cfg)

	ms.RegisterConsumer(purchase_consumer.InitPurchaseTransactionConsumer(ms, cfg))

	ms.RegisterConsumer(saleinvoice_consumer.InitSaleInvoiceTransactionConsumer(ms, cfg))

	ms.RegisterConsumer(journal.InitJournalTransactionConsumer(ms, cfg))

	ms.RegisterConsumer(warehouse.InitWarehouseConsumer(ms, cfg))

	consumeServices := []ConsumerRegister{}

	task.NewTaskConsumer(ms, cfg).RegisterConsumer()
	productbarcode.NewProductBarcodeConsumer(ms, cfg).RegisterConsumer()

	serviceStartConsumer(ms, consumeServices...)

	toolSvc := tools.NewToolsService(ms, cfg)

	toolSvc.RegisterHttp()

	ms.Echo().GET("/routes", func(ctx echo.Context) error {
		data, err := json.MarshalIndent(ms.Echo().Routes(), "", "  ")

		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		os.WriteFile("routes.json", data, 0644)

		// ctx.JSON(http.StatusOK, data)
		ctx.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": ms.Echo().Routes()})

		return nil
	})

	//ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

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
	RegisterConsumer(*microservice.Microservice)
}

func serviceStartConsumer(ms *microservice.Microservice, services ...ConsumerRegister) {
	for _, service := range services {
		ms.RegisterConsumer(service)
	}
}
