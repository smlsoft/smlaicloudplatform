package journal

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/vfgl/journal/config"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
	"smlcloudplatform/pkg/vfgl/journal/services"
	"time"
	"unicode/utf8"

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
	svcWebsocket := services.NewJournalWebsocketService(cacheRepo, time.Duration(30)*time.Minute)

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
	h.ms.GET("/gl/journal/docref/selected", h.GetAllDocRefPool)
	h.ms.POST("/gl/journal/docref/select", h.SelectDocRefPool)
	h.ms.POST("/gl/journal/docref/unselect", h.UnSelectDocRefPool)
	h.ms.GET("/gl/journal/user-docref", h.GetUserDocRef)
	h.ms.GET("/gl/journal/docref-user", h.GetDocRefUser)

}

func (h JournalWs) WebsocketImage(ctx microservice.IContext) error {

	screenName := config.WEBSOCKET_SCREEN_IMAGE

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	username := userInfo.Username

	if username == "" {
		ctx.Response(http.StatusBadRequest, "username parameter missing")
		return nil
	}

	socketID := utils.NewGUID()

	ws, err := h.ms.Websocket(socketID, ctx.ResponseWriter(), ctx.Request())

	if err != nil {
		return err
	}

	err = h.svcWebsocket.SetWebsocket(shopID, username, screenName)

	if err != nil {
		return err
	}

	err = h.svcWebsocket.ExpireWebsocket(shopID, username)

	if err != nil {
		return err
	}

	cacheMsg, subID, err := h.svcWebsocket.SubDoc(shopID, username, screenName)

	if err != nil {
		return err
	}

	defer func() {
		ws.Close()
		h.ms.WebsocketClose(socketID)
		h.svcWebsocket.UnSub(subID)
		h.svcWebsocket.DelWebsocket(shopID, username, screenName)
		h.ClearDocRef(shopID, username)
	}()

	// Send to client
	for {
		h.svcWebsocket.ExpireWebsocket(shopID, username)

		_, _, err := ws.ReadMessage()
		if err != nil {
			return nil
		}

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

	screenName := config.WEBSOCKET_SCREEN_FORM
	sendScreenName := config.WEBSOCKET_SCREEN_IMAGE
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	username := userInfo.Username

	if username == "" {
		ctx.Response(http.StatusBadRequest, "username parameter missing")
		return nil
	}

	socketID := utils.NewID()

	ws, err := h.ms.Websocket(socketID, ctx.ResponseWriter(), ctx.Request())

	if err != nil {
		return err
	}

	err = h.svcWebsocket.SetWebsocket(shopID, username, screenName)

	if err != nil {
		return err
	}

	err = h.svcWebsocket.ExpireWebsocket(shopID, username)

	if err != nil {
		return err
	}

	cacheMsg, subID, err := h.svcWebsocket.SubDoc(shopID, username, screenName)

	if err != nil {
		return err
	}

	sigClose := make(chan struct{})
	defer func() {
		ws.Close()
		h.ms.WebsocketClose(socketID)
		h.svcWebsocket.UnSub(subID)
		h.svcWebsocket.DelWebsocket(shopID, username, screenName)
		h.ClearDocRef(shopID, username)
	}()

	// Receive from client
	go func(ws *websocket.Conn, sigClose chan struct{}) {
		defer func() {
			sigClose <- struct{}{}
		}()

		for {
			h.svcWebsocket.ExpireWebsocket(shopID, username)

			journalCommand := models.JournalCommand{}
			err := ws.ReadJSON(&journalCommand)

			if err != nil {
				return
			}

			tempRef, _ := json.Marshal(journalCommand)
			h.svcWebsocket.PubDoc(shopID, username, sendScreenName, tempRef)

			switch journalCommand.Command {
			case "save":
				h.svcWebsocket.ClearLastMessage(shopID, username)
				//clear
				h.ClearDocRef(shopID, username)
			}
		}

	}(ws, sigClose)

	// Send to client
	lastMessage, err := h.svcWebsocket.GetLastMessage(shopID, username, screenName)
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
		select {
		case temp := <-cacheMsg:
			if temp != nil {
				err = ws.WriteMessage(websocket.TextMessage, []byte(temp.Payload))
				if err != nil {
					return err
				}
			}
		case <-sigClose:
			return nil
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

func (h JournalWs) ClearDocRef(shopID string, username string) error {

	// isExists, err := h.svcWebsocket.ExistsWebsocket(shopID, username)
	// if err != nil {
	// 	h.ms.Logger.Error(err.Error())
	// }

	// if isExists {
	// 	return nil
	// }

	// lastMessage, err := h.svcWebsocket.GetLastMessage(shopID, username, "form")
	// if err != nil {
	// 	h.ms.Logger.Error(err.Error())
	// }

	// if lastMessage != "" {
	// 	journalRef := models.JournalRef{}
	// 	json.Unmarshal([]byte(lastMessage), &journalRef)

	// 	err = h.svcWebsocket.DelDocRefPool(shopID, username, journalRef.DocRef)
	// 	if err != nil {
	// 		h.ms.Logger.Error(err.Error())
	// 	}
	// }

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

func (h JournalWs) GetUserDocRef(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	username := userInfo.Username

	result, err := h.svcWebsocket.GetDocRefUserPool(shopID, username)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data: &models.JournalRef{
			DocRef: result,
		},
	})
	return nil
}

func (h JournalWs) GetDocRefUser(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	docRef := ctx.QueryParam("docref")

	if docRef == "" {
		ctx.ResponseError(http.StatusBadRequest, "docref is empty")
		return nil
	}

	if utf8.RuneCountInString(docRef) > 100 {
		ctx.ResponseError(http.StatusBadRequest, "docref invalid")
		return nil
	}

	result, err := h.svcWebsocket.GetDocRefPool(shopID, docRef)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data: map[string]string{
			"username": result,
		},
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

	forceSelectRaw := ctx.QueryParam("force")
	forceSelect := false
	if forceSelectRaw == "1" {
		forceSelect = true
	}

	result, err := h.svcWebsocket.DocRefSelectForce(shopID, username, docReq.DocRef, forceSelect)

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

	// input := ctx.ReadInput()

	// docReq := &models.JournalRef{}
	// err := json.Unmarshal([]byte(input), &docReq)

	// if err != nil {
	// 	ctx.ResponseError(400, err.Error())
	// 	return err
	// }

	result, err := h.svcWebsocket.DocRefUnSelect(shopID, username)

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

func (h JournalWs) WebsocketConnectCount(ctx microservice.IContext) error {

	ctx.Response(http.StatusOK, h.ms.WebsocketCount())

	return nil
}
