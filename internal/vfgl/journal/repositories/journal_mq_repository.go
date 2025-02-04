package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/vfgl/journal/config"
	"smlaicloudplatform/internal/vfgl/journal/models"
	"smlaicloudplatform/pkg/microservice"
)

type IJournalMqRepository interface {
	Create(doc models.JournalDoc) error
	Update(doc models.JournalDoc) error
	Delete(doc models.JournalDoc) error
	CreateInBatch(docList []models.JournalDoc) error
	DeleteInBatch(docList []models.JournalDoc) error
}

type JournalMqRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.JournalDoc]
}

func NewJournalMqRepository(prod microservice.IProducer) JournalMqRepository {
	mqKey := ""

	insRepo := JournalMqRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.JournalDoc](prod, config.JournalMessageQueueConfig{}, "")
	return insRepo
}
