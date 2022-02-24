package microservice

import (
	"fmt"
	"runtime"
	"smlcloudplatform/internal/microservice/models"
	"strings"
	"time"
)

// ConsumerContext implement IContext it is context for Consumer
type ConsumerContext struct {
	ms      *Microservice
	message string
}

// NewConsumerContext is the constructor function for ConsumerContext
func NewConsumerContext(ms *Microservice, message string) *ConsumerContext {
	return &ConsumerContext{
		ms:      ms,
		message: message,
	}
}

// Log will log a message
func (ctx *ConsumerContext) Log(message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println("Consumer:", fns[len(fns)-1], line, message)
}

// Param return parameter by name (empty in case of Consumer)
func (ctx *ConsumerContext) Param(name string) string {
	return ""
}

// QueryParam return empty in consumer
func (ctx *ConsumerContext) QueryParam(name string) string {
	return ""
}

// ReadInput return message
func (ctx *ConsumerContext) ReadInput() string {
	return ctx.message
}

// ReadInputs return nil in case Consumer
func (ctx *ConsumerContext) ReadInputs() []string {
	return nil
}

func (ctx *ConsumerContext) Validate(model interface{}) error {
	return nil
}

// Response return response to client
func (ctx *ConsumerContext) Response(responseCode int, responseData interface{}) {
	return
}

func (ctx *ConsumerContext) ResponseS(responseCode int, responseData string) {
	return
}

func (ctx *ConsumerContext) ResponseError(responseCode int, errorMessage string) {
	return
}

// Header return header value by key
func (ctx *ConsumerContext) Header(attribute string) string {
	return ""
}

func (ctx *ConsumerContext) UserInfo() models.UserInfo {
	return models.UserInfo{}

}

// Persister return perister
func (ctx *ConsumerContext) Persister(cfg IPersisterConfig) IPersister {
	return ctx.ms.Persister(cfg)
}

// Now return now
func (ctx *ConsumerContext) Now() time.Time {
	return time.Now()
}

// Cacher return cacher
func (ctx *ConsumerContext) Cacher(cacheConfig ICacherConfig) ICacher {
	return ctx.ms.Cacher(cacheConfig)
}

// Producer return producer
func (ctx *ConsumerContext) Producer(servers string) IProducer {
	return ctx.ms.getProducer(servers)
}

// MQ return MQ
func (ctx *ConsumerContext) MQ(servers string) IMQ {
	return NewMQ(servers, ctx.ms)
}
