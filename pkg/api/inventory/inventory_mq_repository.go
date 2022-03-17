package inventory

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type IInventoryMQRepository interface {
	Create(doc models.InventoryRequest) error
}

type InventoryMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewInventoryMQRepository(prod microservice.IProducer) IInventoryMQRepository {
	mqKey := ""

	return &InventoryMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo InventoryMQRepository) Create(doc models.InventoryRequest) error {
	err := repo.prod.SendMessage(MQ_TOPIC_INVENTORY_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
