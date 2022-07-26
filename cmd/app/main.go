package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/authentication"
	"smlcloudplatform/pkg/mastersync"
	"smlcloudplatform/pkg/member"
	"smlcloudplatform/pkg/product/category"
	"smlcloudplatform/pkg/product/inventory"
	"smlcloudplatform/pkg/product/inventoryimport"
	"smlcloudplatform/pkg/restaurant/kitchen"
	"smlcloudplatform/pkg/restaurant/shopprinter"
	"smlcloudplatform/pkg/restaurant/shoptable"
	"smlcloudplatform/pkg/restaurant/shopzone"
	"smlcloudplatform/pkg/shop"
	"smlcloudplatform/pkg/shop/employee"
	"smlcloudplatform/pkg/shopdesign/zonedesign"
	"smlcloudplatform/pkg/syncdata"
	"smlcloudplatform/pkg/tools"
	"smlcloudplatform/pkg/transaction/purchase"
	"smlcloudplatform/pkg/transaction/saleinvoice"
	"smlcloudplatform/pkg/vfgl/accountgroup"
	"smlcloudplatform/pkg/vfgl/journal"
	"smlcloudplatform/pkg/vfgl/journalbook"

	// _ "net/http/pprof"

	"github.com/labstack/echo/v4"
)

func main() {

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	// go func() {
	// 	log.Println(http.ListenAndServe(":6060", nil))
	// }()

	cacher := ms.Cacher(cfg.CacherConfig())
	// jwtService := microservice.NewJwtService(cacher, cfg.JwtSecretKey(), 24*3)
	authService := microservice.NewAuthService(cacher, 24*3)

	publicPath := []string{
		"/login",
		"/register",
		"/list-shop",
		"/select-shop",
		"/create-shop",
		"/employee/login",
		"/healthz",
		"/metrics",
	}

	ms.HttpPreRemoveTrailingSlash()
	ms.HttpUsePrometheus()

	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher, publicPath...))

	ms.RegisterLivenessProbeEndpoint("/healthz")

	authHttp := authentication.NewAuthenticationHttp(ms, cfg)
	authHttp.RouteSetup()

	shopHttp := shop.NewShopHttp(ms, cfg)
	shopHttp.RouteSetup()

	categoryHttp := category.NewCategoryHttp(ms, cfg)
	categoryHttp.RouteSetup()

	inventoryHttp := inventory.NewInventoryHttp(ms, cfg)
	inventoryHttp.RouteSetup()

	inventory.StartInventoryAsync(ms, cfg)
	inventory.StartInventoryComsumeCreated(ms, cfg)

	saleinvoiceHttp := saleinvoice.NewSaleinvoiceHttp(ms, cfg)
	saleinvoiceHttp.RouteSetup()

	saleinvoice.StartSaleinvoiceAsync(ms, cfg)

	purchaseHttp := purchase.NewPurchaseHttp(ms, cfg)
	purchaseHttp.RouteSetup()

	purchase.StartPurchaseAsync(ms, cfg)

	syncDataHttp := syncdata.NewSyncDataHttp(ms, cfg)
	syncDataHttp.RouteSetup()

	memberHttp := member.NewMemberHttp(ms, cfg)
	memberHttp.RouteSetup()

	emphttp := employee.NewEmployeeHttp(ms, cfg)
	emphttp.RouteSetup()

	inventoryImportHttp := inventoryimport.NewInventoryImportHttp(ms, cfg)
	inventoryImportHttp.RouteSetup()

	inventoryOptionImportHttp := inventoryimport.NewInventoryImporOptionMaintHttp(ms, cfg)
	inventoryOptionImportHttp.RouteSetup()

	categoryImportHttp := inventoryimport.NewCategoryImportHttp(ms, cfg)
	categoryImportHttp.RouteSetup()

	shopzonehttp := shopzone.NewShopZoneHttp(ms, cfg)
	shopzonehttp.RouteSetup()

	shoptablehttp := shoptable.NewShopTableHttp(ms, cfg)
	shoptablehttp.RouteSetup()

	shopprinterhttp := shopprinter.NewShopPrinterHttp(ms, cfg)
	shopprinterhttp.RouteSetup()

	kitchenhttp := kitchen.NewKitchenHttp(ms, cfg)
	kitchenhttp.RouteSetup()

	journalhttp := journal.NewJournalHttp(ms, cfg)
	journalhttp.RouteSetup()

	journalWs := journal.NewJournalWs(ms, cfg)
	journalWs.RouteSetup()

	accountGroupHttp := accountgroup.NewAccountGroupHttp(ms, cfg)
	accountGroupHttp.RouteSetup()

	journalBookhttp := journalbook.NewJournalBookHttp(ms, cfg)
	journalBookhttp.RouteSetup()
	zonedesignhttp := zonedesign.NewZoneDesignHttp(ms, cfg)
	zonedesignhttp.RouteSetup()

	mastersync := mastersync.NewMasterSyncHttp(ms, cfg)
	mastersync.RouteSetup()

	toolSvc := tools.NewToolsService(ms, cfg)

	toolSvc.RouteSetup()

	ms.Echo().GET("/routes", func(ctx echo.Context) error {
		data, err := json.MarshalIndent(ms.Echo().Routes(), "", "  ")

		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		ioutil.WriteFile("routes.json", data, 0644)

		// ctx.JSON(http.StatusOK, data)
		ctx.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": ms.Echo().Routes()})

		return nil
	})

	//ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	ms.Start()
}
