package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/purchase/config"
	"smlcloudplatform/pkg/transaction/purchase/models"
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
