package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stockreceiveproduct/config"
	"smlcloudplatform/pkg/transaction/stockreceiveproduct/models"
)

type IStockReceiveProductMessageQueueRepository interface {
	Create(doc models.StockReceiveProductDoc) error
	Update(doc models.StockReceiveProductDoc) error
	Delete(doc models.StockReceiveProductDoc) error
	CreateInBatch(docList []models.StockReceiveProductDoc) error
	UpdateInBatch(docList []models.StockReceiveProductDoc) error
	DeleteInBatch(docList []models.StockReceiveProductDoc) error
}

type StockReceiveProductMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.StockReceiveProductDoc]
}

func NewStockReceiveProductMessageQueueRepository(prod microservice.IProducer) StockReceiveProductMessageQueueRepository {
	mqKey := ""

	insRepo := StockReceiveProductMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.StockReceiveProductDoc](prod, config.StockReceiveProductMessageQueueConfig{}, "")
	return insRepo
}
