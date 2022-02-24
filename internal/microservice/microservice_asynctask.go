package microservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// startAsyncTaskConsumer read async task message from message queue and execute with handler
func (ms *Microservice) startAsyncTaskConsumer(path string, cacheConfig ICacherConfig, mqServers string, h ServiceHandleFunc) {
	topic := escapeName(path)
	mq := NewMQ(mqServers, ms)
	err := mq.CreateTopicR(topic, 5, 1, time.Hour*24*30) // retain message for 30 days
	if err != nil {
		ms.Log("ATASK", err.Error())
	}

	ms.Consume(mqServers, topic, "atask", -1, func(ctx IContext) error {
		message := map[string]interface{}{}
		err := json.Unmarshal([]byte(ctx.ReadInput()), &message)
		if err != nil {
			return err
		}
		ref, _ := message["ref"].(string)
		input, _ := message["input"].(string)
		return h(NewAsyncTaskContext(ms, cacheConfig, ref, input))
	})
}

// handleAsyncTaskRequest accept async task request and send it to message queue
func (ms *Microservice) handleAsyncTaskRequest(path string, cacheConfig ICacherConfig, mqServers string, ctx IContext) error {
	topic := escapeName(path)

	// 1. Read Input
	input := ctx.ReadInput()

	// 2. Generate REF
	ref := fmt.Sprintf("atask-%s", randString())

	// 3. Set Status in Cache
	cacher := ctx.Cacher(cacheConfig)
	status := map[string]interface{}{
		"status": "processing",
	}
	expire := time.Minute * 30
	cacher.Set(ref, status, expire)

	// 4. Send Message to MQ
	prod := ctx.Producer(mqServers)
	message := map[string]interface{}{
		"ref":   ref,
		"input": input,
	}
	prod.SendMessage(topic, "", message)

	// 5. Response REF
	res := map[string]string{
		"ref": ref,
	}
	ctx.Response(http.StatusOK, res)
	return nil
}

func (ms *Microservice) handleAsyncTaskResponse(path string, cacheConfig ICacherConfig, ctx IContext) error {
	// 1. ReadInput (REF from query string)
	ref := ctx.QueryParam("ref")

	// 2. Read Status from Cache
	cacher := ctx.Cacher(cacheConfig)
	statusJS, err := cacher.Get(ref)
	if err != nil {
		return err
	}

	// 3. Return Status
	status := map[string]interface{}{}
	err = json.Unmarshal([]byte(statusJS), &status)
	if err != nil {
		return err
	}
	ctx.Response(http.StatusOK, status)
	return nil
}

// AsyncPOST register async task service for HTTP POST
func (ms *Microservice) AsyncPOST(path string, cacheConfig ICacherConfig, mqServers string, h ServiceHandleFunc) {
	ms.startAsyncTaskConsumer(path, cacheConfig, mqServers, h)
	ms.GET(path, func(ctx IContext) error {
		return ms.handleAsyncTaskResponse(path, cacheConfig, ctx)
	})
	ms.POST(path, func(ctx IContext) error {
		return ms.handleAsyncTaskRequest(path, cacheConfig, mqServers, ctx)
	})
}

// AsyncPUT register async task service for HTTP PUT
func (ms *Microservice) AsyncPUT(path string, cacheConfig ICacherConfig, mqServers string, h ServiceHandleFunc) {
	ms.startAsyncTaskConsumer(path, cacheConfig, mqServers, h)
	ms.GET(path, func(ctx IContext) error {
		return ms.handleAsyncTaskResponse(path, cacheConfig, ctx)
	})
	ms.PUT(path, func(ctx IContext) error {
		return ms.handleAsyncTaskRequest(path, cacheConfig, mqServers, ctx)
	})
}
