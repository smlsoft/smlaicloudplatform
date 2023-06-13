package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stocktransfer/config"
	"smlcloudplatform/pkg/transaction/stocktransfer/models"
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
