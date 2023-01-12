package repositories

import (
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/chartofaccount/config"
	"smlcloudplatform/pkg/vfgl/chartofaccount/models"
)

type IChartOfAccountMQRepository interface {
	Create(doc models.ChartOfAccountDoc) error
	Update(doc models.ChartOfAccountDoc) error
	Delete(doc common.Identity) error
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
