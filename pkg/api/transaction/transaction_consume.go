package transaction

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"time"
)

func StartTransactionComsumeCreated(ms *microservice.Microservice, cfg microservice.IConfig) {
	topic := "when-transaction-created"
	groupID := "elk-transaction-consumer"
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	//Consume transaction Created
	ms.Consume(mqConfig.URI(), topic, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume transaction created"

		elkPst := ms.ElkPersister(cfg.ElkPersisterConfig())

		msg := ctx.ReadInput()

		trans := models.TransactionRequest{}
		err := json.Unmarshal([]byte(msg), &trans)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		err = elkPst.CreateWithId(trans.GuidFixed, &trans)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}
		return nil
	})

}
