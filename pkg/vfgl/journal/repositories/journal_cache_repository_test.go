package repositories_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
	"testing"
)

const MockShopID = "TESTSHOP"
const prefixName = "ws"

var cache *microservice.Cacher
var repoCacheMock repositories.JournalCacheRepository

func init() {
	cacheConfig := mock.NewCacherConfig()
	cache = microservice.NewCacher(cacheConfig)
	repoCacheMock = *repositories.NewJournalCacheRepository(cache)
	// repoMock = category.NewCategoryRepository(mongoPersister)
}

func TestMSet(t *testing.T) {
	data := map[string]interface{}{
		"ws:s1": "mdata1",
		"ws:s2": "mdata2",
	}
	repoCacheMock.HSet("sh1", "p1", prefixName, data)
}

func TestMGet(t *testing.T) {

	result, err := repoCacheMock.HGet("sh1", "p1", prefixName, "m1")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}

func TestHGetAll(t *testing.T) {

	result, err := repoCacheMock.HFields("sh1", "p1", prefixName, "ws:*")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}
