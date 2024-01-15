package stockprocess_test

import (
	"testing"
	productbarcoderepository "vfapi/internal/product/productbarcode/repositories"
	"vfapi/internal/stockprocess"
	"vfapi/internal/stockprocess/repositories"
	"vfapi/pkg/config"
	"vfapi/pkg/microservice"
)

var repo repositories.IStockProcessPGRepository
var productBarcodePGRepository productbarcoderepository.IProductBarcodePGRepository

func init() {
	cfg := config.NewConfig()
	persister := microservice.NewPersister(cfg.PersisterConfig())
	repo = repositories.NewStockProcessPGRepository(persister)
	productBarcodePGRepository = productbarcoderepository.NewProductBarcodePGRepository(persister)
}

func TestStockProcessRealDBTest(t *testing.T) {

	// stockLists, err := repo.GetStockTransactionList("2IZS0jFeRXWPidSupyXN7zQIlaS", "888555")
	// assert.Nil(t, err)
	// assert.NotNil(t, stockLists)
	// assert.Equal(t, 2, len(stockLists))

	process := stockprocess.NewStockCalculator(repo, productBarcodePGRepository)
	process.CalculatorStock("2IZS0jFeRXWPidSupyXN7zQIlaS", "888555")
}
