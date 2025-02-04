package repositories

import (
	"smlaicloudplatform/internal/pos/shift/config"
	"smlaicloudplatform/internal/pos/shift/models"
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/pkg/microservice"
)

type IShiftMessageQueueRepository interface {
	Create(doc models.ShiftDoc) error
	Update(doc models.ShiftDoc) error
	Delete(doc models.ShiftDoc) error
	CreateInBatch(docList []models.ShiftDoc) error
	UpdateInBatch(docList []models.ShiftDoc) error
	DeleteInBatch(docList []models.ShiftDoc) error
}

type ShiftMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.ShiftDoc]
}

func NewShiftMessageQueueRepository(prod microservice.IProducer) ShiftMessageQueueRepository {
	mqKey := ""

	insRepo := ShiftMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.ShiftDoc](prod, config.ShiftMessageQueueConfig{}, "")
	return insRepo
}
