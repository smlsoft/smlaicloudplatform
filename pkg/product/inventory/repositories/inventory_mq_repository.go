package repositories

import (
	"context"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/inventory/config"
	"smlcloudplatform/pkg/product/inventory/models"
)

type IInventoryMQRepository interface {
	Create(ctx context.Context, doc models.InventoryData) error
	Update(doc models.InventoryData) error
	Delete(doc common.Identity) error
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

func (repo InventoryMQRepository) Create(ctx context.Context, doc models.InventoryData) error {
	err := repo.prod.SendMessage(config.MQ_TOPIC_INVENTORY_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryMQRepository) Update(doc models.InventoryData) error {
	err := repo.prod.SendMessage(config.MQ_TOPIC_INVENTORY_UPDATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryMQRepository) Delete(doc common.Identity) error {
	err := repo.prod.SendMessage(config.MQ_TOPIC_INVENTORY_DELETED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
