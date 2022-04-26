package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/authentication"
	"smlcloudplatform/pkg/api/category"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/api/inventoryimport"
	"smlcloudplatform/pkg/api/member"
	"smlcloudplatform/pkg/api/purchase"
	"smlcloudplatform/pkg/api/shop"
	"smlcloudplatform/pkg/api/shop/employee"
	"smlcloudplatform/pkg/api/syncdata"
	"smlcloudplatform/pkg/api/tools"
	"smlcloudplatform/pkg/api/transaction"

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
	}

	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher, publicPath...))

	ms.RegisterLivenessProbeEndpoint("/healthz")

	svcAuth := authentication.NewAuthenticationHttp(ms, cfg)
	svcAuth.RouteSetup()

	svcShop := shop.NewShopHttp(ms, cfg)
	svcShop.RouteSetup()

	categoryHttp := category.NewCategoryHttp(ms, cfg)
	categoryHttp.RouteSetup()

	inventoryapi := inventory.NewInventoryHttp(ms, cfg)
	inventoryapi.RouteSetup()

	inventory.StartInventoryAsync(ms, cfg)
	inventory.StartInventoryComsumeCreated(ms, cfg)

	transapi := transaction.NewTransactionHttp(ms, cfg)
	transapi.RouteSetup()

	transaction.StartTransactionAsync(ms, cfg)

	purchaseapi := purchase.NewPurchaseHttp(ms, cfg)
	purchaseapi.RouteSetup()

	purchase.StartPurchaseAsync(ms, cfg)

	syncapi := syncdata.NewSyncDataHttp(ms, cfg)
	syncapi.RouteSetup()

	memberhttp := member.NewMemberHttp(ms, cfg)
	memberhttp.RouteSetup()

	emphttp := employee.NewEmployeeHttp(ms, cfg)
	emphttp.RouteSetup()

	invImp := inventoryimport.NewInventoryImportHttp(ms, cfg)
	invImp.RouteSetup()

	invOptionImp := inventoryimport.NewInventoryImporOptionMaintHttp(ms, cfg)
	invOptionImp.RouteSetup()

	catImp := inventoryimport.NewCategoryImportHttp(ms, cfg)
	catImp.RouteSetup()

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
