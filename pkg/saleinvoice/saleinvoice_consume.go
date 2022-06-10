package saleinvoice

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/saleinvoice/config"
	"smlcloudplatform/pkg/saleinvoice/models"
	"smlcloudplatform/pkg/saleinvoice/repositories"
	"smlcloudplatform/pkg/saleinvoice/services"
	"time"
)

func MigrationSaleinvoiceTable(ms *microservice.Microservice, cfg microservice.IConfig) error {
	pst := ms.Persister(cfg.PersisterConfig())
	pst.AutoMigrate(
		models.SaleinvoicePg{},
		models.SaleinvoiceDetailPg{},
	)
	return nil
}

func StartSaleinvoiceComsumeCreated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := config.MQ_TOPIC_SALEINVOICE_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume saleinvoice created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Saleinvoice Created : %v", msg)
		doc := models.SaleinvoiceDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		repo := repositories.NewSaleinvoicePgRepository(pst)
		svc := services.NewSaleinvoiceConsumeService(repo)

		_, err = svc.UpSert(doc.ShopID, doc.DocNo, doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}
		return nil
	})

}

func StartSaleinvoiceComsumeUpdated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := config.MQ_TOPIC_SALEINVOICE_UPDATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume saleinvoice update"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Saleinvoice update : %v", msg)
		doc := models.SaleinvoiceDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		repo := repositories.NewSaleinvoicePgRepository(pst)
		svc := services.NewSaleinvoiceConsumeService(repo)

		_, err = svc.UpSert(doc.ShopID, doc.DocNo, doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		return nil
	})

}

func StartSaleinvoiceComsumeDeleted(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := config.MQ_TOPIC_SALEINVOICE_DELETED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume saleinvoice created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Saleinvoice Created : %v", msg)
		doc := models.SaleinvoiceDoc{}
		err := json.Unmarshal([]byte(msg), &doc)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		repo := repositories.NewSaleinvoicePgRepository(pst)
		svc := services.NewSaleinvoiceConsumeService(repo)

		err = svc.Delete(doc.ShopID, doc.DocNo)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}
		return nil
	})

}

func StartSaleinvoiceComsumeBlukCreated(ms *microservice.Microservice, cfg microservice.IConfig, groupID string) {

	topicCreated := config.MQ_TOPIC_SALEINVOICE_BULK_CREATED
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)

	mq.CreateTopicR(topicCreated, 5, 1, time.Hour*24*7)

	ms.Consume(mqConfig.URI(), topicCreated, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume saleinvoice created"

		pst := ms.Persister(cfg.PersisterConfig())
		msg := ctx.ReadInput()

		ms.Logger.Debugf("Consume Saleinvoice Created : %v", msg)
		docList := []models.SaleinvoiceDoc{}
		err := json.Unmarshal([]byte(msg), &docList)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}

		repo := repositories.NewSaleinvoicePgRepository(pst)
		svc := services.NewSaleinvoiceConsumeService(repo)

		err = svc.SaveInBatch(docList)

		if err != nil {
			ms.Logger.Errorf(moduleName, err.Error())
		}
		return nil
	})

}
