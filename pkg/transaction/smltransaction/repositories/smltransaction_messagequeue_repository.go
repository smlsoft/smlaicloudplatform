package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/transaction/smltransaction/config"
	"smlcloudplatform/pkg/transaction/smltransaction/models"
)

type ISMLTransactionMessageQueueRepository interface {
	Save(doc models.SMLTransactionRequest) error
	Delete(doc models.SMLTransactionKeyRequest) error
}

type SMLTransactionMessageQueueRepository struct {
	prod   microservice.IProducer
	mqKey  string
	config config.SMLTransactionMessageQueueConfig
}

func NewSMLTransactionMessageQueueRepository(prod microservice.IProducer) SMLTransactionMessageQueueRepository {
	mqKey := ""

	insRepo := SMLTransactionMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.config = config.SMLTransactionMessageQueueConfig{}
	return insRepo
}

func (repo SMLTransactionMessageQueueRepository) Save(doc models.SMLTransactionRequest) error {
	err := repo.prod.SendMessage(repo.config.TopicSaved(), repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo SMLTransactionMessageQueueRepository) Delete(doc models.SMLTransactionKeyRequest) error {
	err := repo.prod.SendMessage(repo.config.TopicDeleted(), repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
