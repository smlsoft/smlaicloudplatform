package purchase

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"time"
)

func StartPurchaseComsume(ms *microservice.Microservice, cfg microservice.IConfig) {
	topic := MQ_TOPIC_PURCHASE_CREATED
	groupID := "elk-purchase-consumer"
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topic, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume purchase created"
		elkPst := ms.ElkPersister(cfg.ElkPersisterConfig())

		msg := ctx.ReadInput()

		doc := models.PurchaseRequest{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		err = elkPst.CreateWithID(doc.GuidFixed, &doc)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}
		return nil
	})
}
