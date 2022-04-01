package transaction_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/transaction"
	"smlcloudplatform/pkg/models"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	give := models.TransactionDoc{}

	give.ShopID = "mx01"
	give.GuidFixed = "fx01"
	give.Items = &[]models.TransactionDetail{
		{
			ItemSku:      "sku01",
			CategoryGuid: "xxx",
		},
	}

	notWant := primitive.NilObjectID

	idx, err := repo.Create(give)

	if err != nil {
		t.Error(err)
		return
	}

	if idx == notWant {
		t.Error("create failed")
		return
	}

}

func TestUpdateTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	trans := models.TransactionDoc{}

	trans.ShopID = "mx01"
	trans.GuidFixed = "fx02"
	trans.Items = &[]models.TransactionDetail{
		{
			InventoryID:  "inv01",
			ItemSku:      "sku01",
			CategoryGuid: "xxx",
		},
	}

	give := models.TransactionDoc{}

	give.ShopID = "mx01"
	give.GuidFixed = "fx02"
	give.Items = &[]models.TransactionDetail{
		{
			InventoryID:  "inv01",
			ItemSku:      "sku01",
			CategoryGuid: "xxx",
		},
		{
			InventoryID:  "inv02",
			ItemSku:      "sku02",
			CategoryGuid: "xxx2",
		},
	}

	notWant := primitive.NilObjectID

	idx, err := repo.Create(trans)

	if err != nil {
		t.Error(err)
		return
	}

	if idx == notWant {
		t.Error("create failed")
		return
	}

	err = repo.Update(give.GuidFixed, give)

	if err != nil {
		t.Error(err)
		return
	}

	transFind, err := repo.FindByGuid(give.GuidFixed, give.ShopID)

	if err != nil {
		t.Error(err)
		return
	}

	if len(*transFind.Items) < 2 {
		t.Error("Update failed")
	}

}

func TestDeleteTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	give := models.TransactionDoc{}

	give.ShopID = "mx01"
	give.GuidFixed = "fx03"
	give.Items = &[]models.TransactionDetail{
		{
			InventoryID:  "inv01",
			ItemSku:      "sku01",
			CategoryGuid: "xxx",
		},
	}

	notWant := primitive.NilObjectID

	idx, err := repo.Create(give)

	if err != nil {
		t.Error(err)
		return
	}

	if idx == notWant {
		t.Error("create failed")
		return
	}

	err = repo.Delete(give.GuidFixed, give.ShopID, "test")

	if err != nil {
		t.Error(err)
		return
	}

	transFind, err := repo.FindByGuid(give.GuidFixed, give.ShopID)

	if err != nil {
		t.Error(err)
		return
	}

	if transFind.GuidFixed != "" {
		t.Error("Delete failed")
	}

}

func TestFindTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	shopID := "mx01"
	give := "fx01"

	want := "fx01"
	trans, err := repo.FindByGuid(give, shopID)

	if err != nil {
		t.Error(err)
		return
	}

	if trans.GuidFixed != want {
		t.Error("find failed")
		return
	}
}

func TestFindPageTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	shopID := "mx01"
	// give := "fx01"

	want := 1
	trans, _, err := repo.FindPage(shopID, "", 1, 20)

	if err != nil {
		t.Error(err)
		return
	}

	if len(trans) < want {
		t.Error("find failed")
		return
	}
}
