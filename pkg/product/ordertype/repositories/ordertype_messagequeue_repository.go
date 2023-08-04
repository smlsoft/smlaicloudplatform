package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/ordertype/config"
	"smlcloudplatform/pkg/product/ordertype/models"
	"smlcloudplatform/pkg/repositories"
)

type IOrderTypeMessageQueueRepository interface {
	Create(doc models.OrderTypeDoc) error
	Update(doc models.OrderTypeDoc) error
	Delete(doc models.OrderTypeDoc) error
	CreateInBatch(docList []models.OrderTypeDoc) error
	UpdateInBatch(docList []models.OrderTypeDoc) error
	DeleteInBatch(docList []models.OrderTypeDoc) error
}

type OrderTypeMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.OrderTypeDoc]
}

func NewOrderTypeMessageQueueRepository(prod microservice.IProducer) OrderTypeMessageQueueRepository {
	mqKey := ""

	insRepo := OrderTypeMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.OrderTypeDoc](prod, config.OrderTypeMessageQueueConfig{}, "")
	return insRepo
}
