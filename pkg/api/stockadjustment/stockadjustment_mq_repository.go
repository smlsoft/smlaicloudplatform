package stockadjustment

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type IStockAdjustmentMQRepository interface {
	Create(doc models.StockAdjustmentRequest) error
}

type StockAdjustmentMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewStockAdjustmentMQRepository(prod microservice.IProducer) StockAdjustmentMQRepository {
	mqKey := ""

	return StockAdjustmentMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo StockAdjustmentMQRepository) Create(doc models.StockAdjustmentRequest) error {
	err := repo.prod.SendMessage(MQ_TOPIC_STOCK_ADJUSTMENT_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
