package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stockreturnproduct/config"
	"smlcloudplatform/pkg/transaction/stockreturnproduct/models"
)

type IStockReturnProductMessageQueueRepository interface {
	Create(doc models.StockReturnProductDoc) error
	Update(doc models.StockReturnProductDoc) error
	Delete(doc models.StockReturnProductDoc) error
	CreateInBatch(docList []models.StockReturnProductDoc) error
	UpdateInBatch(docList []models.StockReturnProductDoc) error
	DeleteInBatch(docList []models.StockReturnProductDoc) error
}

type StockReturnProductMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.StockReturnProductDoc]
}

func NewStockReturnProductMessageQueueRepository(prod microservice.IProducer) StockReturnProductMessageQueueRepository {
	mqKey := ""

	insRepo := StockReturnProductMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.StockReturnProductDoc](prod, config.StockReturnProductMessageQueueConfig{}, "")
	return insRepo
}
