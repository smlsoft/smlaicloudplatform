package repositories_test

import (
	"context"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/repositories"
	"testing"
	"time"
)

var repoMock repositories.AccountPeriodMasterRepository

func init() {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repoMock = *repositories.NewAccountPeriodMasterRepository(mongoPersister)
}

func TestFindByDateRange(t *testing.T) {
	repoMock.FindByDateRange(context.TODO(), "shopID", time.Now(), time.Now())
}
