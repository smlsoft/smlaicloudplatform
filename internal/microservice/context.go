package microservice

import (
	"mime/multipart"
	"net/http"
	"smlcloudplatform/internal/microservice/models"

	"github.com/labstack/echo/v4"
)

type IContext interface {
	Log(message string)
	UserInfo() models.UserInfo
	Header(attribute string) string
	Param(name string) string
	QueryParam(name string) string
	ReadInput() string
	Response(responseCode int, responseData interface{})
	ResponseError(responseCode int, errorMessage string)
	Validate(model interface{}) error
	FormFile(field string) (*multipart.FileHeader, error)
	ResponseWriter() http.ResponseWriter
	Request() *http.Request

	Persister(cfg IPersisterConfig) IPersister
	Cacher(cacherConfig ICacherConfig) ICacher
	Producer(servers IMQConfig) IProducer
	MQ(servers IMQConfig) IMQ

	EchoContext() echo.Context
}
