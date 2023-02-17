package task

import (
	"encoding/json"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	repositoriesDocumentImage "smlcloudplatform/pkg/documentwarehouse/documentimage/repositories"
	servicesDocumentImage "smlcloudplatform/pkg/documentwarehouse/documentimage/services"
	"smlcloudplatform/pkg/task/repositories"
	"smlcloudplatform/pkg/task/services"
	"time"
)

const (
	MQ_GROUP_TASK          string = "consume-task-3"
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
	svc         services.ITaskHttpService
}

func NewTaskConsumer(ms *microservice.Microservice, cfg microservice.IConfig) *TaskConsumer {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	repo := repositories.NewTaskRepository(pst)
	repoImage := repositoriesDocumentImage.NewDocumentImageRepository(pst)
	repoImageGroup := repositoriesDocumentImage.NewDocumentImageGroupRepository(pst)
	azureblob := microservice.NewPersisterAzureBlob()

	repoImageMessagequeue := repositoriesDocumentImage.NewDocumentImageMessageQueueRepository(prod)

	svcImage := servicesDocumentImage.NewDocumentImageService(repoImage, repoImageGroup, repoImageMessagequeue, azureblob)
	svc := services.NewTaskHttpService(repo, repoImageGroup, svcImage)

	return &TaskConsumer{
		ms:          ms,
		cfg:         cfg,
		svc:         svc,
		consumerCfg: TaskMessageQueueConfig{},
	}
}

// imprement microservice consumer
func (c *TaskConsumer) RegisterConsumer() {

	mqConfig := c.cfg.MQConfig()
	timeout := time.Duration(-1)

	// create topic
	mq := microservice.NewMQ(mqConfig, c.ms.Logger)
	mq.CreateTopicR(c.consumerCfg.TopicTaskChanged(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(c.consumerCfg.TopicTaskRejected(), 5, 1, time.Hour*24*7)

	c.ms.Consume(mqConfig.URI(), c.consumerCfg.TopicTaskChanged(), c.consumerCfg.ConsumerGroup(), timeout, c.ConsumeOnProductCreated)
	c.ms.Consume(mqConfig.URI(), c.consumerCfg.TopicTaskRejected(), c.consumerCfg.ConsumerGroup(), timeout, c.ConsumeOnProductRejected)
}

func (c *TaskConsumer) ConsumeOnProductCreated(ctx microservice.IContext) error {

	msg := ctx.ReadInput()

	doc := models.DocumentImageTaskChangeMessage{}

	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// fmt.Println("consume task changed")
	// fmt.Printf("%v\n", doc)
	// fmt.Printf("task: %s \n", doc.TaskGUID)
	// fmt.Printf("count: %d \n", doc.Count)

	err = c.svc.UpdateTaskTotalImage(doc)

	if err != nil {
		return err
	}

	return nil
}

func (c *TaskConsumer) ConsumeOnProductRejected(ctx microservice.IContext) error {

	msg := ctx.ReadInput()

	doc := models.DocumentImageTaskRejectMessage{}

	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// fmt.Println("consume task rejected")
	// fmt.Printf("%v\n", doc)
	// fmt.Printf("task: %s \n", doc.TaskGUID)
	// fmt.Printf("count: %d \n", doc.Count)

	err = c.svc.UpdateTaskTotalRejectImage(doc)

	if err != nil {
		return err
	}

	return nil
}
