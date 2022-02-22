package main

import (
	"fmt"

	_ "smlcloudplatform/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {

}

// @title           SML Cloud Platform API
// @version         1.0
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080

func main() {
	fmt.Println("Start Swagger API")

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
