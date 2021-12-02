package microservice

import "smlcloudplatform/internal/microservice/models"

type IServiceContext interface {
	Log(message string)
	UserInfo() models.UserInfo
	Param(name string) string
	QueryParam(name string) string
	ReadInput() string
	Response(responseCode int, responseData interface{})
	ResponseS(responseCode int, responseData string)
	ResponseError(responseCode int, errorMessage string)
	Validate(model interface{}) error

	Persister(cfg IPersisterConfig) IPersister
	Producer(servers string) IProducer
	MQ(servers string) IMQ
}
