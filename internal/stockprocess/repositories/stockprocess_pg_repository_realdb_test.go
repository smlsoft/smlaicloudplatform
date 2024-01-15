package repositories_test

import (
	"testing"
	"vfapi/internal/stockprocess/repositories"
	"vfapi/pkg/config"
	"vfapi/pkg/microservice"

	"github.com/stretchr/testify/assert"
)

var repo repositories.IStockProcessPGRepository

func init() {
	cfg := config.NewConfig()
	persister := microservice.NewPersister(cfg.PersisterConfig())
	repo = repositories.NewStockProcessPGRepository(persister)
}

func TestGetStockProcessList(t *testing.T) {

	stockLists, err := repo.GetStockTransactionList("2IZS0jFeRXWPidSupyXN7zQIlaS", "888555")
	assert.Nil(t, err)
	assert.NotNil(t, stockLists)
	assert.Equal(t, 2, len(stockLists))
}

func TestUpdateStockTransaction(t *testing.T) {
	stockLists, err := repo.GetStockTransactionList("2IZS0jFeRXWPidSupyXN7zQIlaS", "888555")
	assert.Nil(t, err)

	for i, _ := range stockLists {
		stockLists[i].AverageCost = float64(1)
		stockLists[i].SumOfCost = float64(1)
	}

	err = repo.UpdateStockTransactionChange(stockLists)
	assert.Nil(t, err)
}
