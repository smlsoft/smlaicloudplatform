package repositories_test

import (
	"encoding/json"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/mastersync/repositories"
	"testing"
	"time"
)

const MockShopID = "TESTSHOP"

var cache *microservice.Cacher
var repoCacheMock repositories.MasterSyncCacheRepository

func init() {
	cacheConfig := mock.NewCacherConfig()
	cache = microservice.NewCacher(cacheConfig)
	repoCacheMock = repositories.NewMasterSyncCacheRepository(cache, "test")
	// repoMock = category.NewCategoryRepository(mongoPersister)
}

func TestSetCache(t *testing.T) {
	err := repoCacheMock.Save(MockShopID)

	if err != nil {
		t.Error(err)
	}
}
func TestGetCache(t *testing.T) {
	val, err := repoCacheMock.Get(MockShopID)

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(val)
}

func TestSetCacheWithModule(t *testing.T) {
	err := repoCacheMock.SaveWithModule(MockShopID, "XTEST")

	if err != nil {
		t.Error(err)
	}
}

func TestGetCacheWithModule(t *testing.T) {
	val, err := repoCacheMock.GetWithModule(MockShopID, "XTEST")

	if err != nil {
		t.Error(err)
	}

	t.Log(val)
}

func TestTimeStr(t *testing.T) {
	value := time.Now()
	fmt.Println("xxxx")

	tempVal := getType(value)

	tx, err := json.Marshal(tempVal)

	fmt.Println("-------")
	fmt.Println(string(tx))

	if err != nil {
		t.Error(err)
	}
}

func getType(v interface{}) interface{} {
	return v
}
