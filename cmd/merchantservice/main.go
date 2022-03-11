package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/merchant"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	cacher := ms.Cacher(cfg.CacherConfig())
	authService := microservice.NewAuthService(cacher, 24*3)

	ms.HttpMiddleware(authService.MWFuncWithMerchant(cacher))
	svc := merchant.NewMerchantHttp(ms, cfg)
	svc.RouteSetup()

	//ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	ms.Start()
}
