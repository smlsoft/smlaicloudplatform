package stocktransaction

import (
	stockprocessmodels "smlaicloudplatform/internal/stockprocess/models"
	stocktransactionmodels "smlaicloudplatform/internal/transaction/models"
)

type IStockProcessTransactionPhaser interface {
	PhaseStockTransactionProcess(stocktransactionmodels.StockTransaction) (error, *[]stockprocessmodels.StockProcessRequest)
}

type StockProcessTransactionPhaser struct{}

func (s StockProcessTransactionPhaser) PhaseStockTransactionProcess(data stocktransactionmodels.StockTransaction) (error, *[]stockprocessmodels.StockProcessRequest) {

	details := make([]stockprocessmodels.StockProcessRequest, len(*data.Details))

	for i, detail := range *data.Details {
		details[i] = stockprocessmodels.StockProcessRequest{
			ShopID:  data.ShopID,
			Barcode: detail.Barcode,
		}
	}
	return nil, &details
}
