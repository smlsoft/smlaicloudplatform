package inventoryimport

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/inventoryimport/models"
	"smlcloudplatform/pkg/utils"

	common "smlcloudplatform/pkg/models"
)

type CategoryImportHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc ICategoryImportService
}

type ICategoryImportHttp interface {
	RouteSetup()
}

func NewCategoryImportHttp(ms *microservice.Microservice, cfg microservice.IConfig) ICategoryImportHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := NewCategoryImportRepository(pst)
	svc := NewCategoryImportService(repo)

	return &CategoryImportHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h *CategoryImportHttp) RouteSetup() {
	h.ms.GET("/import/category", h.ListCategoryImport)
	h.ms.POST("/import/category", h.CreateCategoryImport)
	h.ms.DELETE("/import/category", h.DeleteCategoryImport)
}

// List Category Import godoc
// @Description get struct array by ID
// @Tags		Import
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{array}	models.CategoryImportPageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /import/category [get]
func (h *CategoryImportHttp) ListCategoryImport(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)
	docList, pagination, err := h.svc.List(shopID, pageable)

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

// Create Catagory Import (Bulk) godoc
// @Description Create Catagory Import
// @Tags		Import
// @Param		Catagory  body      []models.CategoryImport  true  "Catagory"
// @Accept 		json
// @Success		201	{object}	models.ResponseSuccess
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /import/category [post]
func (h *CategoryImportHttp) CreateCategoryImport(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := []models.CategoryImport{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.CreateInBatch(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete Category Import godoc
// @Description Delete Category
// @Tags		Import
// @Param		id  body      []string  true  "Category Import ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccess
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /import/category [delete]
func (h *CategoryImportHttp) DeleteCategoryImport(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	docReq := []string{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.Delete(shopID, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}
