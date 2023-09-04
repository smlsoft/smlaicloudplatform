package services_test

import (
	"smlcloudplatform/pkg/stockbalanceimport/models"
	"smlcloudplatform/pkg/stockbalanceimport/repositories"
	"smlcloudplatform/pkg/stockbalanceimport/services"
	stockbalance_models "smlcloudplatform/pkg/transaction/stockbalance/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewStockBalanceImportService(t *testing.T) {
	cacheRepo := repositories.MockStockBalanceImportCacheRepository{}

	cacheExpire := 60 * time.Minute
	cacheRepo.On("CreateMeta", "shoptest", "t1000", models.StockBalanceImportMeta{
		TaskID:    "t1000",
		TotalItem: 1000,
		Parts: []models.StockBalanceImportPartMeta{
			{
				PartID:     "t1000-1",
				PartNumber: 1,
				Status:     0,
			},
			{
				PartID:     "t1000-2",
				PartNumber: 2,
				Status:     0,
			},
		},
	}, cacheExpire).Return(nil)

	cacheRepo.On("CreatePart", "shoptest", "t1000-1", models.StockBalanceImportPartCache{
		TaskID: "t1000",
		StockBalanceImportPartMeta: models.StockBalanceImportPartMeta{
			PartID:     "t1000-1",
			PartNumber: 1,
			Status:     0,
		},
		Detail: []stockbalance_models.StockBalanceDetail{},
	}, cacheExpire).Return(nil)
	cacheRepo.On("CreatePart", "shoptest", "t1000-2", models.StockBalanceImportPartCache{
		TaskID: "t1000",
		StockBalanceImportPartMeta: models.StockBalanceImportPartMeta{
			PartID:     "t1000-2",
			PartNumber: 2,
			Status:     0,
		},
		Detail: []stockbalance_models.StockBalanceDetail{},
	}, cacheExpire).Return(nil)

	cacheRepo.On("CreateMeta", "shoptest", "123", models.StockBalanceImportMeta{
		TaskID:    "123",
		TotalItem: 800,
		Parts: []models.StockBalanceImportPartMeta{
			{
				PartID:     "123-1",
				PartNumber: 1,
				Status:     0,
			},
			{
				PartID:     "123-2",
				PartNumber: 2,
				Status:     0,
			},
		},
	}, cacheExpire).Return(nil)

	svc := services.NewStockBalanceImportService(
		&cacheRepo,
		nil,
		func(int) string {
			return "123"
		},
	)
	t.Run("Create Task 1000", func(t *testing.T) {

		got := models.StockBalanceImportTaskRequest{
			TotalItem: 1000,
		}

		want := models.StockBalanceImportTask{
			TaskID:    "123",
			TotalItem: 1000,
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

	t.Run("Create Task 800", func(t *testing.T) {

		got := models.StockBalanceImportTaskRequest{
			TotalItem: 800,
		}

		want := models.StockBalanceImportTask{
			TaskID:    "123",
			TotalItem: 800,
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
