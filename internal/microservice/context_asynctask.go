package microservice

import (
	"fmt"
	"runtime"
	"smlcloudplatform/internal/microservice/models"
	"strings"
	"time"
)

// AsyncTaskContext implement IContext it is context for Consumer
type AsyncTaskContext struct {
	ms          *Microservice
	cacheConfig ICacherConfig
	ref         string
	input       string
}

// NewAsyncTaskContext is the constructor function for AsyncTaskContext
func NewAsyncTaskContext(ms *Microservice, cacheConfig ICacherConfig, ref string, input string) *AsyncTaskContext {
	return &AsyncTaskContext{
		ms:          ms,
		cacheConfig: cacheConfig,
		ref:         ref,
		input:       input,
	}
}

// Log will log a message
func (ctx *AsyncTaskContext) Log(message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println("ATASK:", fns[len(fns)-1], line, message)
}

// Param return parameter by name (empty in AsyncTask)
func (ctx *AsyncTaskContext) Param(name string) string {
	return ""
}

// QueryParam return empty in async task
func (ctx *AsyncTaskContext) QueryParam(name string) string {
	return ""
}

// ReadInput return message (return empty in AsyncTask)
func (ctx *AsyncTaskContext) ReadInput() string {
	return ctx.input
}

// ReadInputs return messages in batch (return nil in AsyncTask)
func (ctx *AsyncTaskContext) ReadInputs() []string {
	return nil
}

func (ctx *AsyncTaskContext) Validate(model interface{}) error {
	return nil
}

// Response return response to client
func (ctx *AsyncTaskContext) Response(responseCode int, responseData interface{}) {
	return
}

func (ctx *AsyncTaskContext) ResponseS(responseCode int, responseData string) {
	return
}

func (ctx *AsyncTaskContext) ResponseError(responseCode int, errorMessage string) {
	return
}

// Header return header value by key
func (ctx *AsyncTaskContext) Header(attribute string) string {
	return ""
}

func (ctx *AsyncTaskContext) UserInfo() models.UserInfo {
	return models.UserInfo{}

}

// Persister return perister
func (ctx *AsyncTaskContext) Persister(cfg IPersisterConfig) IPersister {
	return ctx.ms.Persister(cfg)
}

// Now return now
func (ctx *AsyncTaskContext) Now() time.Time {
	return time.Now()
}

// Cacher return cacher
func (ctx *AsyncTaskContext) Cacher(cacherConfig ICacherConfig) ICacher {
	return ctx.ms.Cacher(cacherConfig)
}

// Producer return producer
func (ctx *AsyncTaskContext) Producer(servers string) IProducer {
	return ctx.ms.getProducer(servers)
}

// MQ return MQ
func (ctx *AsyncTaskContext) MQ(servers string) IMQ {
	return NewMQ(servers, ctx.ms)
}
