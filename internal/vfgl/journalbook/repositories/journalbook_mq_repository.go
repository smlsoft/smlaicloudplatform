package repositories

import (
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/vfgl/journalbook/config"
	"smlaicloudplatform/internal/vfgl/journalbook/models"
	"smlaicloudplatform/pkg/microservice"
)

type IJournalBookMqRepository interface {
	Create(doc models.JournalBookDoc) error
	Update(doc models.JournalBookDoc) error
	Delete(doc common.Identity) error
	CreateInBatch(docList []models.JournalBookDoc) error
}

type JournalBookMqRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.JournalBookDoc]
}

func NewJournalBookMqRepository(prod microservice.IProducer) JournalBookMqRepository {
	mqKey := ""

	insRepo := JournalBookMqRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.JournalBookDoc](prod, config.JournalBookMessageQueueConfig{}, "")
	return insRepo
}
