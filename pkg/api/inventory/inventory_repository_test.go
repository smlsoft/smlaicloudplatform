package inventory_test

import (
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/inventory"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repoMock inventory.InventoryRepository

func init() {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repoMock = inventory.NewInventoryRepository(mongoPersister)
}

func TestFindByID(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}

	idx, _ := primitive.ObjectIDFromHex("62398ea81e4743ecba54da23")

	doc, err := repoMock.FindByID(idx)

	if err != nil {
		t.Error(err)
	}

	t.Log(doc)

	if doc.DeletedAt.IsZero() {
		t.Log("is zero")
	}

}

func TestFindByItemGuid(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	doc, err := repoMock.FindByItemGuid("27daMDw274R5hHejrjkHDuI91ag", "ix001x")

	if err != nil {
		t.Error(err)
	}

	t.Log(doc)

	// if doc.DeletedAt.IsZero() {
	// 	t.Log("is zero")
	// }

}
