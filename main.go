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
	"smlcloudplatform/pkg/api/inventorysearchconsumer"
	"smlcloudplatform/pkg/api/member"
	"smlcloudplatform/pkg/api/purchase"
	"smlcloudplatform/pkg/api/restaurant/kitchen"
	"smlcloudplatform/pkg/api/restaurant/shopprinter"
	"smlcloudplatform/pkg/api/restaurant/shoptable"
	"smlcloudplatform/pkg/api/restaurant/shopzone"
	"smlcloudplatform/pkg/api/saleinvoice"
	"smlcloudplatform/pkg/api/shop"
	"smlcloudplatform/pkg/api/shop/employee"
	"smlcloudplatform/pkg/vfgl/chartofaccount"
	"smlcloudplatform/pkg/vfgl/journal"
	"smlcloudplatform/pkg/vfgl/journalreport"

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
	devApiMode := os.Getenv("DEV_API_MODE")

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

	if devApiMode == "" || devApiMode == 2 {

		ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

		cacher := ms.Cacher(cfg.CacherConfig())
		authService := microservice.NewAuthService(cacher, 24*3)
		publicPath := []string{
			"/swagger",
			"/login",
			"/register",

			"/employee/login",

			"/images*",
			"/productimage",

			"/healthz",
		}

		exceptShopPath := []string{
			"/shop",
			"/list-shop",
			"/select-shop",
			"/create-shop",
		}
		ms.HttpMiddleware(authService.MWFuncWithRedisMixShop(cacher, exceptShopPath, publicPath...))
		ms.RegisterLivenessProbeEndpoint("/healthz")

		//migration.StartMigrateModel(ms, cfg)
		journal.MigrationJournalTable(ms, cfg)

		authHttp := authentication.NewAuthenticationHttp(ms, cfg)
		authHttp.RouteSetup()

		shopHttp := shop.NewShopHttp(ms, cfg)
		shopHttp.RouteSetup()

		empHttp := employee.NewEmployeeHttp(ms, cfg)
		empHttp.RouteSetup()

		memberapi := member.NewMemberHttp(ms, cfg)
		memberapi.RouteSetup()

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

		shopzonehttp := shopzone.NewShopZoneHttp(ms, cfg)
		shopzonehttp.RouteSetup()

		shoptablehttp := shoptable.NewShopTableHttp(ms, cfg)
		shoptablehttp.RouteSetup()

		shopprinterhttp := shopprinter.NewShopPrinterHttp(ms, cfg)
		shopprinterhttp.RouteSetup()

		kitchenhttp := kitchen.NewKitchenHttp(ms, cfg)
		kitchenhttp.RouteSetup()

		purchase := purchase.NewPurchaseHttp(ms, cfg)
		purchase.RouteSetup()

		saleinvoiceHttp := saleinvoice.NewSaleinvoiceHttp(ms, cfg)
		saleinvoiceHttp.RouteSetup()

		chartHttp := chartofaccount.NewChartOfAccountHttp(ms, cfg)
		chartHttp.RouteSetup()

		journalhttp := journal.NewJournalHttp(ms, cfg)
		journalhttp.RouteSetup()

		journalReportHttp := journalreport.NewJournalReportHttp(ms, cfg)
		journalReportHttp.RouteSetup()
	}

	if devApiMode == 1 || devApiMode == 2 {
		inventorysearchconsumer.StartInventorySearchComsumerOnProductCreated(ms, cfg)
		inventorysearchconsumer.StartInventorySearchComsumerOnProductUpdated(ms, cfg)
		inventorysearchconsumer.StartInventorySearchComsumerOnProductDeleted(ms, cfg)

		saleinvoice.StartSaleinvoiceComsumeCreated(ms, cfg)

		journal.StartJournalComsumeCreated(ms, cfg, "99")
		journal.StartJournalComsumeBlukCreated(ms, cfg, "11")
	}

	ms.Start()
}
