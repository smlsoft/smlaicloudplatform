package repositories

import (
	"context"
	common "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/internal/vfgl/journal/config"
	"smlcloudplatform/internal/vfgl/journal/models"
	"smlcloudplatform/pkg/microservice"
)

type IJournalMqRepository interface {
	Create(ctx context.Context, doc models.JournalDoc) error
	Update(doc models.JournalDoc) error
	Delete(doc common.Identity) error
	CreateInBatch(ctx context.Context, docList []models.JournalDoc) error
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
