package merchant_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/merchant"
	"smlcloudplatform/pkg/models"
	"testing"
)

type TestPersisterMongoConfig struct{}

func (TestPersisterMongoConfig) MongodbURI() string {
	return "mongodb://root:rootx@localhost:27017/"
}

func (TestPersisterMongoConfig) DB() string {
	return "micro_test"
}

var merchantUserRepo merchant.IMerchantUserRepository

func setup() {
	pst := microservice.NewPersisterMongo(&TestPersisterMongoConfig{})

	merchantUserRepo = merchant.NewMerchantUserRepository(pst)
}

func TestMerchantMemberSave(t *testing.T) {
	setup()

	err := merchantUserRepo.Save("mx1", "ux3", models.ROLE_OWNER)

	if err != nil {
		t.Error(err)
	}
}

func TestMerchantMemberFindRole(t *testing.T) {
	setup()

	role, err := merchantUserRepo.FindRole("mx1", "ux3")

	if err != nil {
		t.Error(err)
	}

	t.Log(role)
}

func TestMerchantMemberFindByMerchant(t *testing.T) {
	setup()

	members, err := merchantUserRepo.FindByMerchantId("mx1")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(members)
}

func TestMerchantMemberFindByUsername(t *testing.T) {
	setup()

	members, err := merchantUserRepo.FindByUsername("ux1")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(members)
}
