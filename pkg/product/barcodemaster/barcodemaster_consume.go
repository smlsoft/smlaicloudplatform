package barcodemaster

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/barcodemaster/config"
	"smlcloudplatform/pkg/product/barcodemaster/models"
	"time"
)

func StartBarcodeMasterComsumeCreated(ms *microservice.Microservice, cfg microservice.IConfig) {
	groupID := "postgres-barcodemaster-consumer"

	topicCreated := config.MQ_TOPIC_BARCODEMASTER_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)
	//Consume barcodemaster Created
	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume barcodemaster created"

		pst := ms.Persister(cfg.PersisterConfig())

		msg := ctx.ReadInput()

		inv := models.BarcodeMasterData{}
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

func startBarcodeMasterComsumeCreated(ms *microservice.Microservice, cfg microservice.IConfig) {
	groupID := "postgres-barcodemaster-consumer"

	topicCreated := config.MQ_TOPIC_BARCODEMASTER_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)
	//Consume barcodemaster Created
	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume barcodemaster created"

		pst := ms.Persister(cfg.PersisterConfig())

		msg := ctx.ReadInput()

		inv := models.BarcodeMasterData{}
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

func startBarcodeMasterComsumeUpdated(ms *microservice.Microservice, cfg microservice.IConfig) {
	groupID := "postgres-barcodemaster-consumer"

	topicCreated := config.MQ_TOPIC_BARCODEMASTER_UPDATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)
	//Consume barcodemaster Created
	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume barcodemaster updated"

		pst := ms.Persister(cfg.PersisterConfig())

		msg := ctx.ReadInput()

		inv := models.BarcodeMasterData{}
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
