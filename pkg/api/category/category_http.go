package category

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type ICategoryHttp interface{}

type CategoryHttp struct {
	ms          *microservice.Microservice
	cfg         microservice.IConfig
	cateService ICategoryService
}

func NewCategoryHttp(ms *microservice.Microservice, cfg microservice.IConfig) CategoryHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	cateRepo := NewCategoryRepository(pst)
	cateService := NewCategoryService(cateRepo)

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
	h.ms.PUT("/category/:id", h.UpdateCategory)
	h.ms.DELETE("/category/:id", h.DeleteCategory)
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

	err = h.cateService.UpdateCategory(id, shopID, authUsername, *categoryReq)

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

	err := h.cateService.DeleteCategory(id, shopID, authUsername)

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
	doc, err := h.cateService.InfoCategory(id, shopID)

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
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}
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
