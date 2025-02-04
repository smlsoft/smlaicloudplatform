package stocktransaction_test

import (
	pkgModels "smlaicloudplatform/internal/models"
	stockprocessmodels "smlaicloudplatform/internal/stockprocess/models"
	stocktransactionmodels "smlaicloudplatform/internal/transaction/models"
	stocktransaction "smlaicloudplatform/internal/transaction/transactionconsumer/stocktransaction"
	"testing"

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
