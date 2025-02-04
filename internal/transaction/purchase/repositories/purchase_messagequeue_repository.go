package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/transaction/purchase/config"
	"smlaicloudplatform/internal/transaction/purchase/models"
	"smlaicloudplatform/pkg/microservice"
)

type IPurchaseMessageQueueRepository interface {
	Create(doc models.PurchaseDoc) error
	Update(doc models.PurchaseDoc) error
	Delete(doc models.PurchaseDoc) error
	CreateInBatch(docList []models.PurchaseDoc) error
	UpdateInBatch(docList []models.PurchaseDoc) error
	DeleteInBatch(docList []models.PurchaseDoc) error
}

type PurchaseMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.PurchaseDoc]
}

func NewPurchaseMessageQueueRepository(prod microservice.IProducer) PurchaseMessageQueueRepository {
	mqKey := ""

	insRepo := PurchaseMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.PurchaseDoc](prod, config.PurchaseMessageQueueConfig{}, "")
	return insRepo
}
