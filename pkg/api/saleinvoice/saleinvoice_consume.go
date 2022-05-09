package saleinvoice

import (
	"encoding/json"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"time"
)

const MQ_CONSUMER_SALEINVOICE_CREATED_GROUPNAME string = "saleinvoice-postgresql-oncreated-consumer"
const MQ_CONSUMER_SALEINVOICE_UPDATED_GROUPNAME string = "saleinvoice-postgresql-onupdated-consumer"
const MQ_CONSUMER_SALEINVOICE_DELETED_GROUPNAME string = "saleinvoice-postgresql-ondeleted-consumer"

func StartSaleinvoiceComsumeCreated(ms *microservice.Microservice, cfg microservice.IConfig) {
	topic := MQ_TOPIC_TRANSACTION_CREATED
	groupID := MQ_CONSUMER_SALEINVOICE_CREATED_GROUPNAME
	timeout := time.Duration(-1)

	mqConfig := cfg.MQConfig()

	mq := microservice.NewMQ(mqConfig, ms.Logger)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*7)

	//Consume saleinvoice Created
	ms.Consume(mqConfig.URI(), topic, groupID, timeout, func(ctx microservice.IContext) error {
		moduleName := "comsume saleinvoice created"

		msg := ctx.ReadInput()

		saleinvoiceJson := models.Saleinvoice{}
		err := json.Unmarshal([]byte(msg), &saleinvoiceJson)

		if err != nil {
			ms.Log(moduleName, err.Error())
		}

		// postgresql
		// saleInvoiceTable := SaleInvoiceTable{}
		// err = json.Unmarshal([]byte(msg), &saleInvoiceTable)

		// if err != nil {
		// 	ms.Log(moduleName, err.Error())
		// }

		// pst := ms.Persister(cfg.PersisterConfig())
		// err = pst.Create(&saleInvoiceTable)

		// if err != nil {
		// 	ms.Log(moduleName, err.Error())
		// }

		// elk
		// elkPst := ms.ElkPersister(cfg.ElkPersisterConfig())
		// err = elkPst.CreateWithID(trans.GuidFixed, &trans)

		// if err != nil {
		// 	ms.Log(moduleName, err.Error())
		// }
		return nil
	})

}

func StartSaleInvoiceConsumerUpdated(ms *microservice.Microservice, cfg microservice.IConfig) error {

	return nil
}
