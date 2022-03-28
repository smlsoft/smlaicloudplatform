package transaction

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type ITransactionMQRepository interface {
	Create(doc models.TransactionData) error
	Update(doc models.TransactionData) error
	Delete(doc models.Identity) error
}

type TransactionMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewTransactionMQRepository(prod microservice.IProducer) TransactionMQRepository {
	mqKey := ""

	return TransactionMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo TransactionMQRepository) Create(doc models.TransactionData) error {
	err := repo.prod.SendMessage(MQ_TOPIC_TRANSACTION_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
func (repo TransactionMQRepository) Update(doc models.TransactionData) error {
	err := repo.prod.SendMessage(MQ_TOPIC_TRANSACTION_UPDATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo TransactionMQRepository) Delete(doc models.Identity) error {
	err := repo.prod.SendMessage(MQ_TOPIC_TRANSACTION_DELETED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
