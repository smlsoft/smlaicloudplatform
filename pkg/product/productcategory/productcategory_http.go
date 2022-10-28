package productcategory

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/productcategory/models"
	"smlcloudplatform/pkg/product/productcategory/repositories"
	"smlcloudplatform/pkg/product/productcategory/services"
	"smlcloudplatform/pkg/utils"
)

type IProductCategoryHttp interface{}

type ProductCategoryHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IProductCategoryHttpService
}

func NewProductCategoryHttp(ms *microservice.Microservice, cfg microservice.IConfig) ProductCategoryHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := repositories.NewProductCategoryRepository(pst)

	svc := services.NewProductCategoryHttpService(repo)

	return ProductCategoryHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ProductCategoryHttp) RouteSetup() {

	h.ms.POST("/product/category/bulk", h.SaveBulk)

	h.ms.GET("/product/category", h.SearchProductCategoryPage)
	h.ms.GET("/product/category/list", h.SearchProductCategoryLimit)
	h.ms.POST("/product/category", h.CreateProductCategory)
	h.ms.GET("/product/category/:id", h.InfoProductCategory)
	h.ms.PUT("/product/category/:id", h.UpdateProductCategory)
	h.ms.DELETE("/product/category/:id", h.DeleteProductCategory)
}

// Create ProductCategory godoc
// @Description Create ProductCategory
// @Tags		ProductCategory
// @Param		ProductCategory  body      models.ProductCategory  true  "ProductCategory"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/category [post]
func (h ProductCategoryHttp) CreateProductCategory(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.ProductCategory{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateProductCategory(shopID, authUsername, *docReq)

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

// Update ProductCategory godoc
// @Description Update ProductCategory
// @Tags		ProductCategory
// @Param		id  path      string  true  "ProductCategory ID"
// @Param		ProductCategory  body      models.ProductCategory  true  "ProductCategory"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/category/{id} [put]
func (h ProductCategoryHttp) UpdateProductCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.ProductCategory{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateProductCategory(shopID, id, authUsername, *docReq)

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

// Delete ProductCategory godoc
// @Description Delete ProductCategory
// @Tags		ProductCategory
// @Param		id  path      string  true  "ProductCategory ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/category/{id} [delete]
func (h ProductCategoryHttp) DeleteProductCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteProductCategory(shopID, id, authUsername)

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

// Get ProductCategory godoc
// @Description get struct array by ID
// @Tags		ProductCategory
// @Param		id  path      string  true  "ProductCategory ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/category/{id} [get]
func (h ProductCategoryHttp) InfoProductCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get ProductCategory %v", id)
	doc, err := h.svc.InfoProductCategory(shopID, id)

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

// List ProductCategory godoc
// @Description get struct array by ID
// @Tags		ProductCategory
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/category [get]
func (h ProductCategoryHttp) SearchProductCategoryPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchProductCategory(shopID, q, page, limit, sort)

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

// List ProductCategory godoc
// @Description search limit offset
// @Tags		ProductCategory
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/category/list [get]
func (h ProductCategoryHttp) SearchProductCategoryLimit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	offset, limit := utils.GetParamOffsetLimit(ctx.QueryParam)
	sorts := utils.GetSortParam(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchProductCategoryStep(shopID, lang, q, offset, limit, sorts)

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

// Create ProductCategory Bulk godoc
// @Description Create ProductCategory
// @Tags		ProductCategory
// @Param		ProductCategory  body      []models.ProductCategory  true  "ProductCategory"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/category/bulk [post]
func (h ProductCategoryHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.ProductCategory{}
	err := json.Unmarshal([]byte(input), &dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.SaveInBatch(shopID, authUsername, dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		common.ApiResponse{
			Success: true,
		},
	)

	return nil
}
