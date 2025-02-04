package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/transaction/stocktransfer/config"
	"smlaicloudplatform/internal/transaction/stocktransfer/models"
	"smlaicloudplatform/pkg/microservice"
)

type IStockTransferMessageQueueRepository interface {
	Create(doc models.StockTransferDoc) error
	Update(doc models.StockTransferDoc) error
	Delete(doc models.StockTransferDoc) error
	CreateInBatch(docList []models.StockTransferDoc) error
	UpdateInBatch(docList []models.StockTransferDoc) error
	DeleteInBatch(docList []models.StockTransferDoc) error
}

type StockTransferMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.StockTransferDoc]
}

func NewStockTransferMessageQueueRepository(prod microservice.IProducer) StockTransferMessageQueueRepository {
	mqKey := ""

	insRepo := StockTransferMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.StockTransferDoc](prod, config.StockTransferMessageQueueConfig{}, "")
	return insRepo
}
