package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type IJournalWebsocketService interface {
	PubDoc(shopID string, processID string, screen string, message interface{}) error
	SubDoc(shopID string, processID string, screen string) (<-chan *redis.Message, string, error)

	UnSub(subID string) error

	SaveLastMessage(shopID string, processID string, screen string, message string) error
	GetLastMessage(shopID string, processID string, screen string) (string, error)
	ClearLastMessage(shopID string, processID string) error

	SetWebsocket(shopID string, processID string, screen string) error
	DelWebsocket(shopID string, processID string, screen string) error
	ExistsWebsocket(shopID string, processID string) (bool, error)
	ExpireWebsocket(shopID string, processID string) error

	DocRefPool(shopID string, username string, ws *websocket.Conn) error
	SetDocRefPool(shopID string, username string, docRef string) error
	ExistsDocRefPool(shopID string, docRef string) (bool, error)
	GetDocRefPool(shopID string, docRef string) (string, error)
	GetAllDocRefPool(shopID string) (map[string]string, error)
	DelDocRefPool(shopID string, username string, docRef string) error

	ExistsDocRefUserPool(shopID string, username string) (bool, error)
	GetDocRefUserPool(shopID string, username string) (string, error)

	DocRefSelect(shopID string, username string, docRef string) (bool, error)
	DocRefUnSelect(shopID string, username string) (bool, error)
	DocRefSelectForce(shopID string, username string, docRef string, forceSelect bool) (bool, error)
}

type JournalWebsocketService struct {
	cacheChannelDoc     string
	cacheChannelDocRef  string
	cachePoolDocRef     string
	cachePoolDocRefUser string
	cacheMessageName    string
	cacheWebsocketName  string
	cacheExpire         time.Duration
	repo                repositories.IJournalCacheRepository
}

func NewJournalWebsocketService(repo repositories.IJournalCacheRepository, cacheExpire time.Duration) *JournalWebsocketService {

	return &JournalWebsocketService{
		cacheChannelDoc:     "chdoc",
		cacheChannelDocRef:  "chdocref",
		cachePoolDocRef:     "wsdocref",
		cachePoolDocRefUser: "wsdocuser",
		cacheMessageName:    "wsmsg",
		cacheWebsocketName:  "wssc",
		repo:                repo,
		cacheExpire:         cacheExpire,
	}
}

func (svc JournalWebsocketService) PubDoc(shopID string, processID string, screen string, message interface{}) error {
	channel := svc.getChannelDoc(shopID, processID, screen, svc.cacheChannelDoc)
	return svc.repo.Pub(channel, message)
}

func (svc JournalWebsocketService) SubDoc(shopID string, processID string, screen string) (<-chan *redis.Message, string, error) {
	channel := svc.getChannelDoc(shopID, processID, screen, svc.cacheChannelDoc)
	return svc.repo.Sub(channel)
}

func (svc JournalWebsocketService) PubDocRef(shopID string, message interface{}) error {
	channel := svc.getChannelDocRef(shopID, svc.cacheChannelDocRef)
	return svc.repo.Pub(channel, message)
}

func (svc JournalWebsocketService) SubDocRef(shopID string) (<-chan *redis.Message, string, error) {
	channel := svc.getChannelDocRef(shopID, svc.cacheChannelDocRef)
	return svc.repo.Sub(channel)
}

func (svc JournalWebsocketService) UnSub(subID string) error {
	return svc.repo.Unsub(subID)
}

func (svc JournalWebsocketService) SaveLastMessage(shopID string, processID string, screen string, message string) error {

	keyVal := screen
	data := map[string]interface{}{
		keyVal: message,
	}
	cacheKeyName := svc.getTagID(shopID, processID, svc.cacheMessageName)
	return svc.repo.HSet(cacheKeyName, data)
}

func (svc JournalWebsocketService) GetLastMessage(shopID string, processID string, screen string) (string, error) {

	cacheKeyName := svc.getTagID(shopID, processID, svc.cacheMessageName)
	keyVal := screen
	return svc.repo.HGet(cacheKeyName, keyVal)
}

func (svc JournalWebsocketService) ClearLastMessage(shopID string, processID string) error {
	cacheKeyName := svc.getTagID(shopID, processID, svc.cacheMessageName)
	return svc.repo.Del(cacheKeyName)
}

func (svc JournalWebsocketService) SetWebsocket(shopID string, processID string, screen string) error {
	cacheKeyName := svc.getTagID(shopID, processID, svc.cacheWebsocketName)

	keyVal := screen
	data := map[string]interface{}{
		keyVal: 1,
	}

	return svc.repo.HSet(cacheKeyName, data)
}

func (svc JournalWebsocketService) DelWebsocket(shopID string, processID string, screen string) error {
	cacheKeyName := svc.getTagID(shopID, processID, svc.cacheWebsocketName)
	keyVal := screen
	return svc.repo.HDel(cacheKeyName, keyVal)
}

func (svc JournalWebsocketService) ExistsWebsocket(shopID string, processID string) (bool, error) {
	cacheKeyName := svc.getTagID(shopID, processID, svc.cacheWebsocketName)
	return svc.repo.Exists(cacheKeyName)
}

func (svc JournalWebsocketService) ExpireWebsocket(shopID string, processID string) error {
	cacheKeyName := svc.getTagID(shopID, processID, svc.cacheWebsocketName)
	return svc.repo.Expire(cacheKeyName, svc.cacheExpire)
}

// doc ref
func (svc JournalWebsocketService) SetDocRefPool(shopID string, username string, docRef string) error {

	if len(docRef) < 1 {
		return errors.New("doc ref is empty")
	}

	cacheKeyDocRef := svc.getTagID(shopID, "", svc.cachePoolDocRef)
	cacheKeyUser := svc.getTagID(shopID, "", svc.cachePoolDocRefUser)

	isDocRefSelected, err := svc.repo.HExists(cacheKeyDocRef, docRef)
	if err != nil {
		return err
	}

	if isDocRefSelected {
		return errors.New("doc ref is selected")
	}

	isDocUserSelected, err := svc.repo.HExists(cacheKeyUser, username)
	if err != nil {
		return err
	}

	if isDocUserSelected {
		return errors.New("user is selected")
	}

	dataUser := map[string]interface{}{
		username: docRef,
	}

	svc.repo.HSet(cacheKeyUser, dataUser)

	if err != nil {
		return err
	}

	dataDocRef := map[string]interface{}{
		docRef: username,
	}

	return svc.repo.HSet(cacheKeyDocRef, dataDocRef)
}

func (svc JournalWebsocketService) ExistsDocRefPool(shopID string, docRef string) (bool, error) {
	cacheKeyName := svc.getTagID(shopID, "", svc.cachePoolDocRef)
	return svc.repo.HExists(cacheKeyName, docRef)
}

func (svc JournalWebsocketService) GetDocRefPool(shopID string, docRef string) (string, error) {
	cacheKeyName := svc.getTagID(shopID, "", svc.cachePoolDocRef)
	return svc.repo.HGet(cacheKeyName, docRef)
}

func (svc JournalWebsocketService) GetAllDocRefPool(shopID string) (map[string]string, error) {
	cacheKeyName := svc.getTagID(shopID, "", svc.cachePoolDocRef)
	return svc.repo.HGetAll(cacheKeyName)
}

func (svc JournalWebsocketService) DelDocRefPool(shopID string, username string, docRef string) error {
	cacheKeyDocRef := svc.getTagID(shopID, "", svc.cachePoolDocRef)
	cacheKeyUser := svc.getTagID(shopID, "", svc.cachePoolDocRefUser)

	err := svc.repo.HDel(cacheKeyUser, username)

	if err != nil {
		return err
	}

	return svc.repo.HDel(cacheKeyDocRef, docRef)
}

// user
func (svc JournalWebsocketService) ExistsDocRefUserPool(shopID string, username string) (bool, error) {
	cacheKeyName := svc.getTagID(shopID, "", svc.cachePoolDocRefUser)
	return svc.repo.HExists(cacheKeyName, username)
}

func (svc JournalWebsocketService) GetDocRefUserPool(shopID string, username string) (string, error) {
	cacheKeyName := svc.getTagID(shopID, "", svc.cachePoolDocRefUser)
	return svc.repo.HGet(cacheKeyName, username)
}

func (svc JournalWebsocketService) getChannelDocRef(shopID string, prefix string) string {
	tempChannel := fmt.Sprintf("%s-%s", prefix, shopID)
	return tempChannel
}

func (svc JournalWebsocketService) getChannelDoc(shopID string, processID string, prefix string, screen string) string {
	tempID := svc.getTagID(shopID, processID, prefix)
	return fmt.Sprintf("%s:%s", tempID, screen)
}

func (JournalWebsocketService) getTagID(shopID string, processID string, prefix string) string {
	// tempID := utils.FastHash(fmt.Sprintf("%s%s", shopID, processID))
	tempID := fmt.Sprintf("%s-%s%s", prefix, shopID, processID)
	return tempID
}

func (svc JournalWebsocketService) DocRefUnSelect(shopID string, username string) (bool, error) {

	docRef, err := svc.GetDocRefUserPool(shopID, username)
	if err != nil {
		return false, err
	}

	if docRef == "" {
		return false, errors.New("user is not selected")
	}

	err = svc.DelDocRefPool(shopID, username, docRef)

	if err != nil {
		return false, err
	}

	// send websocket to user
	svc.ClearLastMessage(shopID, username)

	tempDocRef, _ := json.Marshal(models.JournalRef{
		DocRef: "",
	})

	svc.PubDoc(shopID, username, "form", tempDocRef)

	docRefEvent := models.DocRefEvent{
		DocRef:   docRef,
		Username: username,
		Status:   "unselected",
	}

	err = svc.pubDocRefSelect(shopID, docRefEvent)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (svc JournalWebsocketService) DocRefSelectForce(shopID string, username string, docRef string, forceSelect bool) (bool, error) {
	if forceSelect {
		_, err := svc.DocRefUnSelect(shopID, username)
		if err != nil {
			return false, err
		}
	}

	return svc.DocRefSelect(shopID, username, docRef)

}

func (svc JournalWebsocketService) DocRefSelect(shopID string, username string, docRef string) (bool, error) {
	isExists, err := svc.ExistsDocRefPool(shopID, docRef)

	if err != nil {
		return false, err
	}

	if isExists {
		return false, errors.New("user is selected")
	}

	err = svc.SetDocRefPool(shopID, username, docRef)

	if err != nil {
		return false, err
	}

	// send websocket to user
	tempDocRef, _ := json.Marshal(models.JournalRef{
		DocRef: docRef,
	})
	svc.PubDoc(shopID, username, "form", tempDocRef)
	err = svc.SaveLastMessage(shopID, username, "form", string(tempDocRef))
	if err != nil {
		return false, err
	}

	docRefEvent := models.DocRefEvent{
		DocRef:   docRef,
		Username: username,
		Status:   "selected",
	}

	err = svc.pubDocRefSelect(shopID, docRefEvent)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (svc JournalWebsocketService) pubDocRefSelect(shopID string, docRefEvent models.DocRefEvent) error {
	tempData, err := json.Marshal(docRefEvent)
	if err != nil {
		return err
	}

	err = svc.PubDocRef(shopID, tempData)

	if err != nil {
		return err
	}

	return nil
}

func (svc JournalWebsocketService) DocRefPool(shopID string, username string, ws *websocket.Conn) error {

	cacheMsg, subID, err := svc.SubDocRef(shopID)

	if err != nil {
		return err
	}

	defer func(ws *websocket.Conn, svc JournalWebsocketService) {
		ws.Close()
		svc.UnSub(subID)
	}(ws, svc)

	for {
		temp := <-cacheMsg
		if temp != nil {

			err = ws.WriteMessage(websocket.TextMessage, []byte(temp.Payload))
			if err != nil {
				return err
			}
		}
	}
}
