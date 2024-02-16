package stockbalance

import (
	"context"
	stockBalanceModels "smlcloudplatform/internal/transaction/stockbalance/models"
	stockReceiveProductRepositories "smlcloudplatform/internal/transaction/stockbalance/repositories"
	"smlcloudplatform/pkg/microservice"
	"time"
)

type IStockBalanceProductTransactionAdminService interface {
	ReSyncStockBalanceProductDoc(shopID string) error
	ReSyncStockBalanceProductDeleteDoc(shopID string) error
}

type StockBalanceProductTransactionAdminService struct {
	mongoRepo       IStockBalanceTransactionAdminRepository
	kafkaRepo       stockReceiveProductRepositories.IStockBalanceMessageQueueRepository
	timeoutDuration time.Duration
}

func NewStockBalanceProductTransactionAdminService(
	pst microservice.IPersisterMongo,
	kfProducer microservice.IProducer,
) IStockBalanceProductTransactionAdminService {

	mongoRepo := NewStockBalanceTransactionAdminRepository(pst)
	kafkaRepo := stockReceiveProductRepositories.NewStockBalanceMessageQueueRepository(kfProducer)

	return &StockBalanceProductTransactionAdminService{
		mongoRepo:       mongoRepo,
		kafkaRepo:       kafkaRepo,
		timeoutDuration: time.Duration(30) * time.Second,
	}
}

func (s *StockBalanceProductTransactionAdminService) ReSyncStockBalanceProductDoc(shopID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), s.timeoutDuration)
	defer cancel()

	docs, err := s.mongoRepo.FindStockBalanceDocByShopID(ctx, shopID)
	if err != nil {
		return err
	}

	stockBalanceMessages := []stockBalanceModels.StockBalanceMessage{}

	for _, doc := range docs {
		stockBalanceMessage := stockBalanceModels.StockBalanceMessage{}
		stockBalanceMessage.StockBalance = doc.StockBalance

		stockBalanceMessages = append(stockBalanceMessages, stockBalanceMessage)
	}

	err = s.kafkaRepo.CreateInBatch(stockBalanceMessages)
	if err != nil {
		return err
	}

	return nil
}

func (s *StockBalanceProductTransactionAdminService) ReSyncStockBalanceProductDeleteDoc(shopID string) error {

	ctx, cancel := context.WithTimeout(context.Background(), s.timeoutDuration)
	defer cancel()

	docs, err := s.mongoRepo.FindStockBalanceDocDeleteByShopID(ctx, shopID)

	if err != nil {
		return err
	}

	stockBalanceMessages := []stockBalanceModels.StockBalanceMessage{}

	for _, doc := range docs {
		stockBalanceMessage := stockBalanceModels.StockBalanceMessage{}
		stockBalanceMessage.StockBalance = doc.StockBalance

		stockBalanceMessages = append(stockBalanceMessages, stockBalanceMessage)
	}

	err = s.kafkaRepo.DeleteInBatch(stockBalanceMessages)
	if err != nil {
		return err
	}
	return nil
}
