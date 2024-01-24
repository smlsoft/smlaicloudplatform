package purchase

import (
	"encoding/json"
	pkgConfig "smlcloudplatform/internal/config"
	"smlcloudplatform/internal/logger"
	payConfig "smlcloudplatform/internal/transaction/pay/config"
	"smlcloudplatform/internal/transaction/transactionconsumer/services"
	"smlcloudplatform/pkg/microservice"
	"time"

	trans_models "smlcloudplatform/internal/transaction/models"
	transaction_payment_consume "smlcloudplatform/internal/transaction/transactionconsumer/payment"
)

type PayTransactionConsumer struct {
	ms                         *microservice.Microservice
	cfg                        pkgConfig.IConfig
	transPaymentConsumeUsecase transaction_payment_consume.IPaymentUsecase
}

func NewPayTransactionConsumer(
	ms *microservice.Microservice,
	cfg pkgConfig.IConfig,
	transPaymentConsumeUsecase transaction_payment_consume.IPaymentUsecase,
) services.ITransactionDocConsumer {

	return &PayTransactionConsumer{
		ms:                         ms,
		cfg:                        cfg,
		transPaymentConsumeUsecase: transPaymentConsumeUsecase,
	}
}

func InitPayTransactionConsumer(
	ms *microservice.Microservice,
	cfg pkgConfig.IConfig,
) services.ITransactionDocConsumer {

	persister := ms.Persister(cfg.PersisterConfig())

	transPaymentConsumeUsecase := transaction_payment_consume.InitPayment(persister)

	consumer := NewPayTransactionConsumer(ms, cfg, transPaymentConsumeUsecase)

	return consumer
}

func (t *PayTransactionConsumer) RegisterConsumer(ms *microservice.Microservice) {

	trxConsumerGroup := pkgConfig.GetEnv("TRANSACTION_CONSUMER_GROUP", "transaction-consumer-group-06")
	mq := microservice.NewMQ(t.cfg.MQConfig(), ms.Logger)

	purchaseKafkaConfig := payConfig.CreditorPaymentMessageQueueConfig{}

	mq.CreateTopicR(purchaseKafkaConfig.TopicCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(purchaseKafkaConfig.TopicUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(purchaseKafkaConfig.TopicDeleted(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(purchaseKafkaConfig.TopicBulkCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(purchaseKafkaConfig.TopicBulkUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(purchaseKafkaConfig.TopicBulkDeleted(), 5, 1, time.Hour*24*7)

	ms.Consume(t.cfg.MQConfig().URI(), purchaseKafkaConfig.TopicCreated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), purchaseKafkaConfig.TopicUpdated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), purchaseKafkaConfig.TopicDeleted(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnDelete)
	ms.Consume(t.cfg.MQConfig().URI(), purchaseKafkaConfig.TopicBulkCreated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnBulkCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), purchaseKafkaConfig.TopicBulkUpdated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnBulkCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), purchaseKafkaConfig.TopicBulkDeleted(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnBulkDelete)

}

func (t *PayTransactionConsumer) ConsumeOnCreateOrUpdate(ctx microservice.IContext) error {
	msg := ctx.ReadInput()

	transMQDoc := trans_models.TransactionMessageQueue{}
	err := json.Unmarshal([]byte(msg), &transMQDoc)
	if err != nil {
		logger.GetLogger().Errorf("Cannot Unmarshal Transaction Message Queue : %v", err.Error())
		return err
	}

	err = t.upsertPayment(transMQDoc)

	if err != nil {
		logger.GetLogger().Errorf("Cannot Upsert Transaction Payment : %v", err.Error())
		return err
	}
	return nil
}

func (t *PayTransactionConsumer) ConsumeOnDelete(ctx microservice.IContext) error {

	msg := ctx.ReadInput()

	transMQDoc := trans_models.TransactionMessageQueue{}
	err := json.Unmarshal([]byte(msg), &transMQDoc)
	if err != nil {
		logger.GetLogger().Errorf("Cannot Unmarshal Transaction Message Queue : %v", err.Error())
		return err
	}

	// delete transaction payment
	err = t.transPaymentConsumeUsecase.Delete(transMQDoc.ShopID, transMQDoc.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Delete Transaction Payment : %v", err.Error())
		return err
	}

	return nil
}

func (t *PayTransactionConsumer) ConsumeOnBulkCreateOrUpdate(ctx microservice.IContext) error {
	msg := ctx.ReadInput()

	transMQDoc := trans_models.TransactionMessageQueue{}
	err := json.Unmarshal([]byte(msg), &transMQDoc)
	if err != nil {
		logger.GetLogger().Errorf("Cannot Unmarshal Transaction Message Queue : %v", err.Error())
		return err
	}

	err = t.upsertPayment(transMQDoc)

	if err != nil {
		logger.GetLogger().Errorf("Cannot Upsert Transaction Payment : %v", err.Error())
		return err
	}
	return nil
}

func (t *PayTransactionConsumer) ConsumeOnBulkDelete(ctx microservice.IContext) error {
	msg := ctx.ReadInput()

	// transaction payment
	transMQDocs := []trans_models.TransactionMessageQueue{}

	err := json.Unmarshal([]byte(msg), &transMQDocs)

	if err != nil {
		logger.GetLogger().Errorf("Cannot Unmarshal Transaction Message Queue : %v", err.Error())
		return err
	}

	for _, transMQDoc := range transMQDocs {
		err = t.transPaymentConsumeUsecase.Delete(transMQDoc.ShopID, transMQDoc.DocNo)
		if err != nil {
			logger.GetLogger().Errorf("Cannot Delete Transaction Payment : %v", err.Error())
			return err
		}
	}

	return nil
}

func (t *PayTransactionConsumer) upsertPayment(transMQDoc trans_models.TransactionMessageQueue) error {

	err := t.transPaymentConsumeUsecase.Upsert(transMQDoc)

	if err != nil {
		return err
	}

	return nil
}

func MigrationDatabase(ms *microservice.Microservice, cfg pkgConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		trans_models.PayTransactionPG{},
		trans_models.PayTransactionDetailPG{},
	)
	return nil
}
