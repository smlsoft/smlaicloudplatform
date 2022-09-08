package repositories_test

import (
	"fmt"
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/product/inventory/repositories"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repoMock repositories.InventoryRepository

func init() {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repoMock = repositories.NewInventoryRepository(mongoPersister)
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
	doc, err := repoMock.FindByGuid("27dcEdktOoaSBYFmnN6G6ett4Jb", "2EQsi6PRQ3lAmXORqYD9zJltnsz")

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%v", doc)

	// if doc.DeletedAt.IsZero() {
	// 	t.Log("is zero")
	// }

}
