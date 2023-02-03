package task

import (
	"encoding/json"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"time"
)

const (
	MQ_GROUP_TASK          string = "consume-task"
	MQ_TOPIC_TASK_CHANGED  string = "when-documentimage-task-changed"
	MQ_TOPIC_TASK_REJECTED string = "when-documentimage-task-rejected"
)

type TaskMessageQueueConfig struct{}

func (TaskMessageQueueConfig) ConsumerGroup() string {
	return MQ_GROUP_TASK
}

func (TaskMessageQueueConfig) TopicTaskChanged() string {
	return MQ_TOPIC_TASK_CHANGED
}

func (TaskMessageQueueConfig) TopicTaskRejected() string {
	return MQ_TOPIC_TASK_REJECTED
}

type TaskConsumer struct {
	ms          *microservice.Microservice
	cfg         microservice.IConfig
	consumerCfg TaskMessageQueueConfig
}

func NewTaskConsumer(ms *microservice.Microservice, cfg microservice.IConfig) *TaskConsumer {
	return &TaskConsumer{
		ms:          ms,
		cfg:         cfg,
		consumerCfg: TaskMessageQueueConfig{},
	}
}

func (c *TaskConsumer) Consume() error {
	return nil
}

// imprement microservice consumer
func (c *TaskConsumer) RegisterConsumer() {

	mqConfig := c.cfg.MQConfig()
	timeout := time.Duration(-1)

	// create topic
	mq := microservice.NewMQ(mqConfig, c.ms.Logger)
	mq.CreateTopicR(c.consumerCfg.TopicTaskChanged(), 5, 1, time.Hour*24*7)

	c.ms.Consume(mqConfig.URI(), c.consumerCfg.TopicTaskChanged(), c.consumerCfg.ConsumerGroup(), timeout, c.ConsumeOnProductCreated)
}

func (c *TaskConsumer) ConsumeOnProductCreated(ctx microservice.IContext) error {

	msg := ctx.ReadInput()

	doc := models.DocumentImageTaskChangeMessage{}

	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("consume task changed")
	fmt.Printf("%v", doc)
	// err = c.svc.Insert(product)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return err
	// }
	return nil
}
