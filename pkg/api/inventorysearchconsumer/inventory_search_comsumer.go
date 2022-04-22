package inventorysearchconsumer

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/models"
	"time"
)

type IInventorySearchConsumer interface {
}

const INVENTORYSEARCH_CONSUMER_GROUP_NAME string = "inventory-search-comsumer-group"

func StartInventorySearchComsumerOnProductCreated(ms *microservice.Microservice, cfg microservice.IConfig) {
	topic := inventory.MQ_TOPIC_INVENTORY_CREATED
	groupID := INVENTORYSEARCH_CONSUMER_GROUP_NAME

	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	//Consume transaction Created
	ms.Consume(mqConfig.URI(), topic, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume transaction created"

		elkPst := ms.ElkPersister(cfg.ElkPersisterConfig())

		msg := ctx.ReadInput()

		trans := models.InventorySearch{}
		err := json.Unmarshal([]byte(msg), &trans)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		err = elkPst.CreateWithID(trans.GuidFixed, &trans)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}
		return nil
	})

}
