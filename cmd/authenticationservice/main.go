package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/authentication"
	"smlcloudplatform/pkg/config"
)

func main() {

	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}
	svc := authentication.NewAuthenticationHttp(ms, cfg)

	cacher := ms.Cacher(cfg.CacherConfig())
	// jwtService := microservice.NewJwtService(cacher, cfg.JwtSecretKey(), 24*3)
	authService := microservice.NewAuthService(cacher, 24*3)

	publicPath := []string{
		"/login",
		"/register",
		"/list-shop",
		"/select-shop",
		"/create-shop",
		"/healthz",
	}

	ms.HttpMiddleware(authService.MWFuncWithShop(cacher, publicPath...))
	svc.RouteSetup()

	// ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	ms.Start()
}
