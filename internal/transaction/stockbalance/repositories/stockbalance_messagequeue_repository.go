package repositories

import (
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/internal/transaction/stockbalance/config"
	"smlcloudplatform/internal/transaction/stockbalance/models"
	"smlcloudplatform/pkg/microservice"
)

type IStockBalanceMessageQueueRepository interface {
	Create(doc models.StockBalanceDoc) error
	Update(doc models.StockBalanceDoc) error
	Delete(doc models.StockBalanceDoc) error
	CreateInBatch(docList []models.StockBalanceDoc) error
	UpdateInBatch(docList []models.StockBalanceDoc) error
	DeleteInBatch(docList []models.StockBalanceDoc) error
}

type StockBalanceMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.StockBalanceDoc]
}

func NewStockBalanceMessageQueueRepository(prod microservice.IProducer) StockBalanceMessageQueueRepository {
	mqKey := ""

	insRepo := StockBalanceMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.StockBalanceDoc](prod, config.StockBalanceMessageQueueConfig{}, "")
	return insRepo
}
