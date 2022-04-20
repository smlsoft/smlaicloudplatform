package inventoryimport

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IInventoryImporOptionMaintHttp interface {
	RouteSetup()
}

type InventoryImporOptionMaintHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc IInventoryOptionMainImportService
}

func NewInventoryImporOptionMaintHttp(ms *microservice.Microservice, cfg microservice.IConfig) InventoryImporOptionMaintHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	invImportOptionMainRepo := NewInventoryOptionMainImportRepository(pst)
	invImportOptionMainService := NewInventoryOptionMainImportService(invImportOptionMainRepo)

	return InventoryImporOptionMaintHttp{
		ms:  ms,
		cfg: cfg,
		svc: invImportOptionMainService,
	}
}

func (h *InventoryImporOptionMaintHttp) RouteSetup() {

	h.ms.GET("/optionimport", h.ListInventoryOptionMain)
	h.ms.POST("/optionimport", h.CreateInventoryOptionMain)
	h.ms.DELETE("/optionimport", h.DeleteInventoryOptionMain)

}

// Create Inventory Option godoc
// @Description Create Inventory Option
// @Tags		Import
// @Param		Option  body      models.InventoryOptionMain  true  "Option"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
// @Security     AccessToken
// @Router /optionimport [post]
func (h InventoryImporOptionMaintHttp) CreateInventoryOptionMain(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := []models.InventoryOptionMainImport{}
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

// Delete Option godoc
// @Description Delete Option
// @Tags		Import
// @Param		id  path      string  true  "Option ID"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
// @Security     AccessToken
// @Router /optionimport/{id} [delete]
func (h InventoryImporOptionMaintHttp) DeleteInventoryOptionMain(ctx microservice.IContext) error {
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

// List Inventory Option godoc
// @Description List Inventory Option
// @Tags		Import
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
// @Security     AccessToken
// @Router /optionimport [get]
func (h InventoryImporOptionMaintHttp) ListInventoryOptionMain(ctx microservice.IContext) error {
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
