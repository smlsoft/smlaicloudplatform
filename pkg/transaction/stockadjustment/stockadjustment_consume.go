package stockadjustment

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/transaction/stockadjustment/models"
	"time"
)

func StartStockAdjustmentComsume(ms *microservice.Microservice, cfg microservice.IConfig) {
	topic := MQ_TOPIC_STOCK_ADJUSTMENT_CREATED
	groupID := "elk-stockadjustment-consumer"
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topic, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume stockadjustment created"
		elkPst := ms.ElkPersister(cfg.ElkPersisterConfig())

		msg := ctx.ReadInput()

		doc := models.StockAdjustmentData{}
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
