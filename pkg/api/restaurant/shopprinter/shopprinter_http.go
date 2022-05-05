package shopprinter

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"strconv"
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

	crudRepo := repositories.NewCrudRepository[restaurant.PrinterTerminalDoc](pst)
	searchRepo := repositories.NewSearchRepository[restaurant.PrinterTerminalInfo](pst)
	guidRepo := repositories.NewGuidRepository[restaurant.PrinterTerminalItemGuid](pst)
	activityRepo := repositories.NewActivityRepository[restaurant.PrinterTerminalActivity, restaurant.PrinterTerminalDeleteActivity](pst)
	svc := NewShopPrinterService(crudRepo, searchRepo, guidRepo, activityRepo)

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

func (h ShopPrinterHttp) CreateShopPrinter(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &restaurant.PrinterTerminal{}
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

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      idx,
	})
	return nil
}

func (h ShopPrinterHttp) UpdateShopPrinter(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &restaurant.PrinterTerminal{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateShopPrinter(id, shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

func (h ShopPrinterHttp) DeleteShopPrinter(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteShopPrinter(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

func (h ShopPrinterHttp) InfoShopPrinter(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get ShopPrinter %v", id)
	doc, err := h.svc.InfoShopPrinter(id, shopID)

	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v: %v", id, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

func (h ShopPrinterHttp) SearchShopPrinter(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}
	docList, pagination, err := h.svc.SearchShopPrinter(shopID, q, page, limit)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success:    true,
		Data:       docList,
		Pagination: pagination,
	})
	return nil
}

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

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	docList, pagination, err := h.svc.LastActivity(shopID, lastUpdate, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})

	return nil
}

func (h ShopPrinterHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []restaurant.PrinterTerminal{}
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
		models.BulkReponse{
			Success:    true,
			BulkImport: bulkResponse,
		},
	)

	return nil
}
