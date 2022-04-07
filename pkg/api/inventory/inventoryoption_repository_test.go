package inventory_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/models"
	"testing"
)

func getInventoryOptionMainRepo() inventory.InventoryOptionMainRepository {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := inventory.NewInventoryOptionMainRepository(mongoPersister)
	return repo
}

func TestCreateInventoryOptionMain(t *testing.T) {
	repo := getInventoryOptionMainRepo()

	give := models.InventoryOptionMain{
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
