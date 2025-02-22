package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/transaction/receivableother/config"
	"smlaicloudplatform/internal/transaction/receivableother/models"
	"smlaicloudplatform/pkg/microservice"
)

type IReceivableOtherMessageQueueRepository interface {
	Create(doc models.ReceivableOtherDoc) error
	Update(doc models.ReceivableOtherDoc) error
	Delete(doc models.ReceivableOtherDoc) error
	CreateInBatch(docList []models.ReceivableOtherDoc) error
	UpdateInBatch(docList []models.ReceivableOtherDoc) error
	DeleteInBatch(docList []models.ReceivableOtherDoc) error
}

type ReceivableOtherMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.ReceivableOtherDoc]
}

func NewReceivableOtherMessageQueueRepository(prod microservice.IProducer) IReceivableOtherMessageQueueRepository {
	mqKey := ""

	insRepo := ReceivableOtherMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.ReceivableOtherDoc](prod, config.MessageQueueConfig{}, "")
	return insRepo
}
