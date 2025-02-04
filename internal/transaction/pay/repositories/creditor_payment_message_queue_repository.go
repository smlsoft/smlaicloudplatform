package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/transaction/pay/config"
	"smlaicloudplatform/internal/transaction/pay/models"
	"smlaicloudplatform/pkg/microservice"
)

type ICreditorPaymentMessageQueueRepository interface {
	Create(doc models.PayDoc) error
	Update(doc models.PayDoc) error
	Delete(doc models.PayDoc) error
	CreateInBatch(docList []models.PayDoc) error
	UpdateInBatch(docList []models.PayDoc) error
	DeleteInBatch(docList []models.PayDoc) error
}

type PaidMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.PayDoc]
}

func NewPaidMessageQueueRepository(prod microservice.IProducer) ICreditorPaymentMessageQueueRepository {
	mqKey := ""

	insRepo := PaidMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.PayDoc](prod, config.CreditorPaymentMessageQueueConfig{}, "")
	return insRepo
}
