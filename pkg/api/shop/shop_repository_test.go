package shop_test

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/shop"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"testing"
	"time"
)

func TestFindShop(t *testing.T) {

	// os.Setenv("MONGODB_URI", "mongodb://root:rootx@localhost:27017/")
	// defer os.Unsetenv("MONGODB_URI")

	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repository := shop.NewShopRepository(mongoPersister)

	newGuidFixed := utils.NewGUID()
	createAt := time.Now()

	give := models.ShopDoc{
		Shop: models.Shop{
			GuidFixed: newGuidFixed,
			Name1:     "shop test",
		},
		Activity: models.Activity{
			CreatedBy: "test",
			CreatedAt: createAt,
		},
	}

	want := &models.ShopDoc{
		Shop: models.Shop{
			GuidFixed: newGuidFixed,
			Name1:     "shop test",
		},
		Activity: models.Activity{
			CreatedBy: "test",
			CreatedAt: createAt,
		},
	}

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
