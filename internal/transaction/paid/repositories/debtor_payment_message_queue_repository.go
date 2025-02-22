package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/transaction/paid/config"
	"smlaicloudplatform/internal/transaction/paid/models"
	"smlaicloudplatform/pkg/microservice"
)

type IDebtorPaymentMessageQueueRepository interface {
	Create(doc models.PaidDoc) error
	Update(doc models.PaidDoc) error
	Delete(doc models.PaidDoc) error
	CreateInBatch(docList []models.PaidDoc) error
	UpdateInBatch(docList []models.PaidDoc) error
	DeleteInBatch(docList []models.PaidDoc) error
}

type PaidMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.PaidDoc]
}

func NewPaidMessageQueueRepository(prod microservice.IProducer) IDebtorPaymentMessageQueueRepository {
	mqKey := ""

	insRepo := PaidMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.PaidDoc](prod, config.DebtorPaymentMessageQueueConfig{}, "")
	return insRepo
}
