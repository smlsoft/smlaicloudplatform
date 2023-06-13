package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/purchasereturn/config"
	"smlcloudplatform/pkg/transaction/purchasereturn/models"
)

type IPurchaseReturnMessageQueueRepository interface {
	Create(doc models.PurchaseReturnDoc) error
	Update(doc models.PurchaseReturnDoc) error
	Delete(doc models.PurchaseReturnDoc) error
	CreateInBatch(docList []models.PurchaseReturnDoc) error
	UpdateInBatch(docList []models.PurchaseReturnDoc) error
	DeleteInBatch(docList []models.PurchaseReturnDoc) error
}

type PurchaseReturnMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.PurchaseReturnDoc]
}

func NewPurchaseReturnMessageQueueRepository(prod microservice.IProducer) PurchaseReturnMessageQueueRepository {
	mqKey := ""

	insRepo := PurchaseReturnMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.PurchaseReturnDoc](prod, config.PurchaseReturnMessageQueueConfig{}, "")
	return insRepo
}
