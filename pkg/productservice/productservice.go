package productservice

import (
	"fmt"
	"os"
	"smlcloudplatform/internal/microservice"
)

type ProductService struct {
}

func NewProductService(ms *microservice.Microservice, cfg microservice.IConfig) {

	ms.GET("/", func(ctx microservice.IServiceContext) error {
		ctx.ResponseS(200, os.Getenv("SERVICE_ID"))
		return nil
	})
	a := "s"
	fmt.Print(a)
}
