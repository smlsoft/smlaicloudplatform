package transactionconsumer

import (
	"smlcloudplatform/internal/microservice"
	pkgConfig "smlcloudplatform/pkg/config"
	"smlcloudplatform/pkg/transaction/models"
)

type ITransactionConsumer interface {
	ConsumeOnPurchaseDocCreated(doc interface{}) error
	ConsumeOnPurchaseDocUpdated(doc interface{}) error
	ConsumeOnPurchaseDocDeleted(doc interface{}) error

	// ConsumeOnSaleInvoiceDocCreated(doc interface{}) error
	// ConsumeOnSaleInvoiceDocUpdated(doc interface{}) error
	// ConsumeOnSaleInvoiceDocDeleted(doc interface{}) error
}

type TransactionConsumer struct {
	cfg pkgConfig.IConfig
	ms  *microservice.Microservice
}

func NewTransactionConsumer(ms *microservice.Microservice, cfg pkgConfig.IConfig) ITransactionConsumer {
	return &TransactionConsumer{
		cfg: cfg,
		ms:  ms,
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

func (t *TransactionConsumer) ConsumeOnPurchaseDocCreated(doc interface{}) error {
	return nil
}

func (t *TransactionConsumer) ConsumeOnPurchaseDocUpdated(doc interface{}) error {
	return nil
}

func (t *TransactionConsumer) ConsumeOnPurchaseDocDeleted(doc interface{}) error {
	return nil
}
