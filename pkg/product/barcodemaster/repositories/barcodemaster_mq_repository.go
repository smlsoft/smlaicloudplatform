package repositories

import (
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/barcodemaster/config"
	"smlcloudplatform/pkg/product/barcodemaster/models"
)

type IBarcodeMasterMQRepository interface {
	Create(doc models.BarcodeMasterData) error
	Update(doc models.BarcodeMasterData) error
	Delete(doc common.Identity) error
}
type BarcodeMasterMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewBarcodeMasterMQRepository(prod microservice.IProducer) IBarcodeMasterMQRepository {
	mqKey := ""

	return BarcodeMasterMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo BarcodeMasterMQRepository) Create(doc models.BarcodeMasterData) error {
	err := repo.prod.SendMessage(config.MQ_TOPIC_BARCODEMASTER_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo BarcodeMasterMQRepository) Update(doc models.BarcodeMasterData) error {
	err := repo.prod.SendMessage(config.MQ_TOPIC_BARCODEMASTER_UPDATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo BarcodeMasterMQRepository) Delete(doc common.Identity) error {
	err := repo.prod.SendMessage(config.MQ_TOPIC_BARCODEMASTER_DELETED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
