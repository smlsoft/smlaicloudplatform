package repositories

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"strings"
	"time"
)

type IMasterSyncCacheRepository interface {
	Save(shopID string) error
	SaveWithModule(shopID string, moduleName string) error
	Get(shopID string) (time.Time, error)
	GetWithModule(shopID string, moduleName string) (time.Time, error)
}

type MasterSyncCacheRepository struct {
	cache        microservice.ICacher
	moduleName   string
	allMasterKey string
}

func NewMasterSyncCacheRepository(cache microservice.ICacher, moduleName string) MasterSyncCacheRepository {
	return MasterSyncCacheRepository{
		cache:        cache,
		moduleName:   moduleName,
		allMasterKey: "all",
	}
}

func (repo MasterSyncCacheRepository) Save(shopID string) error {
	return repo.SaveWithModule(shopID, repo.moduleName)
}

func (repo MasterSyncCacheRepository) SaveWithModule(shopID string, moduleName string) error {
	changeTime := time.Now()
	cacheModuleKey := repo.getCacheModuleKeyWithModule(shopID, moduleName)
	cacheModuleAllMasterKey := repo.getCacheModuleKeyWithModule(shopID, repo.allMasterKey)

	repo.cache.SetNoExpire(cacheModuleAllMasterKey, changeTime)
	return repo.cache.SetNoExpire(cacheModuleKey, changeTime)
}

func (repo MasterSyncCacheRepository) Get(shopID string) (time.Time, error) {
	return repo.GetWithModule(shopID, repo.moduleName)
}

func (repo MasterSyncCacheRepository) GetWithModule(shopID string, moduleName string) (time.Time, error) {
	cacheModuleKey := repo.getCacheModuleKeyWithModule(shopID, moduleName)

	strTime, err := repo.cache.Get(cacheModuleKey)
	if err != nil {
		fmt.Println(err.Error())
		return time.Time{}, nil
	}

	strTime = strings.ReplaceAll(strTime, "\"", "")

	fmt.Println(strTime)
	valTime, err := time.Parse(time.RFC3339Nano, strTime)
	if err != nil {
		fmt.Println(err.Error())
		return time.Time{}, nil
	}

	return valTime, nil
}

func (repo MasterSyncCacheRepository) getCacheModuleKeyWithModule(shopID string, moduleName string) string {
	return fmt.Sprintf("mastersync-%s::%s", shopID, moduleName)
}
