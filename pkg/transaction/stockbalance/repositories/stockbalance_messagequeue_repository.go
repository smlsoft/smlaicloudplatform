package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stockbalance/config"
	"smlcloudplatform/pkg/transaction/stockbalance/models"
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
