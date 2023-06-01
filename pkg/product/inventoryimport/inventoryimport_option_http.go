package inventoryimport

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/inventoryimport/models"
	"smlcloudplatform/pkg/utils"
)

type IInventoryImporOptionMaintHttp interface {
	RouteSetup()
}

type InventoryImporOptionMaintHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc IInventoryOptionMainImportService
}

func NewInventoryImporOptionMaintHttp(ms *microservice.Microservice, cfg microservice.IConfig) *InventoryImporOptionMaintHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	invImportOptionMainRepo := NewInventoryOptionMainImportRepository(pst)
	invImportOptionMainService := NewInventoryOptionMainImportService(invImportOptionMainRepo)

	return &InventoryImporOptionMaintHttp{
		ms:  ms,
		cfg: cfg,
		svc: invImportOptionMainService,
	}
}

func (h *InventoryImporOptionMaintHttp) RouteSetup() {

	h.ms.GET("/import/option", h.ListInventoryOptionMain)
	h.ms.POST("/import/option", h.CreateInventoryOptionMain)
	h.ms.DELETE("/import/option", h.DeleteInventoryOptionMain)

}

// Create Inventory Option godoc
// @Description Create Inventory Option
// @Tags		Import
// @Param		Option  body	[]models.InventoryOptionMainImport  true  "Option"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccess
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /import/option [post]
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

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete Option godoc
// @Description Delete Option
// @Tags		Import
// @Param		id  body      []string  true  "Option ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccess
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /import/option [delete]
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

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// List Inventory Option godoc
// @Description List Inventory Option
// @Tags		Import
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.InventoryOptionPageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /import/option [get]
func (h InventoryImporOptionMaintHttp) ListInventoryOptionMain(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)
	docList, pagination, err := h.svc.ListInventory(shopID, pageable)

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
