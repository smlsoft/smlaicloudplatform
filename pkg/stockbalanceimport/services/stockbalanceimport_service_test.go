package services_test

import (
	"smlcloudplatform/pkg/stockbalanceimport/models"
	"smlcloudplatform/pkg/stockbalanceimport/repositories"
	"smlcloudplatform/pkg/stockbalanceimport/services"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStockBalanceImportService(t *testing.T) {
	cacheRepo := repositories.MockStockBalanceImportCacheRepository{}

	cacheRepo.On("Save", "shoptest", "123", models.StockBalanceImportPartCache{}, 0).Return(nil)

	svc := services.NewStockBalanceImportService(
		&cacheRepo,
		nil,
		func(int) string {
			return "123"
		},
	)
	t.Run("Test Create Task", func(t *testing.T) {

		got := models.StockBalanceImportTaskRequest{
			TotalItem: 1000,
		}

		want := models.StockBalanceImportTask{
			TaskID: "123",
			Parts: []models.StockBalanceImportPart{
				{
					PartID:     "123-1",
					PartNumber: 1,
				},
				{
					PartID:     "123-2",
					PartNumber: 2,
				},
			},
		}

		result, err := svc.CreateTask("shoptest", got)

		require.NoError(t, err)
		require.Equal(t, want, result)

	})
}
