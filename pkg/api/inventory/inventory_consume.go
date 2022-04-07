package inventory

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"time"
)

func StartInventoryComsumeCreated(ms *microservice.Microservice, cfg microservice.IConfig) {
	groupID := "postgres-inventory-consumer"

	topicCreated := MQ_TOPIC_INVENTORY_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)
	//Consume inventory Created
	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume inventory created"

		pst := ms.Persister(cfg.PersisterConfig())

		msg := ctx.ReadInput()

		inv := models.InventoryData{}
		err := json.Unmarshal([]byte(msg), &inv)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		err = pst.Create(inv)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		return nil
	})

}

func startInventoryComsumeCreated(ms *microservice.Microservice, cfg microservice.IConfig) {
	groupID := "postgres-inventory-consumer"

	topicCreated := MQ_TOPIC_INVENTORY_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)
	//Consume inventory Created
	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume inventory created"

		pst := ms.Persister(cfg.PersisterConfig())

		msg := ctx.ReadInput()

		inv := models.InventoryData{}
		err := json.Unmarshal([]byte(msg), &inv)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		err = pst.Create(inv)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		return nil
	})

}

func startInventoryComsumeUpdated(ms *microservice.Microservice, cfg microservice.IConfig) {
	groupID := "postgres-inventory-consumer"

	topicCreated := MQ_TOPIC_INVENTORY_UPDATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)
	//Consume inventory Created
	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume inventory updated"

		pst := ms.Persister(cfg.PersisterConfig())

		msg := ctx.ReadInput()

		inv := models.InventoryData{}
		err := json.Unmarshal([]byte(msg), &inv)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		err = pst.Create(inv)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		return nil
	})

}
