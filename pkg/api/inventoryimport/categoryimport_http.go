package inventoryimport

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
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
	h.ms.GET("/categoryimport", h.ListCategoryImport)
	h.ms.POST("/categoryimport", h.CreateCategoryImport)
	h.ms.DELETE("/categoryimport", h.DeleteCategoryImport)
}

func (h *CategoryImportHttp) ListCategoryImport(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}
	docList, pagination, err := h.svc.ListInventory(shopID, page, limit)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success:    true,
		Data:       docList,
		Pagination: pagination,
	})
	return nil
}

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

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
	})
	return nil
}

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

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
	})

	return nil
}
