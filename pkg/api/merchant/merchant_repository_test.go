package merchant_test

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/merchant"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"testing"
	"time"
)

func TestFindMerchant(t *testing.T) {

	// os.Setenv("MONGODB_URI", "mongodb://root:rootx@localhost:27017/")
	// defer os.Unsetenv("MONGODB_URI")

	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repository := merchant.NewMerchantRepository(mongoPersister)

	newGuidFixed := utils.NewGUID()
	createAt := time.Now()

	give := models.Merchant{
		GuidFixed: newGuidFixed,
		Name1:     "merchant test",
		CreatedBy: "test",
		CreatedAt: createAt,
	}

	want := &models.Merchant{
		GuidFixed: newGuidFixed,
		Name1:     "merchant test",
		CreatedBy: "test",
		CreatedAt: createAt,
	}

	get, err := repository.Create(give)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if get == "" {
		t.Error(errors.New("Create merchant Failed"))
	}

	t.Logf("Create merchant Success With ID %v", get)

	getUser, err := repository.FindByGuid(want.GuidFixed)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if getUser.GuidFixed != want.GuidFixed {
		t.Error(errors.New("Create merchant And Find Not Match"))
		return
	}

}
