package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/journal/config"
)

type IJournalMqRepository interface {
	Create(doc vfgl.JournalDoc) error
	Update(doc vfgl.JournalDoc) error
	Delete(doc models.Identity) error
	CreateInBatch(docList []vfgl.JournalDoc) error
}

type JournalMqRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[vfgl.JournalDoc]
}

func NewJournalMqRepository(prod microservice.IProducer) JournalMqRepository {
	mqKey := ""

	insRepo := JournalMqRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[vfgl.JournalDoc](prod, config.JournalMessageQueueConfig{}, "")
	return insRepo
}
