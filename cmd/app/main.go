package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/authentication"
	"smlcloudplatform/pkg/documentwarehouse/documentimage"
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

	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
)

func main() {

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

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
	ms.HttpUseJaeger()

	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher, publicPath...))

	ms.RegisterLivenessProbeEndpoint("/healthz")

	services := []HttpRouteSetup{
		authentication.NewAuthenticationHttp(ms, cfg),
		shop.NewShopHttp(ms, cfg),
		category.NewCategoryHttp(ms, cfg),
		inventory.NewInventoryHttp(ms, cfg),
		saleinvoice.NewSaleinvoiceHttp(ms, cfg),
		purchase.NewPurchaseHttp(ms, cfg),
		syncdata.NewSyncDataHttp(ms, cfg),
		member.NewMemberHttp(ms, cfg),
		employee.NewEmployeeHttp(ms, cfg),
		inventoryimport.NewInventoryImportHttp(ms, cfg),
		inventoryimport.NewInventoryImporOptionMaintHttp(ms, cfg),
		inventoryimport.NewCategoryImportHttp(ms, cfg),
		shopzone.NewShopZoneHttp(ms, cfg),
		shoptable.NewShopTableHttp(ms, cfg),
		shopprinter.NewShopPrinterHttp(ms, cfg),
		kitchen.NewKitchenHttp(ms, cfg),
		//Journal
		journal.NewJournalHttp(ms, cfg),
		journal.NewJournalWs(ms, cfg),
		accountgroup.NewAccountGroupHttp(ms, cfg),
		journalbook.NewJournalBookHttp(ms, cfg),
		zonedesign.NewZoneDesignHttp(ms, cfg),
		mastersync.NewMasterSyncHttp(ms, cfg),
		documentimage.NewDocumentImageHttp(ms, cfg),
	}

	serviceStartHttp(services...)

	inventory.StartInventoryAsync(ms, cfg)
	inventory.StartInventoryComsumeCreated(ms, cfg)

	saleinvoice.StartSaleinvoiceAsync(ms, cfg)
	purchase.StartPurchaseAsync(ms, cfg)

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

type HttpRouteSetup interface {
	RouteSetup()
}

func serviceStartHttp(services ...HttpRouteSetup) {
	for _, service := range services {
		service.RouteSetup()
	}
}
