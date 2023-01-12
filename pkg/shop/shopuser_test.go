package shop_test

import (
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shop"
	"smlcloudplatform/pkg/shop/models"
	"testing"
)

type TestPersisterMongoConfig struct{}

func (TestPersisterMongoConfig) MongodbURI() string {
	return "mongodb://root:rootx@localhost:27017/"
}

func (TestPersisterMongoConfig) DB() string {
	return "micro_test"
}

var shopUserRepo shop.IShopUserRepository

func setup() {
	pst := microservice.NewPersisterMongo(&TestPersisterMongoConfig{})

	shopUserRepo = shop.NewShopUserRepository(pst)
}

func TestShopMemberSave(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	setup()

	err := shopUserRepo.Save("mx1", "ux3", models.ROLE_OWNER)

	if err != nil {
		t.Error(err)
	}
}

func TestShopMemberFindRole(t *testing.T) {
	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	setup()

	role, err := shopUserRepo.FindRole("mx1", "ux3")

	if err != nil {
		t.Error(err)
	}

	t.Log(role)
}

func TestShopMemberFindByShop(t *testing.T) {
	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	setup()

	members, err := shopUserRepo.FindByShopID("mx1")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(members)
}

func TestShopMemberFindByUsername(t *testing.T) {
	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	setup()

	members, err := shopUserRepo.FindByUsername("ux1")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(members)
}
