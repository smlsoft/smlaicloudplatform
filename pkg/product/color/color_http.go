package color

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/color/models"
	"smlcloudplatform/pkg/product/color/repositories"
	"smlcloudplatform/pkg/product/color/services"
	"smlcloudplatform/pkg/utils"
)

type IColorHttp interface{}

type ColorHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IColorHttpService
}

func NewColorHttp(ms *microservice.Microservice, cfg microservice.IConfig) ColorHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := repositories.NewColorRepository(pst)

	svc := services.NewColorHttpService(repo)

	return ColorHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ColorHttp) RouteSetup() {

	h.ms.POST("/color/bulk", h.SaveBulk)

	h.ms.GET("/color", h.SearchColor)
	h.ms.POST("/color", h.CreateColor)
	h.ms.GET("/color/:id", h.InfoColor)
	h.ms.PUT("/color/:id", h.UpdateColor)
	h.ms.DELETE("/color/:id", h.DeleteColor)
}

// Create Color godoc
// @Description Create Color
// @Tags		Color
// @Param		Color  body      models.Color  true  "Color"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /color [post]
func (h ColorHttp) CreateColor(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.Color{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateColor(shopID, authUsername, *docReq)

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

// Update Color godoc
// @Description Update Color
// @Tags		Color
// @Param		id  path      string  true  "Color ID"
// @Param		Color  body      models.Color  true  "Color"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /color/{id} [put]
func (h ColorHttp) UpdateColor(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Color{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateColor(shopID, id, authUsername, *docReq)

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

// Delete Color godoc
// @Description Delete Color
// @Tags		Color
// @Param		id  path      string  true  "Color ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /color/{id} [delete]
func (h ColorHttp) DeleteColor(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteColor(shopID, id, authUsername)

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

// Get Color godoc
// @Description get struct array by ID
// @Tags		Color
// @Param		id  path      string  true  "Color ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /color/{id} [get]
func (h ColorHttp) InfoColor(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Color %v", id)
	doc, err := h.svc.InfoColor(shopID, id)

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

// List Color godoc
// @Description get struct array by ID
// @Tags		Color
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /color [get]
func (h ColorHttp) SearchColor(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchColor(shopID, q, page, limit, sort)

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

// Create Color Bulk godoc
// @Description Create Color
// @Tags		Color
// @Param		Color  body      []models.Color  true  "Color"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /color/bulk [post]
func (h ColorHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Color{}
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
