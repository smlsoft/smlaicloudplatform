package journal

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/vfgl/journal/config"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
	"smlcloudplatform/pkg/vfgl/journal/services"
	"time"
)

func MigrationJournalTable(ms *microservice.Microservice, cfg microservice.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.JournalPg{},
		models.JournalDetailPg{},
	)
	return nil
}

func StartJournalComsumeCreated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := config.MQ_TOPIC_JOURNAL_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume journal created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Journal Created : %v", msg)
		doc := models.JournalDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		repo := repositories.NewJournalPgRepository(pst)
		svc := services.NewJournalConsumeService(repo)

		_, err = svc.UpSert(doc.ShopID, doc.DocNo, doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}
		return nil
	})

}

func StartJournalComsumeUpdated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := config.MQ_TOPIC_JOURNAL_UPDATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume journal update"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Journal update : %v", msg)
		doc := models.JournalDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		repo := repositories.NewJournalPgRepository(pst)
		svc := services.NewJournalConsumeService(repo)

		_, err = svc.UpSert(doc.ShopID, doc.DocNo, doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		return nil
	})

}

func StartJournalComsumeDeleted(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := config.MQ_TOPIC_JOURNAL_DELETED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume journal created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Journal Created : %v", msg)
		doc := models.JournalDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		repo := repositories.NewJournalPgRepository(pst)
		svc := services.NewJournalConsumeService(repo)

		err = svc.Delete(doc.ShopID, doc.DocNo)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}
		return nil
	})

}

func StartJournalComsumeBlukCreated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := config.MQ_TOPIC_JOURNAL_BULK_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume journal created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Journal Created : %v", msg)
		docList := []models.JournalDoc{}
		err := json.Unmarshal([]byte(msg), &docList)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		repo := repositories.NewJournalPgRepository(pst)
		svc := services.NewJournalConsumeService(repo)

		err = svc.SaveInBatch(docList)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}
		return nil
	})

}
