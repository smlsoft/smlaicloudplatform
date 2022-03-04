package transaction

import (
	"encoding/json"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"time"
)

func StartTransactionComsume(ms *microservice.Microservice, cfg microservice.IConfig) {
	topic := "when-transaction-created"
	groupID := "distribute-consumer"
	timeout := time.Duration(-1)

	mqServer := cfg.MQServer()

	mq := microservice.NewMQ(mqServer, ms)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	ms.Consume(mqServer, topic, groupID, timeout, func(ctx microservice.IContext) error {

		// pst := ms.Persister(cfg.PersisterConfig())

		msg := ctx.ReadInput()

		trans := models.Transaction{}
		err := json.Unmarshal([]byte(msg), &trans)

		if err != nil {
			ms.Log("transaction comsume", err.Error())
		}

		fmt.Printf("receive :: %v \n", trans)

		// err = pst.Create(trans)

		// if err != nil {
		// 	ms.Log("transaction comsume", err.Error())
		// }
		return nil
	})
}
