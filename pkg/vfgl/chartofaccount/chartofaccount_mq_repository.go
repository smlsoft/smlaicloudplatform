package chartofaccount

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/repositories"
)

type IChartOfAccountMQRepository interface {
	Create(doc vfgl.ChartOfAccountDoc) error
	Update(doc vfgl.ChartOfAccountDoc) error
	Delete(doc models.Identity) error
	CreateInBatch(docList []vfgl.ChartOfAccountDoc) error
}

type ChartOfAccountMQRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[vfgl.ChartOfAccountDoc]
}

func NewChartOfAccountMQRepository(prod microservice.IProducer) ChartOfAccountMQRepository {
	mqKey := ""

	repo := ChartOfAccountMQRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	repo.KafkaRepository = repositories.NewKafkaRepository[vfgl.ChartOfAccountDoc](prod, ChartOfAccountMessageQueueConfig{}, "")
	return repo
}
