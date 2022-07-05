package purchase

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/transaction/purchase/models"
)

type IPurchaseMQRepository interface {
	Create(doc models.PurchaseData) error
}

type PurchaseMQRepository struct {
	prod  microservice.IProducer
	mqKey string
}

func NewPurchaseMQRepository(prod microservice.IProducer) PurchaseMQRepository {
	mqKey := ""

	return PurchaseMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
}

func (repo PurchaseMQRepository) Create(doc models.PurchaseData) error {
	err := repo.prod.SendMessage(MQ_TOPIC_PURCHASE_CREATED, repo.mqKey, doc)

	if err != nil {
		return err
	}

	return nil
}
