package inventory_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/models"
	"testing"
)

func getInventoryOptionRepo() inventory.InventoryOptionRepository {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := inventory.NewInventoryOptionRepository(mongoPersister)
	return repo
}

func TestCreateInventoryOption(t *testing.T) {
	repo := getInventoryOptionRepo()

	give := models.InventoryOption{
		GuidFixed:     "fx01",
		ShopID:        "mx01",
		InventoryID:   "inv01",
		OptionGroupID: "opts01",
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
