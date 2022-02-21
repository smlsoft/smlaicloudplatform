package merchantservice_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/merchantservice"
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

var memberService merchantservice.IMemberService

func setup() {
	pst := microservice.NewPersisterMongo(&TestPersisterMongoConfig{})

	memberService = merchantservice.NewMemberService(pst)
}

func TestMerchantMemberSave(t *testing.T) {
	setup()

	err := memberService.Save("mx1", "ux3", models.ROLE_OWNER)

	if err != nil {
		t.Error(err)
	}
}

func TestMerchantMemberFindRole(t *testing.T) {
	setup()

	role, err := memberService.FindRole("mx1", "ux1")

	if err != nil {
		t.Error(err)
	}

	t.Log(role)
}

func TestMerchantMemberFindByMerchant(t *testing.T) {
	setup()

	members, err := memberService.FindByMerchantId("mx1")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(members)
}

func TestMerchantMemberFindByUsername(t *testing.T) {
	setup()

	members, err := memberService.FindByUsername("ux1")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(members)
}
