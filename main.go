package main

import (
	"fmt"
	"os"

	"smlcloudplatform/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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
	fmt.Println("Start Swagger API")

	host := os.Getenv("HOST_API")
	if host != "" {
		fmt.Printf("Host: %v\n", host)
		docs.SwaggerInfo.Host = host
	}

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "1323"
	}
	e.Logger.Fatal(e.Start(":" + serverPort))
}
