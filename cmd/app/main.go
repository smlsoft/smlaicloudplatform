package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/authentication"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/api/merchant"
	"smlcloudplatform/pkg/api/purchase"
	"smlcloudplatform/pkg/api/tools"
	"smlcloudplatform/pkg/api/transaction"

	"github.com/labstack/echo/v4"
)

func main() {

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	cacher := ms.Cacher(cfg.CacherConfig())
	// jwtService := microservice.NewJwtService(cacher, cfg.JwtSecretKey(), 24*3)
	authService := microservice.NewAuthService(cacher, 24*3)

	publicPath := []string{
		"/login",
		"/register",
		"/list-merchant",
		"/select-merchant",
		"/healthz",
	}

	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher, publicPath...))

	ms.RegisterLivenessProbeEndpoint("/healthz")

	svcAuth := authentication.NewAuthenticationHttp(ms, cfg)
	svcAuth.RouteSetup()

	svcMerchant := merchant.NewMerchantHttp(ms, cfg)
	svcMerchant.RouteSetup()

	inventoryapi := inventory.NewInventoryHttp(ms, cfg)
	inventoryapi.RouteSetup()

	transapi := transaction.NewTransactionHttp(ms, cfg)
	transapi.RouteSetup()

	transaction.StartTransactionAsync(ms, cfg)

	purchaseapi := purchase.NewPurchaseHttp(ms, cfg)
	purchaseapi.RouteSetup()

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
