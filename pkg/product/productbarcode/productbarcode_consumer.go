package productbarcode

import (
	"context"
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/productbarcode/config"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/product/productbarcode/services"
	"smlcloudplatform/pkg/product/productbarcode/usecases"
	"time"

	msConfig "smlcloudplatform/pkg/config"
)

func MigrationProductBarcodeTable(ms *microservice.Microservice, cfg msConfig.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.ProductBarcodePg{},
	)
	return nil
}

type IProductBarcodeConsumer interface {
	ConsumerOnProductBarcodeCreate(ctx microservice.IContext) error
	ConsumerOnProductBarcodeUpdate(ctx microservice.IContext) error
	ConsumerOnProductBarcodeDelete(ctx microservice.IContext) error
	RegisterConsumer(ms *microservice.Microservice)
}

type ProductBarcodeConsumer struct {
	ms                        *microservice.Microservice
	productMessageQueueConfig config.ProductMessageQueueConfig
	cfg                       msConfig.IConfig
	groupId                   string
	svc                       services.IProductBarcodeConsumeService
}

func NewProductBarcodeConsumer(ms *microservice.Microservice, cfg msConfig.IConfig) *ProductBarcodeConsumer {

	consumerGroupID := msConfig.GetEnv("CONSUMER_GROUP_ID", "consumer-productbarcode-group-01")

	pgPersister := ms.Persister(cfg.PersisterConfig())
	mongoPersister := ms.MongoPersister(cfg.MongoPersisterConfig())

	clickhouseCfg := cfg.ClickHouseConfig()

	var clickhouse microservice.IPersisterClickHouse

	if len(clickhouseCfg.ServerAddress()) > 0 {
		clickhouse = ms.ClickHousePersister(cfg.ClickHouseConfig())
	} else {
		clickhouse = nil
	}

	phaser := usecases.ProductBarcodePhaser{}

	svc := services.NewProductBarcodeConsumerService(pgPersister, mongoPersister, clickhouse, phaser)

	return &ProductBarcodeConsumer{
		ms:                        ms,
		cfg:                       cfg,
		productMessageQueueConfig: config.ProductMessageQueueConfig{},
		svc:                       svc,
		groupId:                   consumerGroupID,
	}
}

func (pbc *ProductBarcodeConsumer) RegisterConsumer(ms *microservice.Microservice) {

	// create topic
	mq := microservice.NewMQ(pbc.cfg.MQConfig(), ms.Logger)
	mq.CreateTopicR(pbc.productMessageQueueConfig.TopicCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(pbc.productMessageQueueConfig.TopicUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(pbc.productMessageQueueConfig.TopicDeleted(), 5, 1, time.Hour*24*7)

	ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productMessageQueueConfig.TopicCreated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductBarcodeCreate)
	ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productMessageQueueConfig.TopicUpdated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductBarcodeUpdate)
	ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productMessageQueueConfig.TopicDeleted(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductBarcodeDelete)
}

func (pbc *ProductBarcodeConsumer) ConsumerOnProductBarcodeCreate(ctx microservice.IContext) error {
	msg := ctx.ReadInput()
	moduleName := "Consumer On Product barcode Created"

	pbc.ms.Logger.Debugf("Consume Product Barcode Create : %v", msg)
	doc := models.ProductBarcodeDoc{}
	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	_, err = pbc.svc.UpSert(doc.ShopID, doc.Barcode, doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}
	return nil
}

func (pbc *ProductBarcodeConsumer) ConsumerOnProductBarcodeUpdate(ctx microservice.IContext) error {

	msg := ctx.ReadInput()
	moduleName := "Consumer On Product barcode Created"

	pbc.ms.Logger.Debugf("Consume Product Barcode Create : %v", msg)
	doc := models.ProductBarcodeDoc{}
	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	_, err = pbc.svc.UpSert(doc.ShopID, doc.Barcode, doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	return nil
}

func (pbc *ProductBarcodeConsumer) ConsumerOnProductBarcodeDelete(ctx microservice.IContext) error {

	msg := ctx.ReadInput()
	moduleName := "Consumer On Product barcode Created"

	pbc.ms.Logger.Debugf("Consume Product Barcode Create : %v", msg)
	doc := models.ProductBarcodeDoc{}
	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	err = pbc.svc.Delete(context.Background(), doc.ShopID, doc.Barcode)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}
	return nil
}
