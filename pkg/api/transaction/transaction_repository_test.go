package transaction_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/api/transaction"
	"smlcloudplatform/pkg/models"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	give := models.Transaction{
		MerchantId: "mx01",
		GuidFixed:  "fx01",
		Items: []models.TransactionDetail{
			{
				InventoryId:  "inv01",
				ItemSku:      "sku01",
				CategoryGuid: "xxx",
			},
		},
	}

	notWant := ""

	idx, err := repo.Create(give)

	if err != nil {
		t.Error(err)
	}

	if idx == notWant {
		t.Error("create failed")
	}

}

func TestFindTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	merchantId := "mx01"
	give := "fx01"

	want := "fx01"
	trans, err := repo.FindByGuid(merchantId, give)

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
	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	merchantId := "mx01"
	// give := "fx01"

	want := 1
	trans, _, err := repo.FindPage(merchantId, "", 1, 20)

	if err != nil {
		t.Error(err)
		return
	}

	if len(trans) < want {
		t.Error("find failed")
		return
	}
}
