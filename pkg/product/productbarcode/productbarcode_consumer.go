package productbarcode

import (
	"context"
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/logger"
	ordertype_config "smlcloudplatform/pkg/product/ordertype/config"
	"smlcloudplatform/pkg/product/productbarcode/config"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/product/productbarcode/services"
	"smlcloudplatform/pkg/product/productbarcode/usecases"
	productgroup_config "smlcloudplatform/pkg/product/productgroup/config"
	producttype_config "smlcloudplatform/pkg/product/producttype/config"
	unit_config "smlcloudplatform/pkg/product/unit/config"
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
	ms                             *microservice.Microservice
	productMessageQueueConfig      config.ProductMessageQueueConfig
	productTypeMessageQueueConfig  producttype_config.ProductTypeMessageQueueConfig
	productGroupMessageQueueConfig productgroup_config.ProductGroupMessageQueueConfig
	unitMessageQueueConfig         unit_config.UnitMessageQueueConfig
	orderTypeMessageQueue          ordertype_config.OrderTypeMessageQueueConfig
	cfg                            msConfig.IConfig
	groupId                        string
	svc                            services.IProductBarcodeConsumeService
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
		ms:                            ms,
		cfg:                           cfg,
		productMessageQueueConfig:     config.ProductMessageQueueConfig{},
		productTypeMessageQueueConfig: producttype_config.ProductTypeMessageQueueConfig{},
		svc:                           svc,
		groupId:                       consumerGroupID,
	}
}

func (pbc *ProductBarcodeConsumer) RegisterConsumer() {

	// create topic
	mq := microservice.NewMQ(pbc.cfg.MQConfig(), pbc.ms.Logger)
	mq.CreateTopicR(pbc.productMessageQueueConfig.TopicCreated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(pbc.productMessageQueueConfig.TopicUpdated(), 5, 1, time.Hour*24*7)
	mq.CreateTopicR(pbc.productMessageQueueConfig.TopicDeleted(), 5, 1, time.Hour*24*7)

	// create topic for product type
	mq.CreateTopicR(pbc.productTypeMessageQueueConfig.TopicUpdated(), 5, 1, time.Hour*24*7)

	// create topic for product group
	mq.CreateTopicR(pbc.productGroupMessageQueueConfig.TopicUpdated(), 5, 1, time.Hour*24*7)

	// create topic for unit
	mq.CreateTopicR(pbc.unitMessageQueueConfig.TopicUpdated(), 5, 1, time.Hour*24*7)

	// create topic for order type
	mq.CreateTopicR(pbc.orderTypeMessageQueue.TopicUpdated(), 5, 1, time.Hour*24*7)

	pbc.ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productMessageQueueConfig.TopicCreated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductBarcodeCreate)
	pbc.ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productMessageQueueConfig.TopicUpdated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductBarcodeUpdate)
	pbc.ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productMessageQueueConfig.TopicDeleted(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductBarcodeDelete)

	pbc.ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productTypeMessageQueueConfig.TopicUpdated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductTypeUpdate)
	pbc.ms.Consume(pbc.cfg.MQConfig().URI(), pbc.productGroupMessageQueueConfig.TopicUpdated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductGroupUpdate)
	pbc.ms.Consume(pbc.cfg.MQConfig().URI(), pbc.unitMessageQueueConfig.TopicUpdated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnUnitUpdate)
	pbc.ms.Consume(pbc.cfg.MQConfig().URI(), pbc.orderTypeMessageQueue.TopicUpdated(), pbc.groupId, time.Duration(-1), pbc.ConsumerOnProductOrderTypeUpdate)
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
	moduleName := "Consumer On Product barcode Updated"

	pbc.ms.Logger.Debugf("Consume Product Barcode Update : %v", msg)
	doc := models.ProductBarcodeDoc{}
	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	err = pbc.svc.UpdateRefBarcode(doc.ShopID, doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	// _, err = pbc.svc.UpSert(doc.ShopID, doc.Barcode, doc)

	// if err != nil {
	// 	pbc.ms.Logger.Errorf(moduleName, err.Error())
	// }

	return nil
}

func (pbc *ProductBarcodeConsumer) ConsumerOnProductBarcodeDelete(ctx microservice.IContext) error {

	msg := ctx.ReadInput()
	moduleName := "Consumer On Product barcode Deleted"

	pbc.ms.Logger.Debugf("Consume Product Barcode Delete : %v", msg)
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

func (pbc *ProductBarcodeConsumer) ConsumerOnProductTypeUpdate(ctx microservice.IContext) error {

	msg := ctx.ReadInput()
	moduleName := "Consumer On Product Type Updated in Product Barcode"

	logger.GetLogger().Debugf("Consume Product Type Update in Product Barcode : %v", msg)
	doc := models.ProductTypeMessageQueueRequest{}
	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	err = pbc.svc.UpdateProductType(doc.ShopID, doc.ToProductType())

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	return nil
}

func (pbc *ProductBarcodeConsumer) ConsumerOnProductGroupUpdate(ctx microservice.IContext) error {

	msg := ctx.ReadInput()
	moduleName := "Consumer On Product Group Updated in Product Barcode"

	logger.GetLogger().Debugf("Consume Product Group Update in Product Barcode : %v", msg)
	doc := models.ProductGroupMessageQueueRequest{}
	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	err = pbc.svc.UpdateProductGroup(doc.ShopID, doc.ToProductGroup())

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	return nil
}

func (pbc *ProductBarcodeConsumer) ConsumerOnUnitUpdate(ctx microservice.IContext) error {

	msg := ctx.ReadInput()
	moduleName := "Consumer On Unit Updated in Product Barcode"

	logger.GetLogger().Debugf("Consume Unit Update in Product Barcode : %v", msg)
	doc := models.ProductUnitMessageQueueRequest{}
	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	err = pbc.svc.UpdateProductUnit(doc.ShopID, doc.ToProductUnit())

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	return nil
}

func (pbc *ProductBarcodeConsumer) ConsumerOnProductOrderTypeUpdate(ctx microservice.IContext) error {

	msg := ctx.ReadInput()
	moduleName := "Consumer On Product Type Updated in Product Barcode"

	logger.GetLogger().Debugf("Consume Product Type Update in Product Barcode : %v", msg)
	doc := models.ProductOrderTypeMessageQueueRequest{}
	err := json.Unmarshal([]byte(msg), &doc)

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	err = pbc.svc.UpdateProductOrderType(doc.ShopID, doc.ToProductOrderType())

	if err != nil {
		pbc.ms.Logger.Errorf(moduleName, err.Error())
	}

	return nil
}
