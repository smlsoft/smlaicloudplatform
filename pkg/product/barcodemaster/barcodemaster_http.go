package barcodemaster

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/barcodemaster/models"
	"smlcloudplatform/pkg/product/barcodemaster/repositories"
	"smlcloudplatform/pkg/product/barcodemaster/services"
	categoryRepo "smlcloudplatform/pkg/product/category/repositories"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"
)

type IBarcodeMasterHttp interface {
	RouteSetup()
	CreateBarcodeMaster(ctx microservice.IContext) error
	UpdateBarcodeMaster(ctx microservice.IContext) error
	DeleteBarcodeMaster(ctx microservice.IContext) error
	InfoBarcodeMaster(ctx microservice.IContext) error
	CreateInBatchBarcodeMaster(ctx microservice.IContext) error
	SearchBarcodeMaster(ctx microservice.IContext) error
	LastActivityBarcodeMaster(ctx microservice.IContext) error

	SearchBarcodeMasterOptionMain(ctx microservice.IContext) error
	CreateBarcodeMasterOptionMain(ctx microservice.IContext) error
	DeleteBarcodeMasterOptionMain(ctx microservice.IContext) error
}

type BarcodeMasterHttp struct {
	ms                           *microservice.Microservice
	cfg                          microservice.IConfig
	invService                   services.IBarcodeMasterService
	barcodemasterCategoryService services.IBarcodeMasterCategoryService
}

func NewBarcodeMasterHttp(ms *microservice.Microservice, cfg microservice.IConfig) *BarcodeMasterHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	invRepo := repositories.NewBarcodeMasterRepository(pst)
	invMqRepo := repositories.NewBarcodeMasterMQRepository(prod)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache, "barcodemaster")
	invService := services.NewBarcodeMasterService(invRepo, invMqRepo, masterSyncCacheRepo)

	categoryRepo := categoryRepo.NewCategoryRepository(pst)
	barcodemasterCategoryService := services.NewBarcodeMastercategoryService(invRepo, categoryRepo, invMqRepo)

	return &BarcodeMasterHttp{
		ms:                           ms,
		cfg:                          cfg,
		invService:                   invService,
		barcodemasterCategoryService: barcodemasterCategoryService,
	}
}

func (h BarcodeMasterHttp) RouteSetup() {
	h.ms.GET("/barcodemaster/:id", h.InfoBarcodeMaster)
	h.ms.GET("/barcodemaster", h.SearchBarcodeMaster)
	h.ms.POST("/barcodemaster", h.CreateBarcodeMaster)
	h.ms.POST("/barcodemaster/bulk", h.CreateInBatchBarcodeMaster)
	h.ms.PUT("/barcodemaster/:id", h.UpdateBarcodeMaster)
	h.ms.DELETE("/barcodemaster/:id", h.DeleteBarcodeMaster)
	h.ms.GET("/barcodemaster/fetchupdate", h.LastActivityBarcodeMaster)

	h.ms.POST("/barcodemaster/categoryupdate/:catid", h.UpdateProductCategory)
}

// Create BarcodeMaster godoc
// @Description Create BarcodeMaster
// @Tags		BarcodeMaster
// @Param		BarcodeMaster  body      models.BarcodeMaster  true  "BarcodeMaster"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /barcodemaster [post]
func (h BarcodeMasterHttp) CreateBarcodeMaster(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	barcodemasterReq := &models.BarcodeMaster{}
	err := json.Unmarshal([]byte(input), &barcodemasterReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	_, guidx, err := h.invService.CreateBarcodeMaster(shopID, authUsername, *barcodemasterReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// docIdx := models.BarcodeMasterIndex{}
	// docIdx.ID = idx
	// docIdx.ShopID = userInfo.ShopID
	// docIdx.GuidFixed = guidx

	// err = h.invService.CreateIndex(docIdx)
	// if err != nil {
	// 	return err
	// }

	ctx.Response(
		http.StatusCreated,
		common.ApiResponse{
			Success: true,
			ID:      guidx,
		})

	return nil
}

// Create BarcodeMaster Bulk godoc
// @Description Create BarcodeMaster
// @Tags		BarcodeMaster
// @Param		BarcodeMaster  body      []models.BarcodeMaster  true  "BarcodeMaster"
// @Accept 		json
// @Success		201	{object}	models.BarcodeMasterBulkInsertResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /barcodemaster/bulk [post]
func (h BarcodeMasterHttp) CreateInBatchBarcodeMaster(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	barcodemasterReq := &[]models.BarcodeMaster{}
	err := json.Unmarshal([]byte(input), &barcodemasterReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	barcodemasterBulkResponse, err := h.invService.CreateInBatch(shopID, authUsername, *barcodemasterReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		models.BarcodeMasterBulkReponse{
			Success:                 true,
			BarcodeMasterBulkImport: barcodemasterBulkResponse,
		},
	)

	return nil
}

// Update BarcodeMaster godoc
// @Description Update BarcodeMaster
// @Tags		BarcodeMaster
// @Param		id  path      string  true  "BarcodeMaster ID"
// @Param		BarcodeMaster  body      models.BarcodeMaster  true  "BarcodeMaster"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /barcodemaster/{id} [put]
func (h BarcodeMasterHttp) UpdateBarcodeMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	barcodemasterReq := &models.BarcodeMaster{}
	err := json.Unmarshal([]byte(input), &barcodemasterReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.invService.UpdateBarcodeMaster(shopID, id, authUsername, *barcodemasterReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		common.ApiResponse{
			Success: true,
			ID:      id,
		})

	return nil
}

// Delete BarcodeMaster godoc
// @Description Delete BarcodeMaster
// @Tags		BarcodeMaster
// @Param		id  path      string  true  "BarcodeMaster ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /barcodemaster/{id} [delete]
func (h BarcodeMasterHttp) DeleteBarcodeMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.invService.DeleteBarcodeMaster(shopID, id, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success: true,
			ID:      id,
		},
	)
	return nil
}

// Get BarcodeMaster godoc
// @Description get struct array by ID
// @Tags		BarcodeMaster
// @Param		id  path      string  true  "BarcodeMaster ID"
// @Accept 		json
// @Success		200	{object}	models.BarcodeMasterInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /barcodemaster/{id} [get]
func (h BarcodeMasterHttp) InfoBarcodeMaster(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.invService.InfoBarcodeMaster(shopID, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success: true,
			Data:    doc,
		},
	)

	return nil
}

// List BarcodeMaster godoc
// @Description get struct array by ID
// @Tags		BarcodeMaster
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}		models.BarcodeMasterPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /barcodemaster [get]
func (h BarcodeMasterHttp) SearchBarcodeMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.invService.SearchBarcodeMaster(shopID, q, page, limit)

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

func (h BarcodeMasterHttp) InfoMongoBarcodeMaster(ctx microservice.IContext) error {

	id := ctx.Param("id")

	doc, err := h.invService.InfoMongoBarcodeMaster(id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success: true,
			Data:    doc,
		},
	)

	return nil
}

// Fetch Update BarcodeMaster By Date godoc
// @Description Fetch Update BarcodeMaster By Date
// @Tags		BarcodeMaster
// @Param		lastUpdate query string true "DateTime YYYY-MM-DDTHH:mm"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept		json
// @Success		200 {object} models.BarcodeMasterFetchUpdateResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/barcodemaster/fetchupdate [get]
func (h BarcodeMasterHttp) LastActivityBarcodeMaster(ctx microservice.IContext) error {
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

	docList, pagination, err := h.invService.LastActivity(shopID, lastUpdate, page, limit)

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

// Update BarcodeMaster Category List godoc
// @Description Update BarcodeMaster Category List
// @Tags		BarcodeMaster
// @Param		catid  path      string  true  "Category GUID"
// @Param		BarcodeMaster  body      []string  true  "BarcodeMaster"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccess
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /barcodemaster/categoryupdate/{catid} [post]
func (h BarcodeMasterHttp) UpdateProductCategory(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	catid := ctx.Param("catid")

	input := ctx.ReadInput()

	var barcodemasterReq []string
	err := json.Unmarshal([]byte(input), &barcodemasterReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.barcodemasterCategoryService.UpdateBarcodeMasterCategoryBulk(shopID, authUsername, catid, barcodemasterReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		common.ApiResponse{
			Success: true,
		})

	return nil
}
