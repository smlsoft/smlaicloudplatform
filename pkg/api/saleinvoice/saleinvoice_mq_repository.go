package saleinvoice

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type ISaleinvoiceMQRepository interface {
	Create(doc models.SaleinvoiceData) error
	Update(doc models.SaleinvoiceData) error
	Delete(doc models.Identity) error
}

type SaleinvoiceMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewSaleinvoiceMQRepository(prod microservice.IProducer) SaleinvoiceMQRepository {
	mqKey := ""

	return SaleinvoiceMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo SaleinvoiceMQRepository) Create(doc models.SaleinvoiceData) error {
	err := repo.prod.SendMessage(MQ_TOPIC_TRANSACTION_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
func (repo SaleinvoiceMQRepository) Update(doc models.SaleinvoiceData) error {
	err := repo.prod.SendMessage(MQ_TOPIC_TRANSACTION_UPDATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo SaleinvoiceMQRepository) Delete(doc models.Identity) error {
	err := repo.prod.SendMessage(MQ_TOPIC_TRANSACTION_DELETED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
