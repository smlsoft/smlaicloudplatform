package saleinvoice

import (
	"encoding/json"
	pkgConfig "smlcloudplatform/internal/config"
	"smlcloudplatform/internal/logger"
	"smlcloudplatform/internal/transaction/models"
	saleInvoiceConfig "smlcloudplatform/internal/transaction/saleinvoice/config"
	"smlcloudplatform/internal/transaction/transactionconsumer/debtortransaction"
	"smlcloudplatform/internal/transaction/transactionconsumer/services"
	"smlcloudplatform/internal/transaction/transactionconsumer/stocktransaction"
	"smlcloudplatform/internal/transaction/transactionconsumer/usecases"
	"smlcloudplatform/pkg/microservice"
	"time"

	trans_models "smlcloudplatform/internal/transaction/models"
	transaction_payment_consume "smlcloudplatform/internal/transaction/transactionconsumer/payment"
)

type SaleInvoiceTransactionConsumer struct {
	ms                         *microservice.Microservice
	cfg                        pkgConfig.IConfig
	svc                        ISaleInvoiceTransactionConsumerService
	trxPhaser                  usecases.ITransactionPhaser[models.SaleInvoiceTransactionPG]
	stockPhaser                usecases.IStockTransactionPhaser[models.SaleInvoiceTransactionPG]
	debtorPhaser               usecases.IDebtorTransactionPhaser[models.SaleInvoiceTransactionPG]
	stockConsumerService       stocktransaction.IStockTransactionConsumerService
	debtorConsumerService      debtortransaction.IDebtorTransactionConsumerService
	transPaymentConsumeUsecase transaction_payment_consume.IPaymentUsecase
}

func NewSaleInvoiceTransactionConsumer(
	ms *microservice.Microservice,
	cfg pkgConfig.IConfig,
	svc ISaleInvoiceTransactionConsumerService,
	stockConsumerService stocktransaction.IStockTransactionConsumerService,
	debtorConsumerService debtortransaction.IDebtorTransactionConsumerService,
	transPaymentConsumeUsecase transaction_payment_consume.IPaymentUsecase,
) services.ITransactionDocConsumer {

	saleInvoiceTrxPhaser := SalesInvoiceTransactionPhaser{}
	saleInvoiceStockPhaser := SaleInvoiceTransactionStockPhaser{}
	saleInvoiceDebtorPhaser := SaleInvoiceDebtorTransactionPhaser{}

	return &SaleInvoiceTransactionConsumer{
		ms:                         ms,
		cfg:                        cfg,
		svc:                        svc,
		trxPhaser:                  saleInvoiceTrxPhaser,
		stockPhaser:                saleInvoiceStockPhaser,
		debtorPhaser:               saleInvoiceDebtorPhaser,
		stockConsumerService:       stockConsumerService,
		debtorConsumerService:      debtorConsumerService,
		transPaymentConsumeUsecase: transPaymentConsumeUsecase,
	}
}

func (t *SaleInvoiceTransactionConsumer) RegisterConsumer(ms *microservice.Microservice) {
	trxConsumerGroup := pkgConfig.GetEnv("TRANSACTION_CONSUMER_GROUP", "transaction-consumer-group-01")
	mq := microservice.NewMQ(t.cfg.MQConfig(), ms.Logger)
	saleInvoiceKafkaConfig := saleInvoiceConfig.SaleInvoiceMessageQueueConfig{}

	mq.CreateTopicR(saleInvoiceKafkaConfig.TopicCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceKafkaConfig.TopicUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceKafkaConfig.TopicDeleted(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceKafkaConfig.TopicBulkCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceKafkaConfig.TopicBulkUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceKafkaConfig.TopicBulkDeleted(), 5, 1, time.Hour*24*7)

	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceKafkaConfig.TopicCreated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceKafkaConfig.TopicUpdated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceKafkaConfig.TopicDeleted(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnDelete)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceKafkaConfig.TopicBulkCreated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnBulkCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceKafkaConfig.TopicBulkUpdated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnBulkCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceKafkaConfig.TopicBulkDeleted(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnBulkDelete)

}

func InitSaleInvoiceTransactionConsumer(ms *microservice.Microservice, cfg pkgConfig.IConfig) services.ITransactionDocConsumer {

	persister := ms.Persister(cfg.PersisterConfig())
	producer := ms.Producer(cfg.MQConfig())

	stockService := stocktransaction.NewStockTransactionConsumerService(persister, producer)
	debtorService := debtortransaction.NewDebtorTransactionService(persister, producer)

	transPaymentConsumeUsecase := transaction_payment_consume.InitPayment(persister)

	saleInvoiceConsumerService := NewSaleInvoiceTransactionConsumerService(NewSaleInvoiceTransactionPGRepository(persister))
	consumer := NewSaleInvoiceTransactionConsumer(ms, cfg, saleInvoiceConsumerService, stockService, debtorService, transPaymentConsumeUsecase)

	return consumer
}

func (t *SaleInvoiceTransactionConsumer) ConsumeOnCreateOrUpdate(ctx microservice.IContext) error {

	msg := ctx.ReadInput()
	transaction, err := t.trxPhaser.PhaseSingleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	logger.GetLogger().Info("ConsumeOnCreateOrUpdate : %v", transaction.DocNo)

	err = t.svc.Upsert(transaction.ShopID, transaction.DocNo, *transaction)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}

	// upsert stock transaction
	stock, err := t.stockPhaser.PhaseSingleDoc(*transaction)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	err = t.stockConsumerService.Upsert(stock.ShopID, stock.DocNo, *stock)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}

	// upsert debtor transaction
	hasCreditorEffectDoc := transaction.HasCreditorEffectDoc()
	if hasCreditorEffectDoc {
		debtor, err := t.debtorPhaser.PhaseSingleDoc(*transaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
			return err
		}

		err = t.debtorConsumerService.Upsert(debtor.ShopID, debtor.DocNo, *debtor)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
			return err
		}
	}

	// upsert transaction payment
	transMQDoc := trans_models.TransactionMessageQueue{}

	err = json.Unmarshal([]byte(msg), &transMQDoc)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Unmarshal Transaction Message Queue : %v", err.Error())
		return err
	}

	err = t.upsertPayment(transMQDoc)
	if err != nil {
		logger.GetLogger().Errorf("Cannot Upsert Payment : %v", err.Error())
		return err
	}

	return nil
}

func (t *SaleInvoiceTransactionConsumer) ConsumeOnDelete(ctx microservice.IContext) error {

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
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}

	// delete debtor transaction
	err = t.debtorConsumerService.Delete(trx.ShopID, trx.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}

	// delete transaction payment
	err = t.transPaymentConsumeUsecase.Delete(trx.ShopID, trx.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Delete Transaction Payment : %v", err.Error())
		return err
	}

	return nil
}

func (t *SaleInvoiceTransactionConsumer) ConsumeOnBulkCreateOrUpdate(ctx microservice.IContext) error {
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

		// upsert stock transaction
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

		// upsert debtor transaction
		debtor, err := t.debtorPhaser.PhaseSingleDoc(transaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
			return err
		}

		if debtor.InquiryType == 0 {
			err = t.debtorConsumerService.Upsert(debtor.ShopID, debtor.DocNo, *debtor)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
				return err
			}
		}
	}

	transMQDocs := []trans_models.TransactionMessageQueue{}
	err = json.Unmarshal([]byte(msg), &transMQDocs)

	if err != nil {
		t.ms.Logger.Errorf("Cannot Unmarshal Transaction Message Queue : %v", err.Error())
		return err
	}

	for _, transMQDoc := range transMQDocs {
		err = t.upsertPayment(transMQDoc)
		if err != nil {
			logger.GetLogger().Errorf("Cannot Upsert Payment : %v", err.Error())
			return err
		}
	}

	return nil
}

func (t *SaleInvoiceTransactionConsumer) ConsumeOnBulkDelete(ctx microservice.IContext) error {
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
			t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
			return err
		}

		// delete debtor transaction
		err = t.debtorConsumerService.Delete(trx.ShopID, trx.DocNo)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
			return err
		}

		// delete transaction payment
		err = t.transPaymentConsumeUsecase.Delete(trx.ShopID, trx.DocNo)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Delete Transaction Payment : %v", err.Error())
			return err
		}
	}

	return nil
}

func (t *SaleInvoiceTransactionConsumer) upsertPayment(transMQDoc trans_models.TransactionMessageQueue) error {
	// transaction payment inquiryType = 1

	if t.HasPaymentEffectDoc(transMQDoc) {
		err := t.transPaymentConsumeUsecase.Upsert(transMQDoc)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *SaleInvoiceTransactionConsumer) HasPaymentEffectDoc(transDoc trans_models.TransactionMessageQueue) bool {
	return transDoc.InquiryType == 1
}

func MigrationDatabase(ms *microservice.Microservice, cfg pkgConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.SaleInvoiceTransactionPG{},
		models.SaleInvoiceTransactionDetailPG{},
	)
	return nil
}
