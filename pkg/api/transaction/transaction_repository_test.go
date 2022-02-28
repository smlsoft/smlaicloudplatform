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
		return
	}

	if idx == notWant {
		t.Error("create failed")
		return
	}

}

func TestUpdateTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	trans := models.Transaction{
		MerchantId: "mx01",
		GuidFixed:  "fx02",
		Items: []models.TransactionDetail{
			{
				InventoryId:  "inv01",
				ItemSku:      "sku01",
				CategoryGuid: "xxx",
			},
		},
	}

	give := models.Transaction{
		MerchantId: "mx01",
		GuidFixed:  "fx02",
		Items: []models.TransactionDetail{
			{
				InventoryId:  "inv01",
				ItemSku:      "sku01",
				CategoryGuid: "xxx",
			},
			{
				InventoryId:  "inv02",
				ItemSku:      "sku02",
				CategoryGuid: "xxx2",
			},
		},
	}

	notWant := ""

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

	transFind, err := repo.FindByGuid(give.GuidFixed, give.MerchantId)

	if err != nil {
		t.Error(err)
		return
	}

	if len(transFind.Items) < 2 {
		t.Error("Update failed")
	}

}

func TestDeleteTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	give := models.Transaction{
		MerchantId: "mx01",
		GuidFixed:  "fx03",
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
		return
	}

	if idx == notWant {
		t.Error("create failed")
		return
	}

	err = repo.Delete(give.GuidFixed, give.MerchantId)

	if err != nil {
		t.Error(err)
		return
	}

	transFind, err := repo.FindByGuid(give.GuidFixed, give.MerchantId)

	if err != nil {
		t.Error(err)
		return
	}

	if transFind.GuidFixed != "" {
		t.Error("Delete failed")
	}

}

func TestFindTransaction(t *testing.T) {
	mongoPersisterConfig := mock.NewPersisterMongo()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := transaction.NewTransactionRepository(mongoPersister)

	merchantId := "mx01"
	give := "fx01"

	want := "fx01"
	trans, err := repo.FindByGuid(give, merchantId)

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
