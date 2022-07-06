package shopprinter

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/shopprinter/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"
)

type IShopPrinterHttp interface{}

type ShopPrinterHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc IShopPrinterService
}

func NewShopPrinterHttp(ms *microservice.Microservice, cfg microservice.IConfig) ShopPrinterHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := NewShopPrinterRepository(pst)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache, "shopprinter")
	svc := NewShopPrinterService(repo, masterSyncCacheRepo)

	return ShopPrinterHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ShopPrinterHttp) RouteSetup() {

	h.ms.POST("/restaurant/printer/bulk", h.SaveBulk)
	h.ms.GET("/restaurant/printer/fetchupdate", h.FetchUpdate)

	h.ms.GET("/restaurant/printer", h.SearchShopPrinter)
	h.ms.POST("/restaurant/printer", h.CreateShopPrinter)
	h.ms.GET("/restaurant/printer/:id", h.InfoShopPrinter)
	h.ms.PUT("/restaurant/printer/:id", h.UpdateShopPrinter)
	h.ms.DELETE("/restaurant/printer/:id", h.DeleteShopPrinter)

}

// Create Restaurant Printer godoc
// @Description Restaurant Printer
// @Tags		Restaurant
// @Param		Printer  body      models.PrinterTerminal  true  "Printer"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/printer [post]
func (h ShopPrinterHttp) CreateShopPrinter(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.PrinterTerminal{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateShopPrinter(shopID, authUsername, *docReq)

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

// Update Restaurant Printer godoc
// @Description Restaurant Printer
// @Tags		Restaurant
// @Param		id  path      string  true  "Printer ID"
// @Param		Printer  body      models.PrinterTerminal  true  "Printer"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/printer/{id} [put]
func (h ShopPrinterHttp) UpdateShopPrinter(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.PrinterTerminal{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateShopPrinter(shopID, id, authUsername, *docReq)

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

// Delete Restaurant Printer godoc
// @Description Restaurant Printer
// @Tags		Restaurant
// @Param		id  path      string  true  "Printer ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/printer/{id} [delete]
func (h ShopPrinterHttp) DeleteShopPrinter(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteShopPrinter(shopID, id, authUsername)

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

// Get Restaurant Printer Infomation godoc
// @Description Get Restaurant Printer
// @Tags		Restaurant
// @Param		id  path      string  true  "Printer Id"
// @Accept 		json
// @Success		200	{object}	models.PrinterTerminalInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/printer/{id} [get]
func (h ShopPrinterHttp) InfoShopPrinter(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get ShopPrinter %v", id)
	doc, err := h.svc.InfoShopPrinter(shopID, id)

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

// List Restaurant Printer godoc
// @Description List Restaurant Printer Category
// @Tags		Restaurant
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.PrinterTerminalPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/printer [get]
func (h ShopPrinterHttp) SearchShopPrinter(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchShopPrinter(shopID, q, page, limit)

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

// Fetch Restaurant Printer Update By Date godoc
// @Description Fetch Restaurant Printer Update By Date
// @Tags		Restaurant
// @Param		lastUpdate query string true "DateTime YYYY-MM-DDTHH:mm"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept		json
// @Success		200 {object} models.PrinterTerminalFetchUpdateResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/restaurant/printer/fetchupdate [get]
func (h ShopPrinterHttp) FetchUpdate(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02T15:04" //
	lastUpdateStr := ctx.QueryParam("lastUpdate")

	lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(lastUpdateStr) < 1 {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return nil
	}

	lastUpdate, err := time.Parse(layout, lastUpdateStr)

	if err != nil {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return err
	}

	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.svc.LastActivity(shopID, lastUpdate, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})

	return nil
}

// Create Printer Bulk godoc
// @Description Printer ShopZone
// @Tags		Restaurant
// @Param		Printer  body      []models.PrinterTerminal  true  "Printer"
// @Accept 		json
// @Success		201	{object}	common.BulkInsertResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/printer/bulk [post]
func (h ShopPrinterHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.PrinterTerminal{}
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
