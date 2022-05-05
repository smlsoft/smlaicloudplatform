package main

import (
	"fmt"
	"os"
	"smlcloudplatform/docs"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/authentication"
	"smlcloudplatform/pkg/api/category"
	"smlcloudplatform/pkg/api/images"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/api/inventoryimport"
	"smlcloudplatform/pkg/api/shop"

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

	// old swagger run
	// fmt.Println("Start Swagger API")

	// host := os.Getenv("HOST_API")
	// if host != "" {
	// 	fmt.Printf("Host: %v\n", host)
	// 	docs.SwaggerInfo.Host = host
	// }

	// e := echo.New()

	// e.GET("/swagger/*", echoSwagger.WrapHandler)

	// serverPort := os.Getenv("SERVER_PORT")
	// if serverPort == "" {
	// 	serverPort = "1323"
	// }
	// e.Logger.Fatal(e.Start(":" + serverPort))

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

	ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	cacher := ms.Cacher(cfg.CacherConfig())
	authService := microservice.NewAuthService(cacher, 24*3)
	publicPath := []string{
		"/swagger",
		"/login",
		"/register",
		"/list-shop",
		"/select-shop",
		"/create-shop",
		"/employee/login",

		"/images*",
		"/productimage",

		"/healthz",
	}
	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher, publicPath...))
	ms.RegisterLivenessProbeEndpoint("/healthz")

	authHttp := authentication.NewAuthenticationHttp(ms, cfg)
	authHttp.RouteSetup()

	shopHttp := shop.NewShopHttp(ms, cfg)
	shopHttp.RouteSetup()

	inventoryapi := inventory.NewInventoryHttp(ms, cfg)
	inventoryapi.RouteSetup()

	categoryHttp := category.NewCategoryHttp(ms, cfg)
	categoryHttp.RouteSetup()

	invImp := inventoryimport.NewInventoryImportHttp(ms, cfg)
	invImp.RouteSetup()

	invOptionImp := inventoryimport.NewInventoryImporOptionMaintHttp(ms, cfg)
	invOptionImp.RouteSetup()

	catImp := inventoryimport.NewCategoryImportHttp(ms, cfg)
	catImp.RouteSetup()

	filePersister := microservice.NewPersisterFile(microservice.NewStorageFileConfig())
	imagePersister := microservice.NewPersisterImage(filePersister)

	ms.Logger.Debugf("Store File Path %v", filePersister.StoreFilePath)
	imageHttp := images.NewImagesHttp(ms, cfg, imagePersister)
	imageHttp.RouteSetup()

	ms.Start()
}
