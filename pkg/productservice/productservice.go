package productservice

import (
	"fmt"
	"os"
	"net/http"
	"smlcloudplatform/internal/microservice"
)

type ProductService struct {
}

func NewProductService(ms *microservice.Microservice, cfg microservice.IConfig) {

	ms.GET("/", func(ctx microservice.IServiceContext) error {
		ctx.ResponseS(200, os.Getenv("SERVICE_ID"))
		return nil
	})
	

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
