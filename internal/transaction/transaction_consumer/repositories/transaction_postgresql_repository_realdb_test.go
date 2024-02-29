package repositories_test

import (
	"os"
	"smlcloudplatform/internal/config"
	"smlcloudplatform/internal/transaction/models"
	"smlcloudplatform/internal/transaction/transaction_consumer/repositories"
	"smlcloudplatform/pkg/microservice"
	"strconv"
	"testing"
	"time"

	commonModel "smlcloudplatform/internal/models"

	"github.com/tj/assert"
)

var stockTransaction = models.StockTransaction{}
var repo repositories.ITransactionPGRepository

func init() {
	os.Setenv("MODE", "test")
	cfg := config.NewConfig()

	pst := microservice.NewPersister(cfg.PersisterConfig())
	repo = repositories.NewTransactionPGRepository(pst)

	stockTransaction = models.StockTransaction{
		ShopIdentity: commonModel.ShopIdentity{
			ShopID: "shoptester",
		},
		DocNo: "TRXTEST",
		Details: &[]models.StockTransactionDetail{
			models.StockTransactionDetail{
				DocNo:   "TRXTEST",
				ShopID:  "shoptester",
				Barcode: "BAR1",
			},
			models.StockTransactionDetail{
				DocNo:   "TRXTEST",
				ShopID:  "shoptester",
				Barcode: "BAR2",
			},
		},
	}
}

func TestCreateProductBarcodeInRealDB(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}

	err := repo.Create(stockTransaction)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}

	bar, err := repo.Get(stockTransaction.ShopID, stockTransaction.DocNo)
	assert.NoError(t, err)

	assert.Equal(t, stockTransaction.ShopID, bar.ShopID)
	assert.Equal(t, stockTransaction.DocNo, bar.DocNo)

}

func TestUpdate(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}

	currentTime := time.Now()
	timeStr := currentTime.Format("20060201150405")
	(*stockTransaction.Details)[0].CostPerUnit, _ = strconv.ParseFloat(timeStr, 64)

	err := repo.Update(stockTransaction.ShopID, stockTransaction.DocNo, stockTransaction)
	assert.NoError(t, err)

	bar, err := repo.Get(stockTransaction.ShopID, stockTransaction.DocNo)
	assert.NoError(t, err)

	assert.Equal(t, (*stockTransaction.Details)[0].CostPerUnit, (*bar.Details)[0].CostPerUnit)
}

func TestDeleteInRealDB(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}

	err := repo.Delete(stockTransaction.ShopID, stockTransaction.DocNo)
	assert.NoError(t, err)
}
