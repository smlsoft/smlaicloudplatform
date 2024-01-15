package saleinvoicereturn

import (
	pkgConfig "smlcloudplatform/internal/config"
	"smlcloudplatform/internal/transaction/models"
	saleInvoiceReturnConfig "smlcloudplatform/internal/transaction/saleinvoicereturn/config"
	"smlcloudplatform/internal/transaction/transactionconsumer/debtortransaction"
	"smlcloudplatform/internal/transaction/transactionconsumer/services"
	"smlcloudplatform/internal/transaction/transactionconsumer/stocktransaction"
	"smlcloudplatform/internal/transaction/transactionconsumer/usecases"
	"smlcloudplatform/pkg/microservice"
	"time"
)

// import (
// 	"smlcloudplatform/internal/transaction/transactionconsumer/services"
// 	"smlcloudplatform/internal/transaction/transactionconsumer/usecases"
// 	"smlcloudplatform/pkg/microservice"
// )

// type ITransactionSaleInvoiceReturnConsumer interface {
// 	ConsumeOnCreateOrUpdate(ctx microservice.IContext) error
// 	ConsumeOnBulkCreateOrUpdate(ctx microservice.IContext) error
// 	ConsumeOnDelete(ctx microservice.IContext) error
// 	ConsumeOnBulkDelete(ctx microservice.IContext) error
// }

type SaleInvoiceReturnTransactionConsumer struct {
	ms                    *microservice.Microservice
	cfg                   pkgConfig.IConfig
	svc                   ISaleInvoiceReturnTransactionConsumerService
	transactionPhaser     usecases.ITransactionPhaser[models.SaleInvoiceReturnTransactionPG]
	stockPhaser           usecases.IStockTransactionPhaser[models.SaleInvoiceReturnTransactionPG]
	debtorPhaser          usecases.IDebtorTransactionPhaser[models.SaleInvoiceReturnTransactionPG]
	stockConsumerService  stocktransaction.IStockTransactionConsumerService
	debtorConsumerService debtortransaction.IDebtorTransactionConsumerService
}

func NewSaleInvoiceReturnTransactionConsumer(
	ms *microservice.Microservice,
	cfg pkgConfig.IConfig,
	svc ISaleInvoiceReturnTransactionConsumerService,
	stockConsumerService stocktransaction.IStockTransactionConsumerService,
	debtorConsumerService debtortransaction.IDebtorTransactionConsumerService,
) services.ITransactionDocConsumer {

	transactionPhaser := SaleInvoiceReturnTransactionPhaser{}
	saleInvoiceStockPhaser := SaleInvoiceReturnTransactionStockPhaser{}
	saleInvoiceDebtorPhaser := SaleInvoiceReturnDebtorTransactionPhaser{}

	return &SaleInvoiceReturnTransactionConsumer{
		ms:                    ms,
		cfg:                   cfg,
		svc:                   svc,
		transactionPhaser:     transactionPhaser,
		stockPhaser:           saleInvoiceStockPhaser,
		debtorPhaser:          saleInvoiceDebtorPhaser,
		stockConsumerService:  stockConsumerService,
		debtorConsumerService: debtorConsumerService,
	}
}

func (t *SaleInvoiceReturnTransactionConsumer) RegisterConsumer(ms *microservice.Microservice) {

	trxConsumerGroup := pkgConfig.GetEnv("TRANSACTION_CONSUMER_GROUP", "transaction-consumer-group-01")
	mq := microservice.NewMQ(t.cfg.MQConfig(), ms.Logger)
	saleInvoiceReturnKafkaConfig := saleInvoiceReturnConfig.SaleInvoiceReturnMessageQueueConfig{}

	mq.CreateTopicR(saleInvoiceReturnKafkaConfig.TopicCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceReturnKafkaConfig.TopicUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceReturnKafkaConfig.TopicDeleted(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceReturnKafkaConfig.TopicBulkCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceReturnKafkaConfig.TopicBulkUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(saleInvoiceReturnKafkaConfig.TopicBulkDeleted(), 5, 1, time.Hour*24*7)

	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceReturnKafkaConfig.TopicCreated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceReturnKafkaConfig.TopicUpdated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceReturnKafkaConfig.TopicDeleted(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnDelete)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceReturnKafkaConfig.TopicBulkCreated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnBulkCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceReturnKafkaConfig.TopicBulkUpdated(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnBulkCreateOrUpdate)
	ms.Consume(t.cfg.MQConfig().URI(), saleInvoiceReturnKafkaConfig.TopicBulkDeleted(), trxConsumerGroup, time.Duration(-1), t.ConsumeOnBulkDelete)

}

func InitSaleInvoiceReturnTransactionConsumer(ms *microservice.Microservice, cfg pkgConfig.IConfig) services.ITransactionDocConsumer {

	persister := ms.Persister(cfg.PersisterConfig())
	producer := ms.Producer(cfg.MQConfig())

	stockService := stocktransaction.NewStockTransactionConsumerService(persister, producer)
	debtorService := debtortransaction.NewDebtorTransactionService(persister, producer)

	saleInvoiceReturnConsumerService := NewSaleInvoiceReturnTransactionConsumerService(NewSaleInvoiceReturnTransactionPGRepository(persister))
	consumer := NewSaleInvoiceReturnTransactionConsumer(ms, cfg, saleInvoiceReturnConsumerService, stockService, debtorService)

	return consumer
}

func (t *SaleInvoiceReturnTransactionConsumer) ConsumeOnCreateOrUpdate(ctx microservice.IContext) error {
	msg := ctx.ReadInput()
	transaction, err := t.transactionPhaser.PhaseSingleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	err = t.svc.Upsert(transaction.ShopID, transaction.DocNo, *transaction)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}

	hasCreditorEffectDoc := transaction.HasDebtorEffectDoc()
	hasStockEffectDoc := transaction.HasStockEffectDoc()

	if hasStockEffectDoc {
		stockTransaction, err := t.stockPhaser.PhaseSingleDoc(*transaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Phase SaleInvoice Return to StockTransaction : %v", err.Error())
			return err
		}

		err = t.stockConsumerService.Upsert(stockTransaction.ShopID, stockTransaction.DocNo, *stockTransaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert SaleInvoice Return StockTransaction : %v", err.Error())
			return err
		}
	}

	if hasCreditorEffectDoc {
		debtorTransaction, err := t.debtorPhaser.PhaseSingleDoc(*transaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Phase SaleInvoice Return to DebtorTransaction : %v", err.Error())
			return err
		}

		err = t.debtorConsumerService.Upsert(debtorTransaction.ShopID, debtorTransaction.DocNo, *debtorTransaction)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert SaleInvoice Return DebtorTransaction : %v", err.Error())
			return err
		}
	}

	return nil
}

func (t *SaleInvoiceReturnTransactionConsumer) ConsumeOnDelete(ctx microservice.IContext) error {
	msg := ctx.ReadInput()

	transaction, err := t.transactionPhaser.PhaseSingleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	err = t.svc.Delete(transaction.ShopID, transaction.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}

	err = t.stockConsumerService.Delete(transaction.ShopID, transaction.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}

	err = t.debtorConsumerService.Delete(transaction.ShopID, transaction.DocNo)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
		return err
	}
	return nil
}

func (t *SaleInvoiceReturnTransactionConsumer) ConsumeOnBulkCreateOrUpdate(ctx microservice.IContext) error {
	msg := ctx.ReadInput()
	transactions, err := t.transactionPhaser.PhaseMultipleDoc(msg)
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

		hasCreditorEffectDoc := transaction.HasDebtorEffectDoc()
		hasStockEffectDoc := transaction.HasStockEffectDoc()

		if hasStockEffectDoc {
			stockTransaction, err := t.stockPhaser.PhaseSingleDoc(transaction)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Phase SaleInvoice Return to StockTransaction : %v", err.Error())
				return err
			}

			err = t.stockConsumerService.Upsert(stockTransaction.ShopID, stockTransaction.DocNo, *stockTransaction)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Insert SaleInvoice Return StockTransaction : %v", err.Error())
				return err
			}
		}

		if hasCreditorEffectDoc {
			debtorTransaction, err := t.debtorPhaser.PhaseSingleDoc(transaction)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Phase SaleInvoice Return to DebtorTransaction : %v", err.Error())
				return err
			}

			err = t.debtorConsumerService.Upsert(debtorTransaction.ShopID, debtorTransaction.DocNo, *debtorTransaction)
			if err != nil {
				t.ms.Logger.Errorf("Cannot Insert SaleInvoice Return DebtorTransaction : %v", err.Error())
				return err
			}
		}
	}
	return nil
}

func (t *SaleInvoiceReturnTransactionConsumer) ConsumeOnBulkDelete(ctx microservice.IContext) error {
	msg := ctx.ReadInput()
	transactions, err := t.transactionPhaser.PhaseMultipleDoc(msg)
	if err != nil {
		t.ms.Logger.Errorf("Cannot Phase PurchaseDoc to StockTransaction : %v", err.Error())
		return err
	}

	for _, transaction := range *transactions {
		err = t.svc.Delete(transaction.ShopID, transaction.DocNo)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
			return err
		}

		err = t.stockConsumerService.Delete(transaction.ShopID, transaction.DocNo)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
			return err
		}

		err = t.debtorConsumerService.Delete(transaction.ShopID, transaction.DocNo)
		if err != nil {
			t.ms.Logger.Errorf("Cannot Insert StockTransaction : %v", err.Error())
			return err
		}
	}
	return nil
}

func MigrationDatabase(ms *microservice.Microservice, cfg pkgConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.SaleInvoiceReturnTransactionPG{},
		models.SaleInvoiceReturnTransactionDetailPG{},
	)
	return nil
}
