package inventoryimport

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type InventoryImportHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc IInventoryImportService
}

type IInventoryImportHttp interface {
	RouteSetup()
}

func NewInventoryImportHttp(ms *microservice.Microservice, cfg microservice.IConfig) IInventoryImportHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	invRepo := NewInventoryImportRepository(pst)
	invImportService := NewInventoryImportService(invRepo)

	return &InventoryImportHttp{
		ms:  ms,
		cfg: cfg,
		svc: invImportService,
	}
}

func (h *InventoryImportHttp) RouteSetup() {

	h.ms.GET("/import/inventory", h.ListInventoryImport)
	h.ms.POST("/import/inventory", h.CreateInventoryImport)
	h.ms.DELETE("/import/inventory", h.DeleteInventoryImport)

}

// List Inventory Import godoc
// @Description get struct array by ID
// @Tags		Import
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{array}	models.InventoryPageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /import/inventory [get]
func (h *InventoryImportHttp) ListInventoryImport(ctx microservice.IContext) error {
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

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success:    true,
		Data:       docList,
		Pagination: pagination,
	})
	return nil
}

// Create Inventory Import (Bulk) godoc
// @Description Create Inventory Import
// @Tags		Import
// @Param		Inventory  body      []models.Inventory  true  "Inventory"
// @Accept 		json
// @Success		201	{object}	models.ResponseSuccess
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /import/inventory [post]
func (h *InventoryImportHttp) CreateInventoryImport(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := []models.InventoryImport{}
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

// Delete Inventory Import godoc
// @Description Delete Inventory
// @Tags		Import
// @Param		id  body      []string  true  "Inventory Import ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccess
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /import/inventory [delete]
func (h *InventoryImportHttp) DeleteInventoryImport(ctx microservice.IContext) error {
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

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})

	return nil
}
