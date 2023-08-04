package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stockadjustment/config"
	"smlcloudplatform/pkg/transaction/stockadjustment/models"
)

type IStockAdjustmentMessageQueueRepository interface {
	Create(doc models.StockAdjustmentDoc) error
	Update(doc models.StockAdjustmentDoc) error
	Delete(doc models.StockAdjustmentDoc) error
	CreateInBatch(docList []models.StockAdjustmentDoc) error
	UpdateInBatch(docList []models.StockAdjustmentDoc) error
	DeleteInBatch(docList []models.StockAdjustmentDoc) error
}

type StockAdjustmentMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.StockAdjustmentDoc]
}

func NewStockAdjustmentMessageQueueRepository(prod microservice.IProducer) StockAdjustmentMessageQueueRepository {
	mqKey := ""

	insRepo := StockAdjustmentMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.StockAdjustmentDoc](prod, config.StockAdjustmentMessageQueueConfig{}, "")
	return insRepo
}
