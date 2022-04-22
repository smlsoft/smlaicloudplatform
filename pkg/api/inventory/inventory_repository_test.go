package inventory_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/inventory"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getInventoryRepo() inventory.InventoryRepository {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := inventory.NewInventoryRepository(mongoPersister)
	return repo
}

func TestFindByID(t *testing.T) {

	repo := getInventoryRepo()
	idx, _ := primitive.ObjectIDFromHex("62398ea81e4743ecba54da23")

	doc, err := repo.FindByID(idx)

	if err != nil {
		t.Error(err)
	}

	t.Log(doc)

	if doc.DeletedAt.IsZero() {
		t.Log("is zero")
	}

}

func TestFindByItemGuid(t *testing.T) {

	repo := getInventoryRepo()

	doc, err := repo.FindByItemGuid("27daMDw274R5hHejrjkHDuI91ag", "ix001x")

	if err != nil {
		t.Error(err)
	}

	t.Log(doc)

	// if doc.DeletedAt.IsZero() {
	// 	t.Log("is zero")
	// }

}
