package journal

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/vfgl/journal/models"
	"smlcloudplatform/pkg/vfgl/journal/repositories"
	"smlcloudplatform/pkg/vfgl/journal/services"

	"github.com/gorilla/websocket"
)

type IJournalHttp interface{}

type JournalHttp struct {
	ms           *microservice.Microservice
	cfg          microservice.IConfig
	svc          services.IJournalHttpService
	svcWebsocket services.IJournalWebsocketService
}

func NewJournalHttp(ms *microservice.Microservice, cfg microservice.IConfig) JournalHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewJournalRepository(pst)
	mqRepo := repositories.NewJournalMqRepository(prod)
	svc := services.NewJournalHttpService(repo, mqRepo)

	cacheRepo := repositories.NewJournalCacheRepository(cache)
	svcWebsocket := services.NewJournalWebsocketService(cacheRepo)

	return JournalHttp{
		ms:           ms,
		cfg:          cfg,
		svc:          svc,
		svcWebsocket: svcWebsocket,
	}
}

func (h JournalHttp) RouteSetup() {

	h.ms.POST("/gl/journal/bulk", h.SaveBulk)

	h.ms.GET("/gl/journal", h.SearchJournal)
	h.ms.POST("/gl/journal", h.CreateJournal)
	h.ms.GET("/gl/journal/:id", h.InfoJournal)
	h.ms.PUT("/gl/journal/:id", h.UpdateJournal)
	h.ms.DELETE("/gl/journal/:id", h.DeleteJournal)
	h.ms.GET("/gl/journal/ws/image", h.WebsocketImage)
	h.ms.GET("/gl/journal/ws/form", h.WebsocketForm)
	h.ms.GET("/pubx", h.Pub)
	h.ms.GET("/checkx", h.Check)
}

// Create Journal godoc
// @Summary		บันทึกข้อมูลรายวัน
// @Description บันทึกข้อมูลรายวัน
// @Tags		GL
// @Param		Journal  body      models.Journal  true  "ข้อมูลรายวัน"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/journal [post]
func (h JournalHttp) CreateJournal(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.Journal{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateJournal(shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})
	return nil
}

// Update Journal godoc
// @Summary		แก้ไขข้อมูลรายวัน
// @Description แก้ไขข้อมูลรายวัน
// @Tags		GL
// @Param		id  path      string  true  "Journal ID"
// @Param		Journal  body      models.Journal  true  "ข้อมูลรายวัน"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/journal/{id} [put]
func (h JournalHttp) UpdateJournal(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Journal{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateJournal(id, shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

// Delete Journal godoc
// @Summary		ลบข้อมูลรายวัน
// @Description ลบข้อมูลรายวัน
// @Tags		GL
// @Param		id  path      string  true  "Journal ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/journal/{id} [delete]
func (h JournalHttp) DeleteJournal(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteJournal(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

// Get Journal Infomation godoc
// @Summary		แสดงรายละเอียดข้อมูลรายวัน
// @Description แสดงรายละเอียดข้อมูลรายวัน
// @Tags		GL
// @Param		id  path      string  true  "Journal Id"
// @Accept 		json
// @Success		200	{object}	models.JournalInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/journal/{id} [get]
func (h JournalHttp) InfoJournal(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Journal %v", id)
	doc, err := h.svc.InfoJournal(id, shopID)

	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v: %v", id, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// List Journal godoc
// @Summary		แสดงรายการข้อมูลรายวัน
// @Description แสดงรายการข้อมูลรายวัน
// @Tags		GL
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.JournalPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/journal [get]
func (h JournalHttp) SearchJournal(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchJournal(shopID, q, page, limit, sort)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Data:       docList,
		Pagination: pagination,
	})
	return nil
}

// Create Journal Bulk godoc
// @Summary		นำเข้าข้อมูลรายวัน
// @Description นำเข้าข้อมูลรายวัน
// @Tags		GL
// @Param		Journal  body      []models.Journal  true  "ข้อมูลรายวัน"
// @Accept 		json
// @Success		201	{object}	common.BulkInsertResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/journal/bulk [post]
func (h JournalHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Journal{}
	err := json.Unmarshal([]byte(input), &dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	bulkResponse, err := h.svc.SaveInBatch(shopID, authUsername, dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		common.BulkReponse{
			Success:    true,
			BulkImport: bulkResponse,
		},
	)

	return nil
}

func (h JournalHttp) Pub(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	docID := userInfo.Username
	q := ctx.QueryParam("q")

	if docID == "" {
		ctx.Response(http.StatusBadRequest, "doc parameter missing")
		return nil
	}

	journalRef := models.JournalRef{
		DocRef: q,
	}

	// journalRef := "send ok " + q

	tempx, _ := json.Marshal(journalRef)

	h.svcWebsocket.Pub(shopID, docID, "image", tempx)

	ctx.Response(http.StatusOK, docID)

	return nil
}

func (h JournalHttp) WebsocketImage(ctx microservice.IContext) error {

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

	defer ws.Close()

	err = h.svcWebsocket.SetWebsocket(shopID, processID, screenName, socketID)

	if err != nil {
		return err
	}

	cacheMsg, subID, err := h.svcWebsocket.Sub(shopID, processID, screenName)

	if err != nil {
		return err
	}

	defer func() {
		h.ms.WebsocketClose(socketID)
		h.svcWebsocket.UnSub(subID)
		h.svcWebsocket.DelWebsocket(shopID, processID, screenName, socketID)

	}()

	// Receive
	go func(ws *websocket.Conn) {
		for {
			journalRef := models.JournalRef{}
			err := ws.ReadJSON(&journalRef)

			if err != nil {
				return
			}

			tempRef, _ := json.Marshal(journalRef)
			h.svcWebsocket.Pub(shopID, processID, sendScreenName, tempRef)

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
			err := h.svcWebsocket.SaveLastMessage(shopID, processID, screenName, temp.Payload)
			if err != nil {
				h.ms.Logger.Error(err.Error())
			}

			err = ws.WriteMessage(websocket.TextMessage, []byte(temp.Payload))
			if err != nil {
				return err
			}

		}
	}

}

func (h JournalHttp) WebsocketForm(ctx microservice.IContext) error {

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

	defer ws.Close()

	err = h.svcWebsocket.SetWebsocket(shopID, processID, screenName, socketID)

	if err != nil {
		return err
	}

	cacheMsg, subID, err := h.svcWebsocket.Sub(shopID, processID, screenName)

	if err != nil {
		return err
	}

	defer func() {
		h.ms.WebsocketClose(socketID)
		h.svcWebsocket.UnSub(subID)
		h.svcWebsocket.DelWebsocket(shopID, processID, screenName, socketID)
	}()

	// Receive
	go func(ws *websocket.Conn) {
		for {
			journalRef := models.JournalRef{}
			err := ws.ReadJSON(&journalRef)

			if err != nil {
				return
			}

			tempRef, _ := json.Marshal(journalRef)
			h.svcWebsocket.Pub(shopID, processID, sendScreenName, tempRef)
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
			err := h.svcWebsocket.SaveLastMessage(shopID, processID, screenName, temp.Payload)
			if err != nil {
				h.ms.Logger.Error(err.Error())
			}

			err = ws.WriteMessage(websocket.TextMessage, []byte(temp.Payload))
			if err != nil {
				return err
			}
		}
	}

}

func (h JournalHttp) Check(ctx microservice.IContext) error {

	ctx.Response(http.StatusOK, h.ms.WebsocketCount())

	return nil
}
