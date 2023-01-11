package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/product/config"
	"smlcloudplatform/pkg/product/product/models"
	"smlcloudplatform/pkg/repositories"
)

type IProductMessageQueueRepository interface {
	Create(doc models.ProductDoc) error
	Update(doc models.ProductDoc) error
	Delete(doc models.ProductDoc) error
	CreateInBatch(docList []models.ProductDoc) error
	UpdateInBatch(docList []models.ProductDoc) error
	DeleteInBatch(docList []models.ProductDoc) error
}

type ProductMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.ProductDoc]
}

func NewProductMessageQueueRepository(prod microservice.IProducer) ProductMessageQueueRepository {
	mqKey := ""

	insRepo := ProductMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.ProductDoc](prod, config.ProductMessageQueueConfig{}, "")
	return insRepo
}
