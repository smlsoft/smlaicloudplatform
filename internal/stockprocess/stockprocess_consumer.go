package stockprocess

import (
	"encoding/json"
	"smlcloudplatform/internal/config"
	productBarcodeRepositories "smlcloudplatform/internal/product/productbarcode/repositories"
	stockProcessConfig "smlcloudplatform/internal/stockprocess/config"
	"smlcloudplatform/internal/stockprocess/models"
	"smlcloudplatform/internal/stockprocess/repositories"
	"smlcloudplatform/pkg/microservice"
	"time"
)

type IStockProcessConsumer interface {
	RegisterConsumer(ms *microservice.Microservice)

	ConsumerStockProcessCreate(ctx microservice.IContext) error
	ConsumerStockProcessBulkCreate(ctx microservice.IContext) error
}

type StockProcessConsumer struct {
	calculator IStockCalculator
	ms         *microservice.Microservice
	cfg        config.IConfig
}

func NewStockProcessConsumer(ms *microservice.Microservice, cfg config.IConfig) IStockProcessConsumer {

	repo := repositories.NewStockProcessPGRepository(ms.Persister(cfg.PersisterConfig()))
	barcodeRepo := productBarcodeRepositories.NewProductBarcodePGRepository(ms.Persister(cfg.PersisterConfig()))
	calculator := NewStockCalculator(repo, barcodeRepo)
	return &StockProcessConsumer{
		ms:         ms,
		cfg:        cfg,
		calculator: calculator,
	}
}

func (c *StockProcessConsumer) RegisterConsumer(ms *microservice.Microservice) {
	// ms.RegisterConsumer("stockprocess", "stockprocess", "stockprocess", c.ConsumeOnStockProcess)

	trxConsumerGroup := config.GetEnv("CONSUMER_GROUP_ID", "consumer-stockprocess-group-01")
	mq := microservice.NewMQ(c.cfg.MQConfig(), ms.Logger)

	stockProcessCfg := stockProcessConfig.StockProcessMessageQueueConfig{}

	mq.CreateTopicR(stockProcessCfg.TopicCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(stockProcessCfg.TopicUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(stockProcessCfg.TopicDeleted(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(stockProcessCfg.TopicBulkCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(stockProcessCfg.TopicBulkUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(stockProcessCfg.TopicBulkDeleted(), 5, 1, time.Hour*24*7)

	ms.Consume(c.cfg.MQConfig().URI(), stockProcessCfg.TopicCreated(), trxConsumerGroup, time.Duration(-1), c.ConsumerStockProcessCreate)
	ms.Consume(c.cfg.MQConfig().URI(), stockProcessCfg.TopicUpdated(), trxConsumerGroup, time.Duration(-1), c.ConsumerStockProcessCreate)
	ms.Consume(c.cfg.MQConfig().URI(), stockProcessCfg.TopicDeleted(), trxConsumerGroup, time.Duration(-1), c.ConsumerStockProcessCreate)
	ms.Consume(c.cfg.MQConfig().URI(), stockProcessCfg.TopicBulkCreated(), trxConsumerGroup, time.Duration(-1), c.ConsumerStockProcessBulkCreate)
	ms.Consume(c.cfg.MQConfig().URI(), stockProcessCfg.TopicBulkUpdated(), trxConsumerGroup, time.Duration(-1), c.ConsumerStockProcessBulkCreate)
	ms.Consume(c.cfg.MQConfig().URI(), stockProcessCfg.TopicBulkDeleted(), trxConsumerGroup, time.Duration(-1), c.ConsumerStockProcessBulkCreate)
}

func (c *StockProcessConsumer) ConsumerStockProcessCreate(ctx microservice.IContext) error {

	input := ctx.ReadInput()
	var stockProcessRequest models.StockProcessRequest
	err := json.Unmarshal([]byte(input), &stockProcessRequest)
	if err != nil {
		return err
	}

	err = c.calculator.CalculatorStock(stockProcessRequest.ShopID, stockProcessRequest.Barcode)
	return err
}

func (c *StockProcessConsumer) ConsumerStockProcessBulkCreate(ctx microservice.IContext) error {

	input := ctx.ReadInput()
	var stockProcessRequests []models.StockProcessRequest
	err := json.Unmarshal([]byte(input), &stockProcessRequests)
	if err != nil {
		return err
	}

	for _, req := range stockProcessRequests {
		go func(r models.StockProcessRequest) {
			err = c.calculator.CalculatorStock(r.ShopID, r.Barcode)
			if err != nil {
				c.ms.Logger.Errorf("Cannot Calculator Stock : %v", err.Error())
				//return err
			}
		}(req)
	}
	// err = c.calculator.CalculatorStock(stockProcessRequest.ShopID, stockProcessRequest.Barcode)
	return err
}
