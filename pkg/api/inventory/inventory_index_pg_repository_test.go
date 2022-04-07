package inventory_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/models"
	"testing"
)

func newInventoryPgRepo() inventory.InventoryIndexPGRepository {
	persisterConfig := mock.NewPersisterPostgresqlConfig()
	pst := microservice.NewPersister(persisterConfig)
	repo := inventory.NewInventoryIndexPGRepository(pst)
	return repo
}

func TestCreate(t *testing.T) {
	repo := newInventoryPgRepo()

	idx := models.InventoryIndex{}
	idx.ID = "134567"
	idx.ShopID = "shopidx001"
	idx.GuidFixed = "fixguid"
	err := repo.Create(idx)

	if err != nil {
		t.Error(err)
	}
}

func TestCount(t *testing.T) {
	repo := newInventoryPgRepo()

	count, err := repo.Count("shopidx001", "fixguid")

	if err != nil {
		t.Error(err)
	}

	t.Log(count)
}

func TestFindByGuid(t *testing.T) {
	repo := newInventoryPgRepo()
	inv, err := repo.FindByGuid("shopidx001", "fixguid")

	if err != nil {
		t.Error(err)
	}

	t.Log(inv)
}

func TestDelete(t *testing.T) {
	repo := newInventoryPgRepo()

	err := repo.Delete("shopidx001", "fixguid")

	if err != nil {
		t.Error(err)
	}
}
