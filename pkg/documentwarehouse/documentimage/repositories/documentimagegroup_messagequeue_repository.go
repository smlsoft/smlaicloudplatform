package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/config"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
)

type DocumentImageMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	topic config.DocumentImageMessageQueueConfig
}

func NewDocumentImageMessageQueueRepository(prod microservice.IProducer) DocumentImageMessageQueueRepository {
	mqKey := ""

	insRepo := DocumentImageMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}

	return insRepo
}

func (repo DocumentImageMessageQueueRepository) TaskChange(doc models.DocumentImageTaskChangeMessage) error {

	err := repo.prod.SendMessage(repo.topic.TopicTaskChanged(), repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo DocumentImageMessageQueueRepository) TaskReject(doc models.DocumentImageTaskRejectMessage) error {
	err := repo.prod.SendMessage(repo.topic.TopicTaskRejected(), repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
