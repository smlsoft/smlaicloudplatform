package repositories_test

import (
	"context"
	"smlaicloudplatform/internal/vfgl/accountperiodmaster/repositories"
	"smlaicloudplatform/mock"
	"smlaicloudplatform/pkg/microservice"
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
