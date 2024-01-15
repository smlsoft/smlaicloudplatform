package purchasereturn

import (
	pkgConfig "smlcloudplatform/internal/config"
	"smlcloudplatform/internal/transaction/models"
	purchaseReturnConfig "smlcloudplatform/internal/transaction/purchasereturn/config"
	"smlcloudplatform/internal/transaction/transactionconsumer/creditortransaction"
	"smlcloudplatform/internal/transaction/transactionconsumer/services"
	"smlcloudplatform/internal/transaction/transactionconsumer/stocktransaction"
	"smlcloudplatform/internal/transaction/transactionconsumer/usecases"
	"smlcloudplatform/pkg/microservice"
	"time"
)

type IPurchaseReturnTransactionConsumer struct {
	ms                      *microservice.Microservice
	cfg                     pkgConfig.IConfig
	svc                     IPurchaseReturnTransactionConsumerService
	trxPhaser               usecases.ITransactionPhaser[models.PurchaseReturnTransactionPG]
	stockPhaser             usecases.IStockTransactionPhaser[models.PurchaseReturnTransactionPG]
	creditorPhaser          usecases.ICreditorTransactionPhaser[models.PurchaseReturnTransactionPG]
	stockConsumerService    stocktransaction.IStockTransactionConsumerService
	creditorConsumerService creditortransaction.ICreditorTransactionConsumerService
}

func NewTransactionPurchaseReturnConsumer(
	ms *microservice.Microservice,
	cfg pkgConfig.IConfig,
	svc IPurchaseReturnTransactionConsumerService,
	stockConsumerService stocktransaction.IStockTransactionConsumerService,
	creditorConsumerService creditortransaction.ICreditorTransactionConsumerService,
) services.ITransactionDocConsumer {

	purchaseReturnPhaser := PurchaseReturnTransactionPhaser{}
	purchaseReturnStockPhaser := PurchaseReturnTransactionStockPhaser{}
	purchaseReturnCreditorPhaser := PurchaseReturnTransactionCreditorPhaser{}

	return &IPurchaseReturnTransactionConsumer{
		ms:                      ms,
		cfg:                     cfg,
		svc:                     svc,
		trxPhaser:               purchaseReturnPhaser,
		stockPhaser:             purchaseReturnStockPhaser,
		creditorPhaser:          purchaseReturnCreditorPhaser,
		stockConsumerService:    stockConsumerService,
		creditorConsumerService: creditorConsumerService,
	}
}

func InitPurchaseReturnTransactionConsumer(ms *microservice.Microservice,
	cfg pkgConfig.IConfig) services.ITransactionDocConsumer {

	persister := ms.Persister(cfg.PersisterConfig())
	producer := ms.Producer(cfg.MQConfig())

	stockService := stocktransaction.NewStockTransactionConsumerService(persister, producer)
	creditorService := creditortransaction.NewCreditorTransactionConsumerService(persister, producer)

	purchaseReturnConsumerService := NewPurchaseReturnTransactionService(NewPurchaseReturnTransactionPGRepository(persister))
	consumer := NewTransactionPurchaseReturnConsumer(ms, cfg, purchaseReturnConsumerService, stockService, creditorService)

	return consumer
}

func (t *IPurchaseReturnTransactionConsumer) RegisterConsumer(ms *microservice.Microservice) {
	trxConsumerGroup := pkgConfig.GetEnv("TRANSACTION_CONSUMER_GROUP", "transaction-consumer-group-01")
	mq := microservice.NewMQ(t.cfg.MQConfig(), ms.Logger)

	purchaseKafkaConfig := purchaseReturnConfig.PurchaseReturnMessageQueueConfig{}

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

func (t *IPurchaseReturnTransactionConsumer) ConsumeOnCreateOrUpdate(ctx microservice.IContext) error {
	msg := ctx.ReadInput()
	transaction, err := t.trxPhaser.PhaseSingleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	err = t.svc.Upsert(transaction.ShopID, transaction.DocNo, *transaction)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}

	hasStockEffectDoc := transaction.HasStockEffectDoc()
	hasCreditorEffectDoc := transaction.HasCreditorEffectDoc()

	// Upsert stock transaction
	if hasStockEffectDoc {
		stockTrx, err := t.stockPhaser.PhaseSingleDoc(*transaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
			return err
		}

		isTransactionStockEffected := stockTrx.InquiryType == 0 || stockTrx.InquiryType == 2
		if isTransactionStockEffected {
			err = t.stockConsumerService.Upsert(stockTrx.ShopID, stockTrx.DocNo, *stockTrx)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
				return err
			}
		}
	}

	// Upsert creditor transaction
	if hasCreditorEffectDoc {
		creditorTrx, err := t.creditorPhaser.PhaseSingleDoc(*transaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to CreditorTransaction : %v", err.Error())
			return err
		}

		err = t.creditorConsumerService.Upsert(creditorTrx.ShopID, creditorTrx.DocNo, *creditorTrx)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert CreditorTransaction : %v", err.Error())
			return err
		}
	}

	return nil
}

func (t *IPurchaseReturnTransactionConsumer) ConsumeOnDelete(ctx microservice.IContext) error {
	msg := ctx.ReadInput()

	trx, err := t.trxPhaser.PhaseSingleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

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
		t.ms.Logger.Errorf("Cannot Delete CreditorTransaction : %v", err.Error())
		return err
	}

	return nil
}

func (t *IPurchaseReturnTransactionConsumer) ConsumeOnBulkCreateOrUpdate(ctx microservice.IContext) error {
	msg := ctx.ReadInput()
	transactions, err := t.trxPhaser.PhaseMultipleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	for _, transaction := range *transactions {
		err = t.svc.Upsert(transaction.ShopID, transaction.DocNo, transaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
			return err
		}

		hasStockEffectDoc := transaction.HasStockEffectDoc()
		hasCreditorEffectDoc := transaction.HasCreditorEffectDoc()

		// stock transaction
		if hasStockEffectDoc {
			stock, err := t.stockPhaser.PhaseSingleDoc(transaction)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
				return err
			}

			err = t.stockConsumerService.Upsert(stock.ShopID, stock.DocNo, *stock)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
				return err
			}

		}

		// creditor transaction
		if hasCreditorEffectDoc {
			creditor, err := t.creditorPhaser.PhaseSingleDoc(transaction)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to CreditorTransaction : %v", err.Error())
				return err
			}

			err = t.creditorConsumerService.Upsert(creditor.ShopID, creditor.DocNo, *creditor)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Insert CreditorTransaction : %v", err.Error())
				return err
			}
		}
	}

	return nil
}

func (t *IPurchaseReturnTransactionConsumer) ConsumeOnBulkDelete(ctx microservice.IContext) error {
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
			t.ms.Logger.Errorf("Cannot Delete CreditorTransaction : %v", err.Error())
			return err
		}
	}

	return nil
}

func MigrationDatabase(ms *microservice.Microservice, cfg pkgConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.PurchaseReturnTransactionPG{},
		models.PurchaseReturnTransactionDetailPG{},
	)
	return nil
}
