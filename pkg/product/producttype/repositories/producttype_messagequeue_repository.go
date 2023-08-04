package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/producttype/config"
	"smlcloudplatform/pkg/product/producttype/models"
	"smlcloudplatform/pkg/repositories"
)

type IProductTypeMessageQueueRepository interface {
	Create(doc models.ProductTypeDoc) error
	Update(doc models.ProductTypeDoc) error
	Delete(doc models.ProductTypeDoc) error
	CreateInBatch(docList []models.ProductTypeDoc) error
	UpdateInBatch(docList []models.ProductTypeDoc) error
	DeleteInBatch(docList []models.ProductTypeDoc) error
}

type ProductTypeMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.ProductTypeDoc]
}

func NewProductTypeMessageQueueRepository(prod microservice.IProducer) ProductTypeMessageQueueRepository {
	mqKey := ""

	insRepo := ProductTypeMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.ProductTypeDoc](prod, config.ProductTypeMessageQueueConfig{}, "")
	return insRepo
}
