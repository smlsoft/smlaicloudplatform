package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/transaction/stockbalancedetail/config"
	"smlaicloudplatform/internal/transaction/stockbalancedetail/models"
	"smlaicloudplatform/pkg/microservice"
)

type IStockBalanceDetailMessageQueueRepository interface {
	Create(doc models.StockBalanceDetailDoc) error
	Update(doc models.StockBalanceDetailDoc) error
	Delete(doc models.StockBalanceDetailDoc) error
	CreateInBatch(docList []models.StockBalanceDetailDoc) error
	UpdateInBatch(docList []models.StockBalanceDetailDoc) error
	DeleteInBatch(docList []models.StockBalanceDetailDoc) error
}

type StockBalanceDetailMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.StockBalanceDetailDoc]
}

func NewStockBalanceDetailMessageQueueRepository(prod microservice.IProducer) StockBalanceDetailMessageQueueRepository {
	mqKey := ""

	insRepo := StockBalanceDetailMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.StockBalanceDetailDoc](prod, config.StockBalanceDetailMessageQueueConfig{}, "")
	return insRepo
}
