package chartofaccount

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/vfgl"
	"time"
)

func MigrationChartOfAccountTable(ms *microservice.Microservice, cfg microservice.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		vfgl.ChartOfAccountPG{},
	)
	return nil
}

func StartChartOfAccountConsumerCreated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := MQ_TOPIC_CHARTOFACCOUNT_CREATED
	timeout := time.Duration(-1)
	mqConfig := cfg.MQConfig()
	mq := microservice.NewMQ(mqConfig, ms.Logger)
	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)
	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume chartofaccount created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Journal Created : %v", msg)
		doc := vfgl.ChartOfAccountDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		pgDocList := []vfgl.ChartOfAccountPG{}

		tmpJsonDoc, err := json.Marshal(doc)
		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		tmpDoc := vfgl.ChartOfAccountPG{}
		err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		pgDocList = append(pgDocList, tmpDoc)

		err = pst.CreateInBatch(pgDocList, len(pgDocList))

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}
		return nil
	})

}

func StartChartOfAccountConsumerBlukCreated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := MQ_TOPIC_CHARTOFACCOUNT_BULK_CREATED
	timeout := time.Duration(-1)
	mqConfig := cfg.MQConfig()
	mq := microservice.NewMQ(mqConfig, ms.Logger)
	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)
	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume chartofaccount created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()
		ms.Logger.Debugf("Consume Journal Created : %v", msg)

		docList := []vfgl.ChartOfAccountDoc{}
		err := json.Unmarshal([]byte(msg), &docList)
		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		pgDocList := []vfgl.ChartOfAccountPG{}
		for _, doc := range docList {
			tmpJsonDoc, err := json.Marshal(doc)
			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}
			tmpDoc := vfgl.ChartOfAccountPG{}
			err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}
			pgDocList = append(pgDocList, tmpDoc)
		}

		err = pst.CreateInBatch(pgDocList, len(pgDocList))
		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}
		return nil
	})

}
