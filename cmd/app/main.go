package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/apikeyservice"
	"smlcloudplatform/pkg/authentication"
	"smlcloudplatform/pkg/channel/salechannel"
	"smlcloudplatform/pkg/channel/transportchannel"
	"smlcloudplatform/pkg/config"
	"smlcloudplatform/pkg/debtaccount/creditor"
	"smlcloudplatform/pkg/debtaccount/creditorgroup"
	"smlcloudplatform/pkg/debtaccount/customer"
	"smlcloudplatform/pkg/debtaccount/customergroup"
	"smlcloudplatform/pkg/debtaccount/debtor"
	"smlcloudplatform/pkg/debtaccount/debtorgroup"
	"smlcloudplatform/pkg/documentwarehouse/documentimage"
	"smlcloudplatform/pkg/mastersync"
	"smlcloudplatform/pkg/member"
	"smlcloudplatform/pkg/notify"
	"smlcloudplatform/pkg/ocr"
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
	"smlcloudplatform/pkg/product/optionpattern"
	"smlcloudplatform/pkg/product/ordertype"
	"smlcloudplatform/pkg/product/product"
	"smlcloudplatform/pkg/product/productbarcode"
	"smlcloudplatform/pkg/product/productcategory"
	"smlcloudplatform/pkg/product/productgroup"
	"smlcloudplatform/pkg/product/producttype"
	"smlcloudplatform/pkg/product/promotion"
	"smlcloudplatform/pkg/product/unit"
	"smlcloudplatform/pkg/productsection/sectionbranch"
	"smlcloudplatform/pkg/productsection/sectionbusinesstype"
	"smlcloudplatform/pkg/productsection/sectiondepartment"
	"smlcloudplatform/pkg/report/reportquerym"
	"smlcloudplatform/pkg/stockbalanceimport"
	"time"

	// "smlcloudplatform/pkg/report/reportquery"
	"smlcloudplatform/pkg/restaurant/device"
	"smlcloudplatform/pkg/restaurant/kitchen"
	"smlcloudplatform/pkg/restaurant/notifier"
	"smlcloudplatform/pkg/restaurant/notifierdevice"
	"smlcloudplatform/pkg/restaurant/printer"
	"smlcloudplatform/pkg/restaurant/settings"
	"smlcloudplatform/pkg/restaurant/staff"
	"smlcloudplatform/pkg/restaurant/table"
	"smlcloudplatform/pkg/restaurant/zone"
	"smlcloudplatform/pkg/shop"

	// "smlcloudplatform/pkg/shop/branch"
	order_device "smlcloudplatform/pkg/order/device"
	order_setting "smlcloudplatform/pkg/order/setting"
	"smlcloudplatform/pkg/pos/media"
	pos_setting "smlcloudplatform/pkg/pos/setting"
	"smlcloudplatform/pkg/pos/shift"
	"smlcloudplatform/pkg/shop/employee"
	"smlcloudplatform/pkg/shopdesign/zonedesign"
	"smlcloudplatform/pkg/smsreceive/smspatterns"
	"smlcloudplatform/pkg/smsreceive/smspaymentsettings"
	"smlcloudplatform/pkg/smsreceive/smstransaction"
	"smlcloudplatform/pkg/storefront"
	"smlcloudplatform/pkg/syncdata"
	"smlcloudplatform/pkg/sysinfo"
	"smlcloudplatform/pkg/task"
	"smlcloudplatform/pkg/tools"
	"smlcloudplatform/pkg/transaction/documentformate"
	"smlcloudplatform/pkg/transaction/paid"
	"smlcloudplatform/pkg/transaction/pay"
	"smlcloudplatform/pkg/transaction/purchase"
	"smlcloudplatform/pkg/transaction/purchasereturn"
	"smlcloudplatform/pkg/transaction/saleinvoice"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn"
	"smlcloudplatform/pkg/transaction/smltransaction"
	"smlcloudplatform/pkg/transaction/stockadjustment"
	"smlcloudplatform/pkg/transaction/stockbalance"
	stockbalancedetaildetail "smlcloudplatform/pkg/transaction/stockbalancedetail"
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
		"/swagger",
		"/login",
		"/register",
		"/refresh",

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
		inventory.NewInventoryHttp(ms, cfg),

		syncdata.NewSyncDataHttp(ms, cfg),
		member.NewMemberHttp(ms, cfg),
		employee.NewEmployeeHttp(ms, cfg),
		inventoryimport.NewInventoryImportHttp(ms, cfg),
		inventoryimport.NewInventoryImporOptionMaintHttp(ms, cfg),
		inventoryimport.NewCategoryImportHttp(ms, cfg),

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
		product.NewProductHttp(ms, cfg),
		productgroup.NewProductGroupHttp(ms, cfg),

		accountperiodmaster.NewAccountPeriodMasterHttp(ms, cfg),

		bankmaster.NewBankMasterHttp(ms, cfg),
		bookbank.NewBookBankHttp(ms, cfg),
		qrpayment.NewQrPaymentHttp(ms, cfg),

		task.NewTaskHttp(ms, cfg),
		smltransaction.NewSMLTransactionHttp(ms, cfg),

		sysinfo.NewSysInfoHttp(ms, cfg),
		// branch.NewBranchHttp(ms, cfg),

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
		stockbalancedetaildetail.NewStockBalanceDetailHttp(ms, cfg),

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
		notify.NewNotifyHttp(ms, cfg),
	}

	serviceStartHttp(ms, httpServices...)

	// inventory.StartInventoryAsync(ms, cfg)
	// inventory.StartInventoryComsumeCreated(ms, cfg)

	consumeServices := []ConsumerRegister{
		task.NewTaskConsumer(ms, cfg),
		productbarcode.NewProductBarcodeConsumer(ms, cfg),
	}

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
	RegisterConsumer()
}

func serviceStartConsumer(ms *microservice.Microservice, services ...ConsumerRegister) {
	for _, service := range services {
		ms.RegisterConsumer(service)
	}
}
