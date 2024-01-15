package transactionconsumer

import (
	pkgConfig "smlcloudplatform/internal/config"
	"smlcloudplatform/internal/transaction/models"
	"smlcloudplatform/internal/transaction/transaction_consumer/services"
	"smlcloudplatform/internal/transaction/transaction_consumer/usecases"
	"smlcloudplatform/pkg/microservice"
)

type ITransactionConsumer interface {
	ConsumeOnPurchaseDocCreated(ctx microservice.IContext) error
	ConsumeOnPurchaseDocUpdated(ctx microservice.IContext) error
	ConsumeOnPurchaseDocDeleted(ctx microservice.IContext) error

	ConsumeOnSaleInvoiceDocCreated(ctx microservice.IContext) error
	ConsumeOnSaleInvoiceDocUpdated(ctx microservice.IContext) error
	ConsumeOnSaleInvoiceDocDeleted(ctx microservice.IContext) error
}

type TransactionConsumer struct {
	cfg    pkgConfig.IConfig
	ms     *microservice.Microservice
	svc    services.ITransactionConsumerService
	phaser usecases.ITransactionPhaser
}

func NewTransactionConsumer(ms *microservice.Microservice, cfg pkgConfig.IConfig) ITransactionConsumer {

	svc := services.NewTransactionConsumerService(ms.Persister(cfg.PersisterConfig()))
	phaser := &usecases.TransactionPhaser{}
	return &TransactionConsumer{
		cfg:    cfg,
		ms:     ms,
		svc:    svc,
		phaser: phaser,
	}

}

func MigrationDatabase(ms *microservice.Microservice, cfg pkgConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.StockTransaction{},
		models.StockTransactionDetail{},
	)
	return nil
}

func (pbc *TransactionConsumer) RegisterConsumer(ms *microservice.Microservice) {
	// mq := microservice.NewMQ(pbc.cfg.MQConfig(), ms.Logger)
	// mq.CreateTopicR(pbc.productMessageQueueConfig.TopicCreated(), 5, 1, time.Hour*24*7)
	// mq.CreateTopicR(pbc.productMessageQueueConfig.TopicUpdated(), 5, 1, time.Hour*24*7)
	// mq.CreateTopicR(pbc.productMessageQueueConfig.TopicDeleted(), 5, 1, time.Hour*24*7)

	// ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productMessageQueueConfig.TopicCreated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductBarcodeCreate)
	// ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productMessageQueueConfig.TopicUpdated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductBarcodeUpdate)
	// ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productMessageQueueConfig.TopicDeleted(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductBarcodeDelete)

}
