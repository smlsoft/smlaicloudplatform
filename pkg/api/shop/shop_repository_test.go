package shop_test

import (
	"errors"
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/shop"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"testing"
	"time"
)

func TestFindShop(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	// os.Setenv("MONGODB_URI", "mongodb://root:rootx@localhost:27017/")
	// defer os.Unsetenv("MONGODB_URI")

	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repository := shop.NewShopRepository(mongoPersister)

	newGuidFixed := utils.NewGUID()
	createAt := time.Now()

	activity := models.ActivityDoc{
		CreatedBy: "test",
		CreatedAt: createAt,
	}

	give := models.ShopDoc{}
	give.GuidFixed = newGuidFixed
	give.Name1 = "shop test"
	give.ActivityDoc = activity

	want := models.ShopDoc{}
	want.GuidFixed = newGuidFixed
	want.Name1 = "shop test"
	want.ActivityDoc = activity

	get, err := repository.Create(give)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if get == "" {
		t.Error(errors.New("Create shop Failed"))
	}

	t.Logf("Create shop Success With ID %v", get)

	getUser, err := repository.FindByGuid(want.GuidFixed)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if getUser.GuidFixed != want.GuidFixed {
		t.Error(errors.New("Create shop And Find Not Match"))
		return
	}

}
