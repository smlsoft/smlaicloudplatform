package transaction

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type ITransactionMQRepository interface {
	Create(doc models.TransactionRequest) error
}

type TransactionMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewTransactionMQRepository(prod microservice.IProducer) ITransactionMQRepository {
	mqKey := ""

	return &TransactionMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo TransactionMQRepository) Create(doc models.TransactionRequest) error {
	err := repo.prod.SendMessage(MQ_TOPIC_TRANSACTION_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
