package repositories

import (
	"smlaicloudplatform/internal/creditorprocess/config"
	"smlaicloudplatform/internal/creditorprocess/models"
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/pkg/microservice"
)

type ICreditorProcessMessageQueueRepository interface {
	Create(doc models.CreditorProcessRequest) error
	Update(doc models.CreditorProcessRequest) error
	Delete(doc models.CreditorProcessRequest) error
	CreateInBatch(docList []models.CreditorProcessRequest) error
	UpdateInBatch(docList []models.CreditorProcessRequest) error
	DeleteInBatch(docList []models.CreditorProcessRequest) error
}

type CreditorProcessMessageQueueRepository struct {
	producer microservice.IProducer
	mqKey    string
	repositories.KafkaRepository[models.CreditorProcessRequest]
}

func NewCreditorProcessMessageQueueRepository(producer microservice.IProducer) ICreditorProcessMessageQueueRepository {
	mqKey := ""

	insRepo := CreditorProcessMessageQueueRepository{
		producer: producer,
		mqKey:    mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.CreditorProcessRequest](producer, config.CreditorProcessMessageQueueConfig{}, "")
	return insRepo
}
