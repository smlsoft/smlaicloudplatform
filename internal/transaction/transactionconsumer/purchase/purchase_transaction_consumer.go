package purchase

import (
	pkgConfig "smlcloudplatform/internal/config"
	"smlcloudplatform/internal/logger"
	"smlcloudplatform/internal/transaction/models"
	purchaseConfig "smlcloudplatform/internal/transaction/purchase/config"
	"smlcloudplatform/internal/transaction/transactionconsumer/creditortransaction"
	"smlcloudplatform/internal/transaction/transactionconsumer/services"
	"smlcloudplatform/internal/transaction/transactionconsumer/stocktransaction"
	"smlcloudplatform/internal/transaction/transactionconsumer/usecases"
	"smlcloudplatform/pkg/microservice"
	"time"
)

type PurchaseTransactionConsumer struct {
	ms                      *microservice.Microservice
	cfg                     pkgConfig.IConfig
	svc                     IPurchaseTransactionConsumerService
	trxPhaser               usecases.ITransactionPhaser[models.PurchaseTransactionPG]
	stockPhaser             usecases.IStockTransactionPhaser[models.PurchaseTransactionPG]
	creditorPhaser          usecases.ICreditorTransactionPhaser[models.PurchaseTransactionPG]
	stockConsumerService    stocktransaction.IStockTransactionConsumerService
	creditorConsumerService creditortransaction.ICreditorTransactionConsumerService
}

func NewPurchaseTransactionConsumer(
	ms *microservice.Microservice,
	cfg pkgConfig.IConfig,
	svc IPurchaseTransactionConsumerService,
	stockConsumerService stocktransaction.IStockTransactionConsumerService,
	creditorConsumerService creditortransaction.ICreditorTransactionConsumerService,
) services.ITransactionDocConsumer {

	purchaseTrxPhaser := PurchaseTransactionPhaser{}
	purchaseStockPhaser := PurchaseTransactionStockPhaser{}
	purchaseCreditorPhaser := PurchaseCreditorTransactionPhaser{}

	return &PurchaseTransactionConsumer{
		ms:                      ms,
		cfg:                     cfg,
		svc:                     svc,
		trxPhaser:               purchaseTrxPhaser,
		stockPhaser:             purchaseStockPhaser,
		creditorPhaser:          purchaseCreditorPhaser,
		stockConsumerService:    stockConsumerService,
		creditorConsumerService: creditorConsumerService,
	}
}

func InitPurchaseTransactionConsumer(
	ms *microservice.Microservice,
	cfg pkgConfig.IConfig,
) services.ITransactionDocConsumer {

	persister := ms.Persister(cfg.PersisterConfig())
	producer := ms.Producer(cfg.MQConfig())

	stockService := stocktransaction.NewStockTransactionConsumerService(persister, producer)
	creditorService := creditortransaction.NewCreditorTransactionConsumerService(persister, producer)

	purchaseConsumerService := NewPurchaseTransactionService(NewPurchaseTransactionPGRepository(persister))
	consumer := NewPurchaseTransactionConsumer(ms, cfg, purchaseConsumerService, stockService, creditorService)

	return consumer
}

func (t *PurchaseTransactionConsumer) RegisterConsumer(ms *microservice.Microservice) {

	trxConsumerGroup := pkgConfig.GetEnv("TRANSACTION_CONSUMER_GROUP", "transaction-consumer-group-06")
	mq := microservice.NewMQ(t.cfg.MQConfig(), ms.Logger)

	purchaseKafkaConfig := purchaseConfig.PurchaseMessageQueueConfig{}

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

func (t *PurchaseTransactionConsumer) ConsumeOnCreateOrUpdate(ctx microservice.IContext) error {
	msg := ctx.ReadInput()

	transaction, err := t.trxPhaser.PhaseSingleDoc(msg)
	if err != nil {
		logger.GetLogger().Errorf("Cannot Phase PurchaseDoc to Purchase Transaction : %v", err.Error())
		return err
	}

	err = t.svc.Upsert(transaction.ShopID, transaction.DocNo, *transaction)
	if err != nil {
		logger.GetLogger().Errorf("Cannot Insert Purchase Transaction : %v", err.Error())
		return err
	}

	// upsert stock transaction
	if transaction.HasStockEffectDoc() {
		stock, err := t.stockPhaser.PhaseSingleDoc(*transaction)
		if err != nil {
			logger.GetLogger().Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
			return err
		}

		err = t.stockConsumerService.Upsert(transaction.ShopID, transaction.DocNo, *stock)
		if err != nil {
			logger.GetLogger().Errorf("Cannot Insert StockTransaction : %v", err.Error())
			return err
		}
	}

	if transaction.HasCreditorEffectDoc() {
		creditor, err := t.creditorPhaser.PhaseSingleDoc(*transaction)
		err = t.creditorConsumerService.Upsert(transaction.ShopID, transaction.DocNo, *creditor)
		if err != nil {
			logger.GetLogger().Errorf("Cannot Insert CreditorTransaction : %v", err.Error())
			return err
		}
	}

	return nil
}

func (t *PurchaseTransactionConsumer) ConsumeOnDelete(ctx microservice.IContext) error {

	msg := ctx.ReadInput()

	trx, err := t.trxPhaser.PhaseSingleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to Purchase Transaction : %v", err.Error())
		return err
	}

	err = t.svc.Delete(trx.ShopID, trx.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert Purchase Transaction : %v", err.Error())
		return err
	}

	// delete stock transaction
	err = t.stockConsumerService.Delete(trx.ShopID, trx.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Delete Stock Transaction : %v", err.Error())
		return err
	}

	// delete creditor transaction
	err = t.creditorConsumerService.Delete(trx.ShopID, trx.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Delete Creditor Transaction : %v", err.Error())
		return err
	}

	return nil
}

func (t *PurchaseTransactionConsumer) ConsumeOnBulkCreateOrUpdate(ctx microservice.IContext) error {
	msg := ctx.ReadInput()
	transactions, err := t.trxPhaser.PhaseMultipleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to Purchase Transaction : %v", err.Error())
		return err
	}

	for _, transaction := range *transactions {
		err = t.svc.Upsert(transaction.ShopID, transaction.DocNo, transaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert Purchase Transaction : %v", err.Error())
			return err
		}

		// stock transaction
		if transaction.HasStockEffectDoc() {
			stock, err := t.stockPhaser.PhaseSingleDoc(transaction)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to Stock Transaction : %v", err.Error())
				return err
			}

			err = t.stockConsumerService.Upsert(stock.ShopID, stock.DocNo, *stock)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Insert Stock Transaction : %v", err.Error())
				return err
			}
		}

		// creditor transaction
		hasCreditorEffectDoc := transaction.HasCreditorEffectDoc()
		if hasCreditorEffectDoc {
			creditor, err := t.creditorPhaser.PhaseSingleDoc(transaction)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to Creditor Transaction : %v", err.Error())
				return err
			}

			err = t.creditorConsumerService.Upsert(creditor.ShopID, creditor.DocNo, *creditor)
			if err != nil {
				logger.GetLogger().Errorf("Cannot Insert CreditorTransaction : %v", err.Error())
				return err
			}
		}
	}

	return nil
}

func (t *PurchaseTransactionConsumer) ConsumeOnBulkDelete(ctx microservice.IContext) error {
	msg := ctx.ReadInput()
	transactions, err := t.trxPhaser.PhaseMultipleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	for _, trx := range *transactions {
		err = t.svc.Delete(trx.ShopID, trx.DocNo)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
			return err
		}

		// delete stock transaction
		err = t.stockConsumerService.Delete(trx.ShopID, trx.DocNo)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Delete StockTransaction : %v", err.Error())
			return err
		}

		// delete creditor transaction
		err = t.creditorConsumerService.Delete(trx.ShopID, trx.DocNo)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Delete Creditor Transaction : %v", err.Error())
			return err
		}
	}

	return nil
}

func MigrationDatabase(ms *microservice.Microservice, cfg pkgConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.PurchaseTransactionPG{},
		models.PurchaseTransactionDetailPG{},
	)
	return nil
}
