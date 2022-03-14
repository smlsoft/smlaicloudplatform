package inventory_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/models"
	"testing"
)

func TestCreateInventoryOption(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := inventory.NewInventoryOptionRepository(mongoPersister)

	give := models.InventoryOption{
		GuidFixed:     "fx01",
		ShopId:        "mx01",
		InventoryId:   "inv01",
		OptionGroupId: "opts01",
	}

	_, err := repo.Create(give)

	if err != nil {
		t.Error(err)
		return
	}

	findDoc, err := repo.FindByGuid("fx01", "mx01")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(findDoc)
}
