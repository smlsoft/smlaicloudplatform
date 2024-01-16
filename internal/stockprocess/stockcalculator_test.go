package stockprocess_test

import (
	"smlcloudplatform/internal/stockprocess"
	stockModel "smlcloudplatform/internal/stockprocess/models"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockStockProcessPGRepository struct {
	mock.Mock
}

func (m *MockStockProcessPGRepository) GetStockTransactionList(shopID string, barcode string) ([]stockModel.StockData, error) {
	ret := m.Called(shopID, barcode)
	return ret.Get(0).([]stockModel.StockData), ret.Error(1)
}

func (m *MockStockProcessPGRepository) UpdateStockTransactionChange(stockData []stockModel.StockData) error {
	ret := m.Called(stockData)
	return ret.Error(0)
}

func (m *MockStockProcessPGRepository) ExecuteUpdateProductBarcodeStockBalance(shopID string, barcode string) error {
	ret := m.Called(shopID, barcode)
	return ret.Error(0)
}

// var repo repositories.IStockProcessPGRepository

// func init() {
// 	cfg := config.NewConfig()
// 	persister := microservice.NewPersister(cfg.PersisterConfig())
// 	repo = repositories.NewStockProcessPGRepository(persister)
// }

func TestStockProcess(t *testing.T) {

	// stockLists, err := repo.GetStockTransactionList("2IZS0jFeRXWPidSupyXN7zQIlaS", "888555")
	// assert.Nil(t, err)
	// assert.NotNil(t, stockLists)
	// assert.Equal(t, 2, len(stockLists))

	var stockDataLists []stockModel.StockData
	stockDataLists = append(stockDataLists, stockModel.StockData{
		ShopID:              "SHOPID",
		Barcode:             "BARCODE",
		DocNo:               "DOC1",
		TransFlag:           12,
		CalcFlag:            1,
		CalcQty:             100,
		StandValue:          1,
		DivideValue:         1,
		SumAmount:           1000,
		SumAmountExcludeVat: 1000,
	})

	stockDataLists = append(stockDataLists, stockModel.StockData{
		ShopID:              "SHOPID",
		Barcode:             "BARCODE",
		DocNo:               "DOC2",
		TransFlag:           12,
		CalcFlag:            1,
		CalcQty:             2,
		StandValue:          1,
		DivideValue:         1,
		SumAmount:           22,
		SumAmountExcludeVat: 22,
	})

	stockDataLists = append(stockDataLists, stockModel.StockData{
		ShopID:              "SHOPID",
		Barcode:             "BARCODE",
		DocNo:               "DOC2",
		TransFlag:           12,
		CalcFlag:            1,
		CalcQty:             3,
		StandValue:          1,
		DivideValue:         1,
		SumAmount:           30,
		SumAmountExcludeVat: 30,
	})

	repo := new(MockStockProcessPGRepository)
	repo.On("GetStockTransactionList", "SHOPID", "BARCODE").Return(stockDataLists, nil)

	repo.On("UpdateStockTransactionChange", mock.Anything).Return(nil)

	process := stockprocess.NewStockCalculator(repo, nil)
	process.CalculatorStock("SHOPID", "BARCODE")

}
