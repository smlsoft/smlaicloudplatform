package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/vfgl/chartofaccount/config"
	"smlaicloudplatform/internal/vfgl/chartofaccount/models"
	"smlaicloudplatform/pkg/microservice"
)

type IChartOfAccountMQRepository interface {
	Create(doc models.ChartOfAccountDoc) error
	Update(doc models.ChartOfAccountDoc) error
	Delete(doc models.ChartOfAccountDoc) error
	CreateInBatch(docList []models.ChartOfAccountDoc) error
}

type ChartOfAccountMQRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.ChartOfAccountDoc]
}

func NewChartOfAccountMQRepository(prod microservice.IProducer) ChartOfAccountMQRepository {
	mqKey := ""

	repo := ChartOfAccountMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	repo.KafkaRepository = repositories.NewKafkaRepository[models.ChartOfAccountDoc](prod, config.ChartOfAccountMessageQueueConfig{}, "")
	return repo
}
