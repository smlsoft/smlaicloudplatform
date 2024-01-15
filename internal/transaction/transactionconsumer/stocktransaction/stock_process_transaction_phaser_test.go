package stocktransaction_test

import (
	"testing"
	pkgModels "vfapi/internal/models"
	stockprocessmodels "vfapi/internal/stockprocess/models"
	stocktransactionmodels "vfapi/internal/transaction/models"
	stocktransaction "vfapi/internal/transaction/transactionconsumer/stocktransaction"

	"github.com/tj/assert"
)

func TestPhaserStockProcessFromStockData(t *testing.T) {

	giveStockTransactionModels := stocktransactionmodels.StockTransaction{
		ShopIdentity: pkgModels.ShopIdentity{
			ShopID: "SHOP0001",
		},
		Details: &[]stocktransactionmodels.StockTransactionDetail{
			{
				Barcode: "BARCODE0001",
			},
		},
	}

	phaser := stocktransaction.StockProcessTransactionPhaser{}

	err, gotStockProcessRequest := phaser.PhaseStockTransactionProcess(giveStockTransactionModels)

	assert.Nil(t, err)
	wantStockTransaction := stockprocessmodels.StockProcessRequest{
		ShopID:  "SHOP0001",
		Barcode: "BARCODE0001",
	}

	assert.Equal(t, wantStockTransaction, (*gotStockProcessRequest)[0], "StockProcessRequest")

}
