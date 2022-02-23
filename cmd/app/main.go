package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api"
	"smlcloudplatform/pkg/api/inventoryservice"
	"smlcloudplatform/pkg/api/merchantservice"
	"smlcloudplatform/pkg/api/toolsservice"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := microservice.NewConfig()
	ms := microservice.NewMicroservice(cfg)

	cacher := ms.Cacher(cfg.CacherConfig())
	jwtService := microservice.NewJwtService(cacher, cfg.JwtSecretKey(), 24*3)

	publicPath := []string{
		"/login",
		"/register",
		"/select-merchant",
	}

	ms.HttpMiddleware(jwtService.MWFuncWithRedis(cacher, publicPath...))

	svcAuth := api.NewAuthenticationService(ms, cfg)
	svcAuth.RouteSetup()

	svcMerchant := merchantservice.NewMerchantHttp(ms, cfg)
	svcMerchant.RouteSetup()

	inventoryapi := inventoryservice.NewInventoryService(ms, cfg)
	inventoryapi.RouteSetup()

	toolSvc := toolsservice.NewToolsService(ms, cfg)

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
