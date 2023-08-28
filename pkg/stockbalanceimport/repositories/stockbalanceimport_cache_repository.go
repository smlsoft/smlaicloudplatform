package repositories

import (
	"encoding/json"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/stockbalanceimport/models"
	"time"
)

type IStockBalanceImportCacheRepository interface {
	CreateMeta(shopID, cacheKey string, value models.StockBalanceImportMeta, expire time.Duration) error
	UpdateMeta(shopID, cacheKey string, value models.StockBalanceImportMeta) error
	GetMeta(shopID, cacheKey string) (models.StockBalanceImportMeta, error)
	CreatePart(shopID, cacheKey string, value models.StockBalanceImportPartCache, expire time.Duration) error
	UpdatePart(shopID, cacheKey string, value models.StockBalanceImportPartCache) error
	GetPart(shopID, cacheKey string) (models.StockBalanceImportPartCache, error)
}

type StockBalanceImportCacheRepository struct {
	prefixCacheKey string
	cache          microservice.ICacher
}

func NewStockBalanceImportCacheRepository(cache microservice.ICacher) StockBalanceImportCacheRepository {
	return StockBalanceImportCacheRepository{
		prefixCacheKey: "stockbalanceimport",
		cache:          cache,
	}
}

func (repo StockBalanceImportCacheRepository) CreateMeta(shopID, cacheKey string, value models.StockBalanceImportMeta, expire time.Duration) error {
	tempCacheKey := repo.generateCacheKey("meta", shopID, cacheKey)
	return repo.cache.SetNX(tempCacheKey, value, expire)
}

func (repo StockBalanceImportCacheRepository) UpdateMeta(shopID, cacheKey string, value models.StockBalanceImportMeta) error {
	tempCacheKey := repo.generateCacheKey("meta", shopID, cacheKey)
	return repo.cache.SetXX(tempCacheKey, value, time.Duration(-1))
}

func (repo StockBalanceImportCacheRepository) GetMeta(shopID, cacheKey string) (models.StockBalanceImportMeta, error) {
	tempCacheKey := repo.generateCacheKey("meta", shopID, cacheKey)
	rawResult, err := repo.cache.Get(tempCacheKey)

	if err != nil {
		return models.StockBalanceImportMeta{}, err
	}

	if rawResult == "" {
		return models.StockBalanceImportMeta{}, nil
	}

	result := models.StockBalanceImportMeta{}

	err = json.Unmarshal([]byte(rawResult), &result)

	if err != nil {
		return models.StockBalanceImportMeta{}, err
	}

	return result, nil
}

func (repo StockBalanceImportCacheRepository) CreatePart(shopID, cacheKey string, value models.StockBalanceImportPartCache, expire time.Duration) error {
	tempCacheKey := repo.generateCacheKey("part", shopID, cacheKey)
	return repo.cache.SetNX(tempCacheKey, value, expire)
}

func (repo StockBalanceImportCacheRepository) UpdatePart(shopID, cacheKey string, value models.StockBalanceImportPartCache) error {
	tempCacheKey := repo.generateCacheKey("part", shopID, cacheKey)
	return repo.cache.SetXX(tempCacheKey, value, time.Duration(-1))
}

func (repo StockBalanceImportCacheRepository) GetPart(shopID, cacheKey string) (models.StockBalanceImportPartCache, error) {
	tempCacheKey := repo.generateCacheKey("part", shopID, cacheKey)
	rawResult, err := repo.cache.Get(tempCacheKey)

	if err != nil {
		return models.StockBalanceImportPartCache{}, err
	}

	if rawResult == "" {
		return models.StockBalanceImportPartCache{}, nil
	}

	result := models.StockBalanceImportPartCache{}

	err = json.Unmarshal([]byte(rawResult), &result)

	if err != nil {
		return models.StockBalanceImportPartCache{}, err
	}

	return result, nil
}

func (repo StockBalanceImportCacheRepository) generateCacheKey(cacheType, shopID, cacheKey string) string {
	return fmt.Sprintf("%s:%s%s:%s", repo.prefixCacheKey, shopID, cacheKey, cacheType)
}
