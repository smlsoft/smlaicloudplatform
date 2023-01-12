package kitchen

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/kitchen/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"
)

type IKitchenHttp interface{}

type KitchenHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc IKitchenService
}

func NewKitchenHttp(ms *microservice.Microservice, cfg microservice.IConfig) KitchenHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := NewKitchenRepository(pst)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := NewKitchenService(repo, masterSyncCacheRepo)

	return KitchenHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h KitchenHttp) RouteSetup() {

	h.ms.POST("/restaurant/kitchen/bulk", h.SaveBulk)
	h.ms.GET("/restaurant/kitchen/fetchupdate", h.FetchUpdate)

	h.ms.GET("/restaurant/kitchen", h.SearchKitchen)
	h.ms.GET("/restaurant/kitchen/list", h.SearchKitchenStep)
	h.ms.POST("/restaurant/kitchen", h.CreateKitchen)
	h.ms.GET("/restaurant/kitchen/:id", h.InfoKitchen)
	h.ms.PUT("/restaurant/kitchen/:id", h.UpdateKitchen)
	h.ms.DELETE("/restaurant/kitchen/:id", h.DeleteKitchen)
}

// Create Restaurant Kitchen godoc
// @Description Restaurant Kitchen
// @Tags		Restaurant
// @Param		Kitchen  body      models.Kitchen  true  "Kitchen"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/kitchen [post]
func (h KitchenHttp) CreateKitchen(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.Kitchen{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateKitchen(shopID, authUsername, *docReq)

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

// Update Restaurant Kitchen godoc
// @Description Restaurant Kitchen
// @Tags		Restaurant
// @Param		id  path      string  true  "Kitchen ID"
// @Param		Kitchen  body      models.Kitchen  true  "Kitchen"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/kitchen/{id} [put]
func (h KitchenHttp) UpdateKitchen(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Kitchen{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateKitchen(shopID, id, authUsername, *docReq)

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

// Delete Restaurant Kitchen godoc
// @Description Restaurant Kitchen
// @Tags		Restaurant
// @Param		id  path      string  true  "Kitchen ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/kitchen/{id} [delete]
func (h KitchenHttp) DeleteKitchen(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteKitchen(shopID, id, authUsername)

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

// Get Restaurant Kitchen Infomation godoc
// @Description Get Restaurant Kitchen
// @Tags		Restaurant
// @Param		id  path      string  true  "Kitchen Id"
// @Accept 		json
// @Success		200	{object}	models.KitchenInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/kitchen/{id} [get]
func (h KitchenHttp) InfoKitchen(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Kitchen %v", id)
	doc, err := h.svc.InfoKitchen(shopID, id)

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

// List Restaurant Kitchen godoc
// @Description List Restaurant Kitchen Category
// @Tags		Restaurant
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.KitchenPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/kitchen [get]
func (h KitchenHttp) SearchKitchen(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchKitchen(shopID, q, page, limit)

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

// List Restaurant Kitchen Search Step godoc
// @Description search limit offset
// @Tags		Restaurant
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/kitchen/list [get]
func (h KitchenHttp) SearchKitchenStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	offset, limit := utils.GetParamOffsetLimit(ctx.QueryParam)
	sorts := utils.GetSortParam(ctx.QueryParam)

	docList, total, err := h.svc.SearchKitchenStep(shopID, "", q, offset, limit, sorts)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    docList,
		Total:   total,
	})
	return nil
}

// Fetch Restaurant Kitchen Update By Date godoc
// @Description Fetch Restaurant Kitchen Update By Date
// @Tags		Restaurant
// @Param		lastUpdate query string true "DateTime YYYY-MM-DDTHH:mm"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept		json
// @Success		200 {object} models.KitchenFetchUpdateResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/restaurant/kitchen/fetchupdate [get]
func (h KitchenHttp) FetchUpdate(ctx microservice.IContext) error {
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

	docList, pagination, err := h.svc.LastActivity(shopID, "all", lastUpdate, page, limit)

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

// Create Kitchen Bulk godoc
// @Description Create Kitchen
// @Tags		Restaurant
// @Param		Kitchen  body      []models.Kitchen  true  "Kitchen"
// @Accept 		json
// @Success		201	{object}	common.BulkInsertResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/kitchen/bulk [post]
func (h KitchenHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Kitchen{}
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
