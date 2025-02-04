package repositories

import (
	"smlaicloudplatform/internal/repositories"
	config "smlaicloudplatform/internal/stockprocess/config"
	"smlaicloudplatform/pkg/microservice"

	models "smlaicloudplatform/internal/stockprocess/models"
)

type IStockProcessMessageQueueRepository interface {
	Create(doc models.StockProcessRequest) error
	Update(doc models.StockProcessRequest) error
	Delete(doc models.StockProcessRequest) error
	CreateInBatch(docList []models.StockProcessRequest) error
	UpdateInBatch(docList []models.StockProcessRequest) error
	DeleteInBatch(docList []models.StockProcessRequest) error
}

type StockProcessMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.StockProcessRequest]
}

func NewStockProcessMessageQueueRepository(prod microservice.IProducer) StockProcessMessageQueueRepository {
	mqKey := ""

	insRepo := StockProcessMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.StockProcessRequest](prod, config.StockProcessMessageQueueConfig{}, "")
	return insRepo
}
