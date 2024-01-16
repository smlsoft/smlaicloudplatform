package stockprocess

import (
	"smlcloudplatform/internal/logger"
	productBarcodeRepositories "smlcloudplatform/internal/product/productbarcode/repositories"
	stockModel "smlcloudplatform/internal/stockprocess/models"
	"smlcloudplatform/internal/stockprocess/repositories"
	"smlcloudplatform/pkg/stockcalculator"
)

var log logger.ILogger

type IStockCalculator interface {
	// CalculateStockPrice(stockData []StockData) float64
	CalculatorStock(shopID string, barcode string) error
	GetStockDataList(shopID string, barcode string) ([]stockModel.StockData, error)
	WriteUpdateStockDataChanged(stockData []stockModel.StockData) error
}

type StockCalculator struct {
	stockMovementRepo  repositories.IStockProcessPGRepository
	productBarcodeRepo productBarcodeRepositories.IProductBarcodePGRepository
}

func NewStockCalculator(
	repo repositories.IStockProcessPGRepository,
	productBarcodeRepo productBarcodeRepositories.IProductBarcodePGRepository,
) IStockCalculator {

	return &StockCalculator{
		stockMovementRepo:  repo,
		productBarcodeRepo: productBarcodeRepo,
	}
}

func (sc *StockCalculator) CalculatorStock(shopID string, barcode string) error {

	stockDataList, err := sc.GetStockDataList(shopID, barcode)
	if err != nil {
		return err
	}

	productBarcode, err := sc.productBarcodeRepo.Get(shopID, barcode)
	if err != nil {
		return err
	}

	var stockDataChangeLists []stockModel.StockData

	if len(stockDataList) > 0 {
		calculator := stockcalculator.NewStockCalculator(shopID, barcode, 2, 0, 0)
		for i, data := range stockDataList {
			if data.CalcFlag == 1 {

				var qty, amountExcludeVat float64
				if data.HasCostFromOtherDoc() {

					var costOut float64

					for _, costData := range stockDataList {
						if costData.DocNo == data.DocRef {
							costOut = costData.AverageCost
							break
						}
					}

					qty = data.CalcQty
					amountExcludeVat = costOut * qty
				} else if data.TransFlag == 66 {
					qty = data.CalcQty
					amountExcludeVat = calculator.AverageCost()
				} else {
					qty = data.CalcQty
					amountExcludeVat = data.SumAmountExcludeVat
				}

				cost, averageCost := calculator.ApplyStock(qty, amountExcludeVat)

				isDataChange := data.SumOfCost != cost || data.AverageCost != averageCost
				if isDataChange {
					stockDataList[i].SumOfCost = cost
					stockDataList[i].AverageCost = averageCost
					stockDataChangeLists = append(stockDataChangeLists, stockDataList[i])
				}
				logger.GetLogger().Debugf("ApplyStock: %+v", data)
			} else {

				var cost, averageCost float64
				if data.HasCostFromOtherDoc() {

					var costOut float64

					for _, costData := range stockDataList {
						if costData.DocNo == data.DocRef {
							costOut = costData.AverageCost
							break
						}
					}
					// find cost from other doc
					cost, averageCost = calculator.ReduceStockWithCost(data.CalcQty, costOut)
				} else {
					cost, averageCost = calculator.ReduceStock(data.CalcQty)
				}

				isDataChange := data.SumOfCost != cost || data.AverageCost != averageCost
				if isDataChange {
					stockDataList[i].SumOfCost = cost
					stockDataList[i].AverageCost = averageCost
					stockDataChangeLists = append(stockDataChangeLists, stockDataList[i])
				}
				logger.GetLogger().Debugf("ReduceStock: %+v", data)
			}

			logger.GetLogger().Debugf("After Place Movement: %+v", calculator)
		}

		// write update movement
		if len(stockDataChangeLists) > 0 {
			err = sc.WriteUpdateStockDataChanged(stockDataChangeLists)
			if err != nil {
				return err
			}
		}

		isProductBalanceNotEqual := calculator.BalanceQty() != productBarcode.BalanceQty || calculator.BalanceAmount() != productBarcode.BalanceAmount || calculator.AverageCost() != productBarcode.AverageCost

		if isProductBalanceNotEqual {

			productBarcode.BalanceQty = calculator.BalanceQty()
			productBarcode.BalanceAmount = calculator.BalanceAmount()
			productBarcode.AverageCost = calculator.AverageCost()

			err = sc.productBarcodeRepo.Update(shopID, barcode, productBarcode)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (sc *StockCalculator) GetStockDataList(shopID string, barcode string) ([]stockModel.StockData, error) {
	return sc.stockMovementRepo.GetStockTransactionList(shopID, barcode)
}

func (sc *StockCalculator) WriteUpdateStockDataChanged(stockData []stockModel.StockData) error {
	return sc.stockMovementRepo.UpdateStockTransactionChange(stockData)
}
