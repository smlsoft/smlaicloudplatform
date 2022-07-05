package stockinout

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/transaction/stockinout/models"
)

type IStockInOutMQRepository interface {
	Create(doc models.StockInOutData) error
}

type StockInOutMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewStockInOutMQRepository(prod microservice.IProducer) StockInOutMQRepository {
	mqKey := ""

	return StockInOutMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo StockInOutMQRepository) Create(doc models.StockInOutData) error {
	err := repo.prod.SendMessage(MQ_TOPIC_STOCK_IN_OUT_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
