package repositories

import (
	"smlcloudplatform/pkg/stockbalanceimport/models"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockStockBalanceImportCacheRepository struct {
	mock.Mock
}

func (m *MockStockBalanceImportCacheRepository) CreateMeta(shopID, cacheKey string, value models.StockBalanceImportMeta, expire time.Duration) error {
	args := m.Called(shopID, cacheKey, value, expire)
	return args.Error(0)
}

func (m *MockStockBalanceImportCacheRepository) UpdateMeta(shopID, cacheKey string, value models.StockBalanceImportMeta) error {
	args := m.Called(shopID, cacheKey, value)
	return args.Error(0)
}

func (m *MockStockBalanceImportCacheRepository) GetMeta(shopID, cacheKey string) (models.StockBalanceImportMeta, error) {
	args := m.Called(shopID, cacheKey)
	return args.Get(0).(models.StockBalanceImportMeta), args.Error(1)
}

func (m *MockStockBalanceImportCacheRepository) CreatePart(shopID, cacheKey string, value models.StockBalanceImportPartCache, expire time.Duration) error {
	args := m.Called(shopID, cacheKey, value, expire)
	return args.Error(0)
}

func (m *MockStockBalanceImportCacheRepository) UpdatePart(shopID, cacheKey string, value models.StockBalanceImportPartCache) error {
	args := m.Called(shopID, cacheKey, value)
	return args.Error(0)
}

func (m *MockStockBalanceImportCacheRepository) GetPart(shopID, cacheKey string) (models.StockBalanceImportPartCache, error) {
	args := m.Called(shopID, cacheKey)
	return args.Get(0).(models.StockBalanceImportPartCache), args.Error(1)
}
