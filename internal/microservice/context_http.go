package microservice

import (
	"fmt"
	"io/ioutil"
	"smlcloudplatform/internal/microservice/models"

	"github.com/labstack/echo/v4"
)

// HTTPContext implement IContext it is context for HTTP
type HTTPContext struct {
	ms *Microservice
	c  echo.Context
}

// NewHTTPContext is the constructor function for HTTPContext
func NewHTTPContext(ms *Microservice, c echo.Context) *HTTPContext {

	return &HTTPContext{
		ms: ms,
		c:  c,
	}
}

// Log will log a message
func (ctx *HTTPContext) Log(message string) {
	fmt.Println("HTTP: ", message)
}

// Param return parameter by name
func (ctx *HTTPContext) Param(name string) string {
	return ctx.c.Param(name)
}

// QueryParam return query param
func (ctx *HTTPContext) QueryParam(name string) string {
	return ctx.c.QueryParam(name)
}

// ReadInput read the request body and return it as string
func (ctx *HTTPContext) ReadInput() string {
	body, err := ioutil.ReadAll(ctx.c.Request().Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// Header return header value by key
func (ctx *HTTPContext) Header(attribute string) string {
	return ctx.c.Request().Header.Get(attribute)
}

func (ctx *HTTPContext) UserInfo() models.UserInfo {
	userInfo := ctx.c.Get("UserInfo")
	if userInfo == nil {
		return models.UserInfo{}
	}
	return userInfo.(models.UserInfo)

}

func (ctx *HTTPContext) Validate(model interface{}) error {
	if err := ctx.c.Validate(model); err != nil {
		return err
	}
	return nil
}

func (ctx *HTTPContext) Response(responseCode int, responseData interface{}) {
	ctx.c.JSON(responseCode, responseData)
}

func (ctx *HTTPContext) ResponseS(responseCode int, responseData string) {
	ctx.c.String(responseCode, responseData)
}

func (ctx *HTTPContext) ResponseError(responseCode int, errorMessage string) {
	ctx.c.JSON(responseCode, map[string]interface{}{
		"success": false,
		"message": errorMessage,
	})
}

// Persister return perister
func (ctx *HTTPContext) Persister(cfg IPersisterConfig) IPersister {
	return ctx.ms.Persister(cfg)
}

// Cacher return cacher
func (ctx *HTTPContext) Cacher(cacheConfig ICacherConfig) ICacher {
	return ctx.ms.Cacher(cacheConfig)
}

// Producer return producer
func (ctx *HTTPContext) Producer(mqConfig IMQConfig) IProducer {
	return ctx.ms.Producer(mqConfig)
}

// MQ return MQ
func (ctx *HTTPContext) MQ(mqConfig IMQConfig) IMQ {
	return NewMQ(mqConfig, ctx.ms.Logger)
}
