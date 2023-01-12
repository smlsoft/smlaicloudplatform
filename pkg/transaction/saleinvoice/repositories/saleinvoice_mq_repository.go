package repositories

import (
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/config"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"
)

type ISaleinvoiceMQRepository interface {
	Create(doc models.SaleinvoiceData) error
	Update(doc models.SaleinvoiceData) error
	Delete(doc common.Identity) error
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
	err := repo.prod.SendMessage(config.MQ_TOPIC_SALEINVOICE_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
func (repo SaleinvoiceMQRepository) Update(doc models.SaleinvoiceData) error {
	err := repo.prod.SendMessage(config.MQ_TOPIC_SALEINVOICE_UPDATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo SaleinvoiceMQRepository) Delete(doc common.Identity) error {
	err := repo.prod.SendMessage(config.MQ_TOPIC_SALEINVOICE_DELETED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
