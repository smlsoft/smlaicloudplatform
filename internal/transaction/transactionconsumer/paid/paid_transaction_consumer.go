package purchase

import (
	"encoding/json"
	pkgConfig "smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/logger"
	paidConfig "smlaicloudplatform/internal/transaction/paid/config"
	"smlaicloudplatform/internal/transaction/transactionconsumer/services"
	"smlaicloudplatform/pkg/microservice"
	"time"

	trans_models "smlaicloudplatform/internal/transaction/models"
	transaction_payment_consume "smlaicloudplatform/internal/transaction/transactionconsumer/payment"
)

type PaidTransactionConsumer struct {
	ms                         *microservice.Microservice
	cfg                        pkgConfig.IConfig
	transPaymentConsumeUsecase transaction_payment_consume.IPaymentUsecase
}

func NewPaidTransactionConsumer(
	ms *microservice.Microservice,
	cfg pkgConfig.IConfig,
	transPaymentConsumeUsecase transaction_payment_consume.IPaymentUsecase,
) services.ITransactionDocConsumer {

	return &PaidTransactionConsumer{
		ms:                         ms,
		cfg:                        cfg,
		transPaymentConsumeUsecase: transPaymentConsumeUsecase,
	}
}

func InitPaidTransactionConsumer(
	ms *microservice.Microservice,
	cfg pkgConfig.IConfig,
) services.ITransactionDocConsumer {

	persister := ms.Persister(cfg.PersisterConfig())

	transPaymentConsumeUsecase := transaction_payment_consume.InitPayment(persister)

	consumer := NewPaidTransactionConsumer(ms, cfg, transPaymentConsumeUsecase)

	return consumer
}

func (t *PaidTransactionConsumer) RegisterConsumer(ms *microservice.Microservice) {

	trxConsumerGroup := pkgConfig.GetEnv("TRANSACTION_CONSUMER_GROUP", "transaction-consumer-group-06")
	mq := microservice.NewMQ(t.cfg.MQConfig(), ms.Logger)

	purchaseKafkaConfig := paidConfig.DebtorPaymentMessageQueueConfig{}

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

func (t *PaidTransactionConsumer) ConsumeOnCreateOrUpdate(ctx microservice.IContext) error {
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

func (t *PaidTransactionConsumer) ConsumeOnDelete(ctx microservice.IContext) error {

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

func (t *PaidTransactionConsumer) ConsumeOnBulkCreateOrUpdate(ctx microservice.IContext) error {
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

func (t *PaidTransactionConsumer) ConsumeOnBulkDelete(ctx microservice.IContext) error {
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

func (t *PaidTransactionConsumer) upsertPayment(transMQDoc trans_models.TransactionMessageQueue) error {

	err := t.transPaymentConsumeUsecase.Upsert(transMQDoc)

	if err != nil {
		return err
	}

	return nil
}

func MigrationDatabase(ms *microservice.Microservice, cfg pkgConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		trans_models.PaidTransactionPG{},
		trans_models.PaidTransactionDetailPG{},
	)
	return nil
}
