package microservice

import "smlcloudplatform/internal/microservice/models"

type IContext interface {
	Log(message string)
	UserInfo() models.UserInfo
	Header(attribute string) string
	Param(name string) string
	QueryParam(name string) string
	ReadInput() string
	Response(responseCode int, responseData interface{})
	ResponseS(responseCode int, responseData string)
	ResponseError(responseCode int, errorMessage string)
	Validate(model interface{}) error

	Persister(cfg IPersisterConfig) IPersister
	Cacher(cacherConfig ICacherConfig) ICacher
	Producer(servers string) IProducer
	MQ(servers string) IMQ
}
