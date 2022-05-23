package chartofaccount

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/vfgl/chartofaccount/repositories"
	"smlcloudplatform/pkg/vfgl/chartofaccount/services"
	"strconv"
)

type IChartOfAccountHttp interface{}

type ChartOfAccountHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IChartOfAccountHttpService
}

func NewChartOfAccountHttp(ms *microservice.Microservice, cfg microservice.IConfig) ChartOfAccountHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	repo := repositories.NewChartOfAccountRepository(pst)
	mqRepo := repositories.NewChartOfAccountMQRepository(prod)

	svc := services.NewChartOfAccountHttpService(repo, mqRepo)

	return ChartOfAccountHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ChartOfAccountHttp) RouteSetup() {
	h.ms.GET("/gl/chartofaccount", h.Search)
	h.ms.POST("/gl/chartofaccount", h.Create)
	h.ms.GET("/gl/chartofaccount/:id", h.Info)
	h.ms.PUT("/gl/chartofaccount/:id", h.Update)
	h.ms.DELETE("/gl/chartofaccount/:id", h.Delete)
	h.ms.POST("/gl/chartofaccount/bulk", h.SaveBulk)
}

// List Chart Of Account godoc
// @Description List Chart Of Account
// @Tags		GL
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	vfgl.ChartOfAccountPageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/chartofaccount [get]
func (h ChartOfAccountHttp) Search(ctx microservice.IContext) error {
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
	docList, pagination, err := h.svc.Search(shopID, q, page, limit)

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

// Create Chart Of Account godoc
// @Description Create Chart Of Account
// @Tags		GL
// @Param		ChartOfAccount  body      vfgl.ChartOfAccount  true  "ChartOfAccount"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/chartofaccount [post]
func (h ChartOfAccountHttp) Create(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &vfgl.ChartOfAccount{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.Create(shopID, authUsername, *docReq)

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

// Get Chart Of Account Infomation godoc
// @Description Get Chart Of Account
// @Tags		GL
// @Param		id  path      string  true  "Id"
// @Accept 		json
// @Success		200	{object}	vfgl.ChartOfAccountInfo
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/chartofaccount/{id} [get]
func (h ChartOfAccountHttp) Info(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Journal %v", id)
	doc, err := h.svc.Info(id, shopID)

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

// Update Chart Of Account godoc
// @Description Update Chart Of Account
// @Tags		GL
// @Param		id  path      string  true  "ID"
// @Param		ChartOfAccount  body      vfgl.ChartOfAccount  true  "ChartOfAccount"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/chartofaccount/{id} [put]
func (h ChartOfAccountHttp) Update(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &vfgl.ChartOfAccount{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.Update(id, shopID, authUsername, *docReq)

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

// Delete Chart Of Account godoc
// @Description Delete Chart Of Account
// @Tags		GL
// @Param		id  path      string  true  "ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/chartofaccount/{id} [delete]
func (h ChartOfAccountHttp) Delete(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.Delete(id, shopID, authUsername)

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

// Create Chart Of Account Bulk godoc
// @Description Chart Of Account Journal
// @Tags		GL
// @Param		ChartOfAccount  body      []vfgl.ChartOfAccount  true  "ChartOfAccount"
// @Accept 		json
// @Success		201	{object}	models.BulkInsertResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/chartofaccount/bulk [post]
func (h ChartOfAccountHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []vfgl.ChartOfAccount{}
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
