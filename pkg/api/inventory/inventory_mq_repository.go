package inventory

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type IInventoryMQRepository interface {
	Create(doc models.InventoryData) error
	Update(doc models.InventoryData) error
	Delete(doc models.Identity) error
}
type InventoryMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewInventoryMQRepository(prod microservice.IProducer) IInventoryMQRepository {
	mqKey := ""

	return InventoryMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo InventoryMQRepository) Create(doc models.InventoryData) error {
	err := repo.prod.SendMessage(MQ_TOPIC_INVENTORY_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryMQRepository) Update(doc models.InventoryData) error {
	err := repo.prod.SendMessage(MQ_TOPIC_INVENTORY_UPDATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryMQRepository) Delete(doc models.Identity) error {
	err := repo.prod.SendMessage(MQ_TOPIC_INVENTORY_DELETED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
