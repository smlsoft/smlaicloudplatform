package repositories_test

import (
	"os"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/mock"
	inventoryModel "smlcloudplatform/pkg/product/inventory/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/repositories"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateSaleinvoice(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := repositories.NewSaleinvoiceRepository(mongoPersister)

	inv := inventoryModel.Inventory{
		ItemSku:      "sku01",
		CategoryGuid: "xxx",
	}

	give := models.SaleinvoiceDoc{}

	give.ShopID = "mx01"
	give.GuidFixed = "fx01"
	give.Items = &[]models.SaleinvoiceDetail{
		{
			InventoryInfo: inventoryModel.InventoryInfo{
				Inventory: inv,
			},
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

func TestUpdateSaleinvoice(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := repositories.NewSaleinvoiceRepository(mongoPersister)

	invInfo := inventoryModel.InventoryInfo{
		Inventory: inventoryModel.Inventory{
			ItemSku:      "sku01",
			CategoryGuid: "xxx",
		},
	}

	trans := models.SaleinvoiceDoc{}

	trans.ShopID = "mx01"
	trans.GuidFixed = "fx02"
	trans.Items = &[]models.SaleinvoiceDetail{
		{
			InventoryInfo: invInfo,
		},
	}

	invGive1 := inventoryModel.InventoryInfo{
		Inventory: inventoryModel.Inventory{
			ItemSku:      "sku01",
			CategoryGuid: "xxx",
		},
	}

	invGive2 := inventoryModel.InventoryInfo{
		Inventory: inventoryModel.Inventory{
			ItemSku:      "sku01",
			CategoryGuid: "xxx",
		},
	}

	give := models.SaleinvoiceDoc{}

	give.ShopID = "mx01"
	give.GuidFixed = "fx02"
	give.Items = &[]models.SaleinvoiceDetail{
		{
			InventoryInfo: invGive1,
		},
		{
			InventoryInfo: invGive2,
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

	err = repo.Update(give.ShopID, give.GuidFixed, give)

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

func TestDeleteSaleinvoice(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}

	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := repositories.NewSaleinvoiceRepository(mongoPersister)

	give := models.SaleinvoiceDoc{}

	invGive1 := inventoryModel.InventoryInfo{
		Inventory: inventoryModel.Inventory{
			ItemSku:      "sku01",
			CategoryGuid: "xxx",
		},
	}

	give.ShopID = "mx01"
	give.GuidFixed = "fx03"
	give.Items = &[]models.SaleinvoiceDetail{
		{
			InventoryInfo: invGive1,
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

func TestFindSaleinvoice(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := repositories.NewSaleinvoiceRepository(mongoPersister)

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

func TestFindPageSaleinvoice(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	mongoPersisterConfig := mock.NewPersisterMongoConfig()
	mongoPersister := microservice.NewPersisterMongo(mongoPersisterConfig)
	repo := repositories.NewSaleinvoiceRepository(mongoPersister)

	shopID := "mx01"
	// give := "fx01"

	pageable := micromodels.Pageable{
		Query: "",
		Page:  1,
		Limit: 20,
	}

	want := 1
	trans, _, err := repo.FindPage(shopID, pageable)

	if err != nil {
		t.Error(err)
		return
	}

	if len(trans) < want {
		t.Error("find failed")
		return
	}
}
