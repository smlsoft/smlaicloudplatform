package category

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"
)

type ICategoryHttp interface{}

type CategoryHttp struct {
	ms          *microservice.Microservice
	cfg         microservice.IConfig
	cateService ICategoryService
}

func NewCategoryHttp(ms *microservice.Microservice, cfg microservice.IConfig) CategoryHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	cateRepo := NewCategoryRepository(pst)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache, "category")
	cateService := NewCategoryService(cateRepo, masterSyncCacheRepo)

	return CategoryHttp{
		ms:          ms,
		cfg:         cfg,
		cateService: cateService,
	}
}

func (h CategoryHttp) RouteSetup() {

	h.ms.GET("/category/:id", h.InfoCategory)
	h.ms.GET("/category", h.SearchCategory)
	h.ms.POST("/category", h.CreateCategory)
	h.ms.POST("/category/bulk", h.CreateInBatchCategory)
	h.ms.PUT("/category/:id", h.UpdateCategory)
	h.ms.DELETE("/category/:id", h.DeleteCategory)
	h.ms.GET("/category/fetchupdate", h.LastActivityCategory)
}

// Create Category godoc
// @Description Create Inventory Category
// @Tags		Inventory
// @Param		Category  body      models.Category  true  "Add Category"
// @Accept 		json
// @Success		201	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /category [post]
func (h CategoryHttp) CreateCategory(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	categoryReq := &models.Category{}
	err := json.Unmarshal([]byte(input), &categoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(categoryReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.cateService.CreateCategory(shopID, authUsername, *categoryReq)

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

// Update Inventory godoc
// @Description Update Inventory
// @Tags		Inventory
// @Param		id  path      string  true  "Category ID"
// @Param		Category  body      models.Category  true  "Category"
// @Accept 		json
// @Success		201	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /category/{id} [put]
func (h CategoryHttp) UpdateCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	categoryReq := &models.Category{}
	err := json.Unmarshal([]byte(input), &categoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(categoryReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.cateService.UpdateCategory(shopID, id, authUsername, *categoryReq)

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

// Delete Category godoc
// @Description Delete Category
// @Tags		Inventory
// @Param		id  path      string  true  "Category ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /category/{id} [delete]
func (h CategoryHttp) DeleteCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.cateService.DeleteCategory(shopID, id, authUsername)

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

// Get Category Infomation godoc
// @Description Get Inventory Category
// @Tags		Inventory
// @Param		id  path      string  true  "Category Id"
// @Accept 		json
// @Success		200	{object}	models.CategoryInfoResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /category/{id} [get]
func (h CategoryHttp) InfoCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Category %v", id)
	doc, err := h.cateService.InfoCategory(shopID, id)

	if err != nil {
		h.ms.Logger.Errorf("Error getting category %v: %v", id, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// List Category godoc
// @Description List Inventory Category
// @Tags		Inventory
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{object}	models.CategoryPageResponse
// @Failure		400 {object}	models.AuthResponseFailed
// @Failure		500 {object}	models.AuthResponseFailed
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /category [get]
func (h CategoryHttp) SearchCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")

	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.cateService.SearchCategory(shopID, q, page, limit)

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

// Fetch Update Category By Date godoc
// @Description Fetch Update Category By Date
// @Tags		Inventory
// @Param		lastUpdate query string true "DateTime"
// @Accept		json
// @Success		200 {object} models.CategoryFetchUpdateResponse
// @Failure		401 {object} models.AuthResponseFailed
// @Security	AccessToken
// @Router		/category/fetchupdate [get]
func (h CategoryHttp) LastActivityCategory(ctx microservice.IContext) error {
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

	docList, pagination, err := h.cateService.LastActivity(shopID, lastUpdate, page, limit)

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

// Create Category Bulk godoc
// @Description Create Category
// @Tags		Inventory
// @Param		Category  body      []models.Category  true  "Category"
// @Accept 		json
// @Success		201	{object}	models.CategoryBulkReponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /category/bulk [post]
func (h CategoryHttp) CreateInBatchCategory(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Category{}
	err := json.Unmarshal([]byte(input), &dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if len(dataReq) < 1 {
		ctx.ResponseError(400, "Require category more than one")
		return err
	}

	bulkResponse, err := h.cateService.SaveInBatch(shopID, authUsername, dataReq)

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
