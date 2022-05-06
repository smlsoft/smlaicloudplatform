package shoptable

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

type IShopTableHttp interface{}

type ShopTableHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc IShopTableService
}

func NewShopTableHttp(ms *microservice.Microservice, cfg microservice.IConfig) ShopTableHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	crudRepo := repositories.NewCrudRepository[restaurant.ShopTableDoc](pst)
	searchRepo := repositories.NewSearchRepository[restaurant.ShopTableInfo](pst)
	guidRepo := repositories.NewGuidRepository[restaurant.ShopTableItemGuid](pst)
	activityRepo := repositories.NewActivityRepository[restaurant.ShopTableActivity, restaurant.ShopTableDeleteActivity](pst)

	svc := NewShopTableService(crudRepo, searchRepo, guidRepo, activityRepo)

	return ShopTableHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ShopTableHttp) RouteSetup() {

	h.ms.POST("/restaurant/table/bulk", h.SaveBulk)
	h.ms.GET("/restaurant/table/fetchupdate", h.FetchUpdate)

	h.ms.GET("/restaurant/table", h.SearchShopTable)
	h.ms.POST("/restaurant/table", h.CreateShopTable)
	h.ms.GET("/restaurant/table/:id", h.InfoShopTable)
	h.ms.PUT("/restaurant/table/:id", h.UpdateShopTable)
	h.ms.DELETE("/restaurant/table/:id", h.DeleteShopTable)

}

// Create Restaurant Shop Table godoc
// @Description Restaurant Shop Table
// @Tags		Restaurant
// @Param		Table  body      restaurant.ShopTable  true  "Table"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/table [post]
func (h ShopTableHttp) CreateShopTable(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &restaurant.ShopTable{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateShopTable(shopID, authUsername, *docReq)

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

// Update Restaurant Shop Table godoc
// @Description Restaurant Shop Table
// @Tags		Restaurant
// @Param		id  path      string  true  "Table ID"
// @Param		Table  body      restaurant.ShopTable  true  "Table"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/table/{id} [put]
func (h ShopTableHttp) UpdateShopTable(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &restaurant.ShopTable{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateShopTable(id, shopID, authUsername, *docReq)

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

// Delete Restaurant Shop Table godoc
// @Description Restaurant Shop Table
// @Tags		Restaurant
// @Param		id  path      string  true  "ShopTable ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/table/{id} [delete]
func (h ShopTableHttp) DeleteShopTable(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteShopTable(id, shopID, authUsername)

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

// Get Restaurant Shop Table Infomation godoc
// @Description Get Restaurant Shop Table
// @Tags		Restaurant
// @Param		id  path      string  true  "ShopTable Id"
// @Accept 		json
// @Success		200	{object}	restaurant.ShopTableInfoResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/table/{id} [get]
func (h ShopTableHttp) InfoShopTable(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get ShopTable %v", id)
	doc, err := h.svc.InfoShopTable(id, shopID)

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

// List Restaurant Shop Table godoc
// @Description List Restaurant Shop Table Category
// @Tags		Restaurant
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	restaurant.ShopTablePageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/table [get]
func (h ShopTableHttp) SearchShopTable(ctx microservice.IContext) error {
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
	docList, pagination, err := h.svc.SearchShopTable(shopID, q, page, limit)

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

// Fetch Restaurant ShopTable Update By Date godoc
// @Description Fetch Restaurant ShopTable Update By Date
// @Tags		Restaurant
// @Param		lastUpdate query string true "DateTime YYYY-MM-DDTHH:mm"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept		json
// @Success		200 {object} restaurant.ShopTableFetchUpdateResponse
// @Failure		401 {object} models.AuthResponseFailed
// @Security	AccessToken
// @Router		/restaurant/table/fetchupdate [get]
func (h ShopTableHttp) FetchUpdate(ctx microservice.IContext) error {
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

// Create ShopTable Bulk godoc
// @Description Create ShopTable
// @Tags		Restaurant
// @Param		ShopTable  body      []restaurant.ShopTable  true  "ShopTable"
// @Accept 		json
// @Success		201	{object}	models.BulkInsertResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/table/bulk [post]
func (h ShopTableHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []restaurant.ShopTable{}
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
