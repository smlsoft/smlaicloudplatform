package transaction

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"time"
)

func StartTransactionComsume(ms *microservice.Microservice, cfg microservice.IConfig) {
	topic := "when-transaction-created"
	groupID := "elk-transaction-consumer"
	timeout := time.Duration(-1)

	mqServer := cfg.MQServer()

	mq := microservice.NewMQ(mqServer, ms)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	ms.Consume(mqServer, topic, groupID, timeout, func(ctx microservice.IContext) error {
		elkPst := ms.ElkPersister(cfg.ElkPersisterConfig())

		msg := ctx.ReadInput()

		trans := models.TransactionRequest{}
		err := json.Unmarshal([]byte(msg), &trans)

		if err != nil {
			ms.Log("transaction comsume", err.Error())
		}

		err = elkPst.CreateWithId(trans.GuidFixed, &trans)

		if err != nil {
			ms.Log("transaction comsume", err.Error())
		}
		return nil
	})
}
