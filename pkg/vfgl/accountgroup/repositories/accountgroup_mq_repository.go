package repositories

import (
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/vfgl/accountgroup/config"
	"smlcloudplatform/pkg/vfgl/accountgroup/models"
)

type IAccountGroupMqRepository interface {
	Create(doc models.AccountGroupDoc) error
	Update(doc models.AccountGroupDoc) error
	Delete(doc common.Identity) error
	CreateInBatch(docList []models.AccountGroupDoc) error
}

type AccountGroupMqRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.AccountGroupDoc]
}

func NewAccountGroupMqRepository(prod microservice.IProducer) AccountGroupMqRepository {
	mqKey := ""

	insRepo := AccountGroupMqRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.AccountGroupDoc](prod, config.AccountGroupMessageQueueConfig{}, "")
	return insRepo
}
