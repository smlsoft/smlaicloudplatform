package services

import (
	"fmt"
	"smlcloudplatform/pkg/vfgl/journal/repositories"

	"github.com/go-redis/redis/v8"
)

type IJournalWebsocketService interface {
	Pub(shopID string, processID string, screen string, message interface{}) error
	Sub(shopID string, processID string, screen string) (<-chan *redis.Message, string, error)
	UnSub(subID string) error

	SaveLastMessage(shopID string, processID string, screen string, message string) error
	GetLastMessage(shopID string, processID string, screen string) (string, error)

	SetWebsocket(shopID string, processID string, screen string, socketID string) error
	DelWebsocket(shopID string, processID string, screen string, socketID string) error
}

type JournalWebsocketService struct {
	cacheChannelName   string
	cacheMessageName   string
	cacheWebsocketName string
	repo               repositories.IJournalCacheRepository
}

func NewJournalWebsocketService(repo repositories.IJournalCacheRepository) *JournalWebsocketService {

	return &JournalWebsocketService{
		cacheChannelName:   "ws",
		cacheMessageName:   "wsmsg",
		cacheWebsocketName: "wssc",
		repo:               repo,
	}
}

func (svc JournalWebsocketService) Pub(shopID string, processID string, screen string, message interface{}) error {
	return svc.repo.Pub(shopID, processID, svc.cacheChannelName, screen, message)
}

func (svc JournalWebsocketService) Sub(shopID string, processID string, screen string) (<-chan *redis.Message, string, error) {
	return svc.repo.Sub(shopID, processID, svc.cacheChannelName, screen)
}

func (svc JournalWebsocketService) UnSub(subID string) error {
	return svc.repo.Unsub(subID)
}

func (svc JournalWebsocketService) SaveLastMessage(shopID string, processID string, screen string, message string) error {

	keyVal := fmt.Sprintf(":%s", screen)
	data := map[string]interface{}{
		keyVal: message,
	}

	return svc.repo.HSet(shopID, processID, svc.cacheMessageName, data)
}

func (svc JournalWebsocketService) GetLastMessage(shopID string, processID string, screen string) (string, error) {

	keyVal := fmt.Sprintf(":%s", screen)

	return svc.repo.HGet(shopID, processID, svc.cacheMessageName, keyVal)
}

func (svc JournalWebsocketService) ClearLastMessage(shopID string, processID string) error {

	return svc.repo.Del(shopID, processID, svc.cacheMessageName)
}

func (svc JournalWebsocketService) SetWebsocket(shopID string, processID string, screen string, socketID string) error {

	keyVal := fmt.Sprintf("%s:%s", socketID, screen)
	data := map[string]interface{}{
		keyVal: 1,
	}

	return svc.repo.HSet(shopID, processID, svc.cacheWebsocketName, data)
}

func (svc JournalWebsocketService) DelWebsocket(shopID string, processID string, screen string, socketID string) error {

	keyVal := fmt.Sprintf("%s:%s", socketID, screen)

	return svc.repo.HDel(shopID, processID, svc.cacheWebsocketName, keyVal)
}
