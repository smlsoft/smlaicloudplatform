package main

import (
	"fmt"
	"os"
	"log"	
	"net/http"
	"github.com/joho/godotenv"
	"github.com/swaggo/echo-swagger"

	_ "smlcloudplatform/cmd/productapi/docs"	
	// "smlcloudplatform/pkg/models"
	"smlcloudplatform/internal/microservice"
	
)

// @title ProductAPI
// @version 1.0
// @description This Api for Manage Product For SMLCloudPlatform

// @contact.name API Support
// @contact.url http://www.smlcloudplatform.io
// @contact.email info@smlcloudplatform.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Start Product Service")
	cfg := microservice.NewConfig()

	ms := microservice.NewMicroservice(cfg)

	serviceID := os.Getenv("SERVICE_ID")

	fmt.Printf("Service ID :: %s \n", serviceID)

	NewProductService(ms, cfg)

	ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)
	ms.Start()
	// find by shop_id

}

func NewProductService(ms *microservice.Microservice, cfg microservice.IConfig) {

	ms.GET("/", FetchProduct)
	

	ms.GET("/fetchupdate", func(ctx microservice.IServiceContext) error {
		ctx.ResponseS(200, os.Getenv("SERVICE_ID"))
		return nil
	})


	ms.POST("/", func(ctx microservice.IServiceContext) error {
		input := ctx.ReadInput()
		ctx.Log("Receive Update Arm : " + input)
		ctx.Response(http.StatusOK, map[string]interface{}{"success": true})
		return nil
	})


	ms.PUT("/", func(ctx microservice.IServiceContext) error {
		ctx.ResponseS(200, os.Getenv("SERVICE_ID"))
		return nil
	})

	fmt.Print("Start Product Service")

}

// FetchProduct godoc
// @Summary Fetch Product
// @Tags Inventory
// @Success 200 {array} models.Inventory
// @Router / [get]
func FetchProduct(ctx microservice.IServiceContext) error {
	ctx.ResponseS(200, "TEST")
	return nil
}