package journal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
	"smlcloudplatform/pkg/vfgl/journal/services"

	"github.com/gorilla/websocket"
)

type JournalWs struct {
	ms           *microservice.Microservice
	cfg          microservice.IConfig
	svcWebsocket services.IJournalWebsocketService
}

func NewJournalWs(ms *microservice.Microservice, cfg microservice.IConfig) JournalWs {
	cache := ms.Cacher(cfg.CacherConfig())

	cacheRepo := repositories.NewJournalCacheRepository(cache)
	svcWebsocket := services.NewJournalWebsocketService(cacheRepo)

	return JournalWs{
		ms:           ms,
		cfg:          cfg,
		svcWebsocket: svcWebsocket,
	}
}

func (h JournalWs) RouteSetup() {

	h.ms.GET("/gl/journal/ws/image", h.WebsocketImage)
	h.ms.GET("/gl/journal/ws/form", h.WebsocketForm)

	h.ms.GET("/gl/journal/ws/docref", h.WebsocketDocRefPool)
	h.ms.GET("/gl/journal/img/select-all", h.GetAllDocRefPool)
	h.ms.POST("/gl/journal/img/select", h.SelectDocRefPool)
	h.ms.POST("/gl/journal/img/unselect", h.UnSelectDocRefPool)

	h.ms.GET("/checkx", h.Check)
}

func (h JournalWs) WebsocketImage(ctx microservice.IContext) error {

	screenName := "image"
	sendScreenName := "form"
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	processID := userInfo.Username

	if processID == "" {
		ctx.Response(http.StatusBadRequest, "processID parameter missing")
		return nil
	}

	socketID := utils.NewGUID()

	ws, err := h.ms.Websocket(socketID, ctx.ResponseWriter(), ctx.Request())

	if err != nil {
		return err
	}

	err = h.svcWebsocket.SetWebsocket(shopID, processID, screenName, socketID)

	if err != nil {
		return err
	}

	cacheMsg, subID, err := h.svcWebsocket.SubDoc(shopID, processID, screenName)

	if err != nil {
		return err
	}

	// Receive
	go func(ws *websocket.Conn) {
		defer func() {
			h.ms.WebsocketClose(socketID)
			h.svcWebsocket.UnSub(subID)
			h.svcWebsocket.DelWebsocket(shopID, processID, screenName, socketID)
			h.ClearDocRef(shopID, processID)
		}()

		for {

			journalRef := models.JournalRef{}
			err := ws.ReadJSON(&journalRef)

			if err != nil {
				return
			}

			tempRef, _ := json.Marshal(journalRef)
			h.svcWebsocket.PubDoc(shopID, processID, sendScreenName, tempRef)

			err = h.svcWebsocket.SaveLastMessage(shopID, processID, sendScreenName, string(tempRef))
			if err != nil {
				h.ms.Logger.Error(err.Error())
			}
		}
	}(ws)

	// Send
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

func (h JournalWs) WebsocketForm(ctx microservice.IContext) error {

	screenName := "form"
	sendScreenName := "image"
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	processID := userInfo.Username

	if processID == "" {
		ctx.Response(http.StatusBadRequest, "processID parameter missing")
		return nil
	}

	socketID := utils.NewID()

	ws, err := h.ms.Websocket(socketID, ctx.ResponseWriter(), ctx.Request())

	if err != nil {
		return err
	}

	err = h.svcWebsocket.SetWebsocket(shopID, processID, screenName, socketID)

	if err != nil {
		return err
	}

	cacheMsg, subID, err := h.svcWebsocket.SubDoc(shopID, processID, screenName)

	if err != nil {
		return err
	}

	// Receive
	go func(ws *websocket.Conn) {
		defer func() {
			ws.Close()
			h.ms.WebsocketClose(socketID)
			h.svcWebsocket.UnSub(subID)
			h.svcWebsocket.DelWebsocket(shopID, processID, screenName, socketID)
			h.ClearDocRef(shopID, processID)
		}()

		for {
			journalCommand := models.JournalCommand{}
			err := ws.ReadJSON(&journalCommand)

			if err != nil {
				return
			}

			tempRef, _ := json.Marshal(journalCommand)
			h.svcWebsocket.PubDoc(shopID, processID, sendScreenName, tempRef)

			switch journalCommand.Command {
			case "save":
				h.svcWebsocket.ClearLastMessage(shopID, processID)
				//clear
				h.ClearDocRef(shopID, processID)
			}
		}

	}(ws)

	// Send
	lastMessage, err := h.svcWebsocket.GetLastMessage(shopID, processID, screenName)
	if err != nil {
		h.ms.Logger.Error(err.Error())
	}

	if lastMessage != "" {
		err = ws.WriteMessage(websocket.TextMessage, []byte(lastMessage))

		if err != nil {
			h.ms.Logger.Error(err.Error())
		}
	}

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

func (h JournalWs) WebsocketDocRefPool(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	username := userInfo.Username

	if username == "" {
		ctx.Response(http.StatusBadRequest, "username is missing")
		return nil
	}

	socketID := utils.NewID()

	ws, err := h.ms.Websocket(socketID, ctx.ResponseWriter(), ctx.Request())

	if err != nil {
		return err
	}

	return h.svcWebsocket.DocRefPool(shopID, username, ws)

}

func (h JournalWs) ClearDocRef(shopID string, processID string) error {

	isExists, err := h.svcWebsocket.ExistsWebsocket(shopID, processID)
	if err != nil {
		h.ms.Logger.Error(err.Error())
	}

	if isExists {
		return nil
	}

	lastMessage, err := h.svcWebsocket.GetLastMessage(shopID, processID, "form")
	if err != nil {
		h.ms.Logger.Error(err.Error())
	}

	if lastMessage != "" {
		journalRef := models.JournalRef{}
		json.Unmarshal([]byte(lastMessage), &journalRef)

		err = h.svcWebsocket.DelDocRefPool(shopID, journalRef.DocRef)
		if err != nil {
			h.ms.Logger.Error(err.Error())
		}

		fmt.Println("Clear docref :: " + journalRef.DocRef)
	}

	return nil
}

func (h JournalWs) GetAllDocRefPool(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	result, err := h.svcWebsocket.GetAllDocRefPool(shopID)

	docRefPool := []models.DocRefPool{}
	for tempDocRef, tempUsername := range result {
		docRefPool = append(docRefPool, models.DocRefPool{
			DocRef:   tempDocRef,
			Username: tempUsername,
		})
	}

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    docRefPool,
	})
	return nil
}

func (h JournalWs) SelectDocRefPool(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	username := userInfo.Username

	input := ctx.ReadInput()

	docReq := &models.JournalRef{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	result, err := h.svcWebsocket.DocRefSelect(shopID, username, docReq.DocRef)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    result,
	})
	return nil
}

func (h JournalWs) UnSelectDocRefPool(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	username := userInfo.Username

	input := ctx.ReadInput()

	docReq := &models.JournalRef{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	result, err := h.svcWebsocket.DocRefUnSelect(shopID, username, docReq.DocRef)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    result,
	})
	return nil
}

func (h JournalWs) Check(ctx microservice.IContext) error {

	ctx.Response(http.StatusOK, h.ms.WebsocketCount())

	return nil
}
