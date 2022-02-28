package merchant_test

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/merchant"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := merchant.NewMerchantUserRepository(mongoPersister)

	err := repo.Save("25H2pZ8v2jRVGwjOLKBAzSaHgOA", "dev01", "owner")

	if err != nil {
		t.Error(err.Error())
		return
	}

	memUser, err := repo.FindByMerchantIdAndUsername("25H2pZ8v2jRVGwjOLKBAzSaHgOA", "dev01")

	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Printf("%v \n", memUser)
	if memUser.Role == "" {
		t.Error("find error")
		return
	}

}
