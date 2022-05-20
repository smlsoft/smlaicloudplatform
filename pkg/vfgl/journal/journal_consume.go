package journal

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/vfgl"
	"time"
)

func MigrationJournalTable(ms *microservice.Microservice, cfg microservice.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		vfgl.JournalPg{},
		vfgl.JournalDetailPg{},
	)
	return nil
}

func StartJournalComsumeCreated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := MQ_TOPIC_JOURNAL_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume journal created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Journal Created : %v", msg)
		doc := vfgl.JournalDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		pgDocList := []vfgl.JournalPg{}
		pgDocDetailList := []vfgl.JournalDetailPg{}

		//for _, doc := range docList
		{

			tmpJsonDoc, err := json.Marshal(doc)

			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}

			tmpDoc := vfgl.JournalPg{}
			err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}

			// err = json.Unmarshal([]byte(msg), &doc)

			// if err != nil {
			// 	ms.Logger.Errorf(moduleName, err.Error())
			// }

			pgDocList = append(pgDocList, tmpDoc)

			docDetailList := doc.AccountBook

			for _, docDetail := range docDetailList {
				tmpDocDetail := vfgl.JournalDetailPg{}

				tmpJsonDocDetail, err := json.Marshal(docDetail)
				if err != nil {
					ms.Logger.Errorf(moduleName, err.Error())
				}
				err = json.Unmarshal([]byte(tmpJsonDocDetail), &tmpDocDetail)

				if err != nil {
					ms.Logger.Errorf(moduleName, err.Error())
				}

				tmpDocDetail.ParID = doc.ParID
				tmpDocDetail.ShopID = doc.ShopID
				tmpDocDetail.Docno = doc.Docno

				pgDocDetailList = append(pgDocDetailList, tmpDocDetail)
			}
		}

		err = pst.CreateInBatch(pgDocList, len(pgDocList))

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		err = pst.CreateInBatch(pgDocDetailList, len(pgDocDetailList))

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		return nil
	})

}

func StartJournalComsumeUpdated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := MQ_TOPIC_JOURNAL_UPDATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume journal created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Journal Created : %v", msg)
		doc := vfgl.JournalDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		pgDocList := []vfgl.JournalPg{}
		pgDocDetailList := []vfgl.JournalDetailPg{}

		//for _, doc := range docList
		{

			tmpJsonDoc, err := json.Marshal(doc)

			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}

			tmpDoc := vfgl.JournalPg{}
			err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}

			pgDocList = append(pgDocList, tmpDoc)

			docDetailList := doc.AccountBook

			for _, docDetail := range docDetailList {
				tmpDocDetail := vfgl.JournalDetailPg{}

				tmpJsonDocDetail, err := json.Marshal(docDetail)
				if err != nil {
					ms.Logger.Errorf(moduleName, err.Error())
				}
				err = json.Unmarshal([]byte(tmpJsonDocDetail), &tmpDocDetail)

				if err != nil {
					ms.Logger.Errorf(moduleName, err.Error())
				}

				tmpDocDetail.ParID = doc.ParID
				tmpDocDetail.ShopID = doc.ShopID
				tmpDocDetail.Docno = doc.Docno

				pgDocDetailList = append(pgDocDetailList, tmpDocDetail)
			}
		}

		err = pst.CreateInBatch(pgDocList, len(pgDocList))

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		err = pst.CreateInBatch(pgDocDetailList, len(pgDocDetailList))

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		return nil
	})

}

func StartJournalComsumeDeleted(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := MQ_TOPIC_JOURNAL_DELETED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume journal created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Journal Created : %v", msg)
		doc := vfgl.JournalDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		pgDocList := []vfgl.JournalPg{}
		pgDocDetailList := []vfgl.JournalDetailPg{}

		//for _, doc := range docList
		{

			tmpJsonDoc, err := json.Marshal(doc)

			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}

			tmpDoc := vfgl.JournalPg{}
			err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}

			// err = json.Unmarshal([]byte(msg), &doc)

			// if err != nil {
			// 	ms.Logger.Errorf(moduleName, err.Error())
			// }

			pgDocList = append(pgDocList, tmpDoc)

			docDetailList := doc.AccountBook

			for _, docDetail := range docDetailList {
				tmpDocDetail := vfgl.JournalDetailPg{}

				tmpJsonDocDetail, err := json.Marshal(docDetail)
				if err != nil {
					ms.Logger.Errorf(moduleName, err.Error())
				}
				err = json.Unmarshal([]byte(tmpJsonDocDetail), &tmpDocDetail)

				if err != nil {
					ms.Logger.Errorf(moduleName, err.Error())
				}

				tmpDocDetail.ParID = doc.ParID
				tmpDocDetail.ShopID = doc.ShopID
				tmpDocDetail.Docno = doc.Docno

				pgDocDetailList = append(pgDocDetailList, tmpDocDetail)
			}
		}

		err = pst.CreateInBatch(pgDocList, len(pgDocList))

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		err = pst.CreateInBatch(pgDocDetailList, len(pgDocDetailList))

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		return nil
	})

}

func StartJournalComsumeBlukCreated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := MQ_TOPIC_JOURNAL_BULK_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume journal created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Journal Created : %v", msg)
		docList := []vfgl.JournalDoc{}
		err := json.Unmarshal([]byte(msg), &docList)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		pgDocList := []vfgl.JournalPg{}
		pgDocDetailList := []vfgl.JournalDetailPg{}

		for _, doc := range docList {

			tmpJsonDoc, err := json.Marshal(doc)

			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}

			tmpDoc := vfgl.JournalPg{}
			err = json.Unmarshal([]byte(tmpJsonDoc), &tmpDoc)
			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}

			err = json.Unmarshal([]byte(msg), &docList)

			if err != nil {
				ms.Logger.Errorf(moduleName, err.Error())
			}

			pgDocList = append(pgDocList, tmpDoc)

			docDetailList := doc.AccountBook

			for _, docDetail := range docDetailList {
				tmpDocDetail := vfgl.JournalDetailPg{}

				tmpJsonDocDetail, err := json.Marshal(docDetail)
				if err != nil {
					ms.Logger.Errorf(moduleName, err.Error())
				}
				err = json.Unmarshal([]byte(tmpJsonDocDetail), &tmpDocDetail)

				if err != nil {
					ms.Logger.Errorf(moduleName, err.Error())
				}

				tmpDocDetail.ParID = doc.ParID
				tmpDocDetail.ShopID = doc.ShopID
				tmpDocDetail.Docno = doc.Docno

				pgDocDetailList = append(pgDocDetailList, tmpDocDetail)
			}
		}

		err = pst.CreateInBatch(pgDocList, len(pgDocList))

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		err = pst.CreateInBatch(pgDocDetailList, len(pgDocDetailList))

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		return nil
	})

}
