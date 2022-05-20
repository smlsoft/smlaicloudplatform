package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/vfgl"
)

type IKafkaRepo interface {
	vfgl.JournalDoc
}

type KafkaConfig interface {
	TopicCreated() string
	TopicUpdated() string
	TopicDeleted() string
	TopicBulkCreated() string
}

type KafkaRepository[T any] struct {
	topic KafkaConfig
	prod  microservice.IProducer
	mqKey string
}

func NewKafkaRepository[T IKafkaRepo](prod microservice.IProducer, topicConfig KafkaConfig, mqKey string) KafkaRepository[T] {
	return KafkaRepository[T]{
		prod:  prod,
		topic: topicConfig,
		mqKey: mqKey,
	}
}

func (repo KafkaRepository[T]) Create(doc T) error {
	err := repo.prod.SendMessage(repo.topic.TopicCreated(), repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo KafkaRepository[T]) Update(doc T) error {
	err := repo.prod.SendMessage(repo.topic.TopicUpdated(), repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo KafkaRepository[T]) Delete(doc models.Identity) error {
	err := repo.prod.SendMessage(repo.topic.TopicDeleted(), repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo KafkaRepository[T]) CreateInBatch(docList []T) error {
	err := repo.prod.SendMessage(repo.topic.TopicBulkCreated(), repo.mqKey, docList)

	if err != nil {
		return err
	}

	return nil
}
