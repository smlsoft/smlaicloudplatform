package repositories

import (
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/journal/config"
	"smlcloudplatform/pkg/vfgl/journal/models"
)

type IJournalMqRepository interface {
	Create(doc models.JournalDoc) error
	Update(doc models.JournalDoc) error
	Delete(doc common.Identity) error
	CreateInBatch(docList []models.JournalDoc) error
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
