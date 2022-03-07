package purchase

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type IPurchaseMQRepository interface {
	Create(doc models.PurchaseRequest) error
}

type PurchaseMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewPurchaseMQRepository(prod microservice.IProducer) IPurchaseMQRepository {
	mqKey := ""

	return &PurchaseMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo PurchaseMQRepository) Create(doc models.PurchaseRequest) error {
	err := repo.prod.SendMessage(MQ_TOPIC_PURCHASE_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
