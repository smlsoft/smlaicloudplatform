package inventorysearchconsumer

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	inventoryConfig "smlcloudplatform/pkg/product/inventory/config"
	"smlcloudplatform/pkg/product/inventorysearchconsumer/models"
	"time"
)

const INVENTORYSEARCH_CONSUMER_GROUP_NAME string = "inventory-search-comsumer-group"

type IInventorySearchConsumer interface {
	ConsumeInsertInventory()
	ConsumeUpdateInventory()
	ConsumeDeleteInventory()
}

type InventorySearchConsumer struct {
	ms                             *microservice.Microservice
	cfg                            config.IConfig
	InventorySearchConsumerService IInventorySearchConsumerService
}

func NewInventorySearchConsumer(ms *microservice.Microservice, cfg config.IConfig) *InventorySearchConsumer {
	consume := &InventorySearchConsumer{
		ms:  ms,
		cfg: cfg,
	}
	openSearchPst := ms.SearchPersister(cfg.OpenSearchPersisterConfig())
	inventorySearchRepository := NewInventorySearchRepository(openSearchPst)
	consume.InventorySearchConsumerService = NewInventorySearchConsumerService(inventorySearchRepository)
	return consume
}

func (c *InventorySearchConsumer) Start() {
	c.ConsumeInsertInventory()
	c.ConsumeUpdateInventory()
	c.ConsumeDeleteInventory()
}

func (c *InventorySearchConsumer) ConsumeInsertInventory() {

	topic := inventoryConfig.MQ_TOPIC_INVENTORY_CREATED
	groupID := INVENTORYSEARCH_CONSUMER_GROUP_NAME

	timeout := time.Duration(-1)

	mqConfig := c.cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, c.ms.Logger)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	//Consume transaction Created
	c.ms.Consume(mqConfig.URI(), topic, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume transaction created"

		openSearchPst := c.ms.SearchPersister(c.cfg.OpenSearchPersisterConfig())

		msg := ctx.ReadInput()

		trans := models.InventorySearch{}
		err := json.Unmarshal([]byte(msg), &trans)

		if err != nil {
			c.ms.Log(moduleName, err.Error())
		}

		err = openSearchPst.CreateWithID(trans.GuidFixed, &trans)

		if err != nil {
			c.ms.Log(moduleName, err.Error())
		}
		return nil
	})
}

func (c *InventorySearchConsumer) ConsumeUpdateInventory() {
	topic := inventoryConfig.MQ_TOPIC_INVENTORY_UPDATED
	groupID := INVENTORYSEARCH_CONSUMER_GROUP_NAME

	timeout := time.Duration(-1)

	mqConfig := c.cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, c.ms.Logger)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	//Consume transaction Created
	c.ms.Consume(mqConfig.URI(), topic, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume transaction created"

		openSearchPst := c.ms.SearchPersister(c.cfg.OpenSearchPersisterConfig())

		msg := ctx.ReadInput()

		trans := models.InventorySearch{}
		err := json.Unmarshal([]byte(msg), &trans)

		if err != nil {
			c.ms.Log(moduleName, err.Error())
		}

		err = openSearchPst.Update(trans.GuidFixed, &trans)

		if err != nil {
			c.ms.Log(moduleName, err.Error())
		}
		return nil
	})
}

func (c *InventorySearchConsumer) ConsumeDeleteInventory() {
	topic := inventoryConfig.MQ_TOPIC_INVENTORY_DELETED
	groupID := INVENTORYSEARCH_CONSUMER_GROUP_NAME

	timeout := time.Duration(-1)

	mqConfig := c.cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, c.ms.Logger)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	//Consume transaction Created
	c.ms.Consume(mqConfig.URI(), topic, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume transaction created"

		openSearchPst := c.ms.SearchPersister(c.cfg.OpenSearchPersisterConfig())

		msg := ctx.ReadInput()

		trans := models.InventorySearch{}
		err := json.Unmarshal([]byte(msg), &trans)

		if err != nil {
			c.ms.Log(moduleName, err.Error())
		}

		err = openSearchPst.Delete(trans.GuidFixed, &trans)

		if err != nil {
			c.ms.Log(moduleName, err.Error())
		}
		return nil
	})
}
