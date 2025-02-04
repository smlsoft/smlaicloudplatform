package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/warehouse/config"
	"smlaicloudplatform/internal/warehouse/models"
	"smlaicloudplatform/pkg/microservice"
)

type IWarehouseMessageQueueRepository interface {
	Create(doc models.WarehouseDoc) error
	Update(doc models.WarehouseDoc) error
	Delete(doc models.WarehouseDoc) error
	CreateInBatch(docList []models.WarehouseDoc) error
	UpdateInBatch(docList []models.WarehouseDoc) error
	DeleteInBatch(docList []models.WarehouseDoc) error
}

type WarehouseMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.WarehouseDoc]
}

func NewWarehouseMessageQueueRepository(prod microservice.IProducer) WarehouseMessageQueueRepository {
	mqKey := ""

	insRepo := WarehouseMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.WarehouseDoc](prod, config.WarehouseMessageQueueConfig{}, "")
	return insRepo
}
