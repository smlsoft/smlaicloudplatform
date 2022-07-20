package repositories

import (
	"fmt"
	"smlcloudplatform/internal/microservice"

	"github.com/go-redis/redis/v8"
)

type IJournalCacheRepository interface {
	Pub(shopID string, processID string, prefix string, screen string, message interface{}) error
	Sub(shopID string, processID string, prefix string, screen string) (<-chan *redis.Message, string, error)
	Unsub(subID string) error

	HSet(shopID string, processID string, prefix string, data map[string]interface{}) error
	HGet(shopID string, processID string, prefix string, field string) (string, error)
	HFields(shopID string, processID string, prefix string, pattern string) ([]string, error)
	HDel(shopID string, processID string, prefix string, fields ...string) error
	Del(shopID string, processID string, prefix string) error
}

type JournalCacheRepository struct {
	cache microservice.ICacher
}

func NewJournalCacheRepository(cache microservice.ICacher) *JournalCacheRepository {
	return &JournalCacheRepository{cache: cache}
}

func (repo JournalCacheRepository) Pub(shopID string, processID string, prefix string, screen string, message interface{}) error {
	channelName := repo.getChannelName(shopID, processID, prefix, screen)
	return repo.cache.Pub(channelName, message)
}

func (repo JournalCacheRepository) Sub(shopID string, processID string, prefix string, screen string) (<-chan *redis.Message, string, error) {
	channelName := repo.getChannelName(shopID, processID, prefix, screen)
	return repo.cache.Sub(channelName)
}

func (repo JournalCacheRepository) Unsub(subID string) error {
	return repo.cache.Unsub(subID)
}

func (repo JournalCacheRepository) HSet(shopID string, processID string, prefix string, data map[string]interface{}) error {
	cacheKeyName := repo.getTagID(shopID, processID, prefix)
	return repo.cache.HMSet(cacheKeyName, data)
}

func (repo JournalCacheRepository) HGet(shopID string, processID string, prefix string, field string) (string, error) {
	cacheKeyName := repo.getTagID(shopID, processID, prefix)
	return repo.cache.HGet(cacheKeyName, field)
}

func (repo JournalCacheRepository) HFields(shopID string, processID string, prefix string, pattern string) ([]string, error) {
	cacheKeyName := repo.getTagID(shopID, processID, prefix)
	return repo.cache.HFields(cacheKeyName, pattern)
}

func (repo JournalCacheRepository) HDel(shopID string, processID string, prefix string, fields ...string) error {
	cacheKeyName := repo.getTagID(shopID, processID, prefix)
	return repo.cache.HDel(cacheKeyName, fields...)
}

func (repo JournalCacheRepository) Del(shopID string, processID string, prefix string) error {
	cacheKeyName := repo.getTagID(shopID, processID, prefix)
	return repo.cache.Del(cacheKeyName)
}

func (repo JournalCacheRepository) getChannelName(shopID string, processID string, prefix string, screen string) string {
	tempID := repo.getTagID(shopID, processID, prefix)
	return fmt.Sprintf("%s:%s", tempID, screen)
}

func (repo JournalCacheRepository) getTagID(shopID string, processID string, prefix string) string {
	// tempID := utils.FastHash(fmt.Sprintf("%s%s", shopID, processID))
	tempID := fmt.Sprintf("%s-%s%s", prefix, shopID, processID)
	return tempID
}
