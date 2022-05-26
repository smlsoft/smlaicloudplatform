package inventory

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/category"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"
)

type IInventoryHttp interface {
	RouteSetup()
	CreateInventory(ctx microservice.IContext) error
	UpdateInventory(ctx microservice.IContext) error
	DeleteInventory(ctx microservice.IContext) error
	InfoInventory(ctx microservice.IContext) error
	CreateInBatchInventory(ctx microservice.IContext) error
	SearchInventory(ctx microservice.IContext) error
	LastActivityInventory(ctx microservice.IContext) error

	SearchInventoryOptionMain(ctx microservice.IContext) error
	CreateInventoryOptionMain(ctx microservice.IContext) error
	DeleteInventoryOptionMain(ctx microservice.IContext) error
}

type InventoryHttp struct {
	ms                       *microservice.Microservice
	cfg                      microservice.IConfig
	invService               IInventoryService
	invOptService            IInventoryOptionMainService
	optGroupService          IOptionGroupService
	inventoryCategoryService IInventoryCategoryService
}

func NewInventoryHttp(ms *microservice.Microservice, cfg microservice.IConfig) IInventoryHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	invRepo := NewInventoryRepository(pst)
	invMqRepo := NewInventoryMQRepository(prod)
	invService := NewInventoryService(invRepo, invMqRepo)

	invOptRepo := NewInventoryOptionMainRepository(pst)
	invOptService := NewInventoryOptionMainService(invOptRepo)

	optGroupRepo := NewOptionGroupRepository(pst)
	optGroupService := NewOptionGroupService(optGroupRepo)

	categoryRepo := category.NewCategoryRepository(pst)
	inventoryCategoryService := NewInventorycategoryService(invRepo, categoryRepo, invMqRepo)

	return &InventoryHttp{
		ms:                       ms,
		cfg:                      cfg,
		invService:               invService,
		invOptService:            invOptService,
		optGroupService:          optGroupService,
		inventoryCategoryService: inventoryCategoryService,
	}
}

func (h InventoryHttp) RouteSetup() {
	h.ms.GET("/inventory/:id", h.InfoInventory)
	h.ms.GET("/inventory", h.SearchInventory)
	h.ms.POST("/inventory", h.CreateInventory)
	h.ms.POST("/inventory/bulk", h.CreateInBatchInventory)
	h.ms.PUT("/inventory/:id", h.UpdateInventory)
	h.ms.DELETE("/inventory/:id", h.DeleteInventory)
	h.ms.GET("/inventory/fetchupdate", h.LastActivityInventory)

	h.ms.GET("/option/:id", h.InfoInventoryOptionMain)
	h.ms.GET("/option", h.SearchInventoryOptionMain)
	h.ms.POST("/option", h.CreateInventoryOptionMain)
	h.ms.PUT("/option/:id", h.UpdateInventoryOptionMain)
	h.ms.DELETE("/option/:id", h.DeleteInventoryOptionMain)

	h.ms.GET("/optgroup/:id", h.InfoOptionGroup)
	h.ms.GET("/optgroup", h.SearchOptionGroup)
	h.ms.POST("/optgroup", h.CreateOptionGroup)
	h.ms.PUT("/optgroup/:id", h.UpdateOptionGroup)
	h.ms.DELETE("/optgroup/:id", h.DeleteOptionGroup)

	h.ms.POST("/inventory/categoryupdate/:catid", h.UpdateProductCategory)
}

// Create Inventory godoc
// @Description Create Inventory
// @Tags		Inventory
// @Param		Inventory  body      models.Inventory  true  "Inventory"
// @Accept 		json
// @Success		201	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /inventory [post]
func (h InventoryHttp) CreateInventory(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	inventoryReq := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	_, guidx, err := h.invService.CreateInventory(shopID, authUsername, *inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// docIdx := models.InventoryIndex{}
	// docIdx.ID = idx
	// docIdx.ShopID = userInfo.ShopID
	// docIdx.GuidFixed = guidx

	// err = h.invService.CreateIndex(docIdx)
	// if err != nil {
	// 	return err
	// }

	ctx.Response(
		http.StatusCreated,
		models.ApiResponse{
			Success: true,
			ID:      guidx,
		})

	return nil
}

// Create Inventory Bulk godoc
// @Description Create Inventory
// @Tags		Inventory
// @Param		Inventory  body      []models.Inventory  true  "Inventory"
// @Accept 		json
// @Success		201	{object}	models.InventoryBulkInsertResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /inventory/bulk [post]
func (h InventoryHttp) CreateInBatchInventory(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	inventoryReq := &[]models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	inventoryBulkResponse, err := h.invService.CreateInBatch(shopID, authUsername, *inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		models.InventoryBulkReponse{
			Success:             true,
			InventoryBulkImport: inventoryBulkResponse,
		},
	)

	return nil
}

// Update Inventory godoc
// @Description Update Inventory
// @Tags		Inventory
// @Param		id  path      string  true  "Inventory ID"
// @Param		Inventory  body      models.Inventory  true  "Inventory"
// @Accept 		json
// @Success		201	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /inventory/{id} [put]
func (h InventoryHttp) UpdateInventory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	inventoryReq := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.invService.UpdateInventory(shopID, id, authUsername, *inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		models.ApiResponse{
			Success: true,
			ID:      id,
		})

	return nil
}

// Delete Inventory godoc
// @Description Delete Inventory
// @Tags		Inventory
// @Param		id  path      string  true  "Inventory ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /inventory/{id} [delete]
func (h InventoryHttp) DeleteInventory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.invService.DeleteInventory(shopID, id, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success: true,
			ID:      id,
		},
	)
	return nil
}

// Get Inventory godoc
// @Description get struct array by ID
// @Tags		Inventory
// @Param		id  path      string  true  "Inventory ID"
// @Accept 		json
// @Success		200	{object}	models.InventoryInfoResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /inventory/{id} [get]
func (h InventoryHttp) InfoInventory(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.invService.InfoInventory(shopID, id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success: true,
			Data:    doc,
		},
	)

	return nil
}

// List Inventory godoc
// @Description get struct array by ID
// @Tags		Inventory
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}	models.InventoryPageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /inventory [get]
func (h InventoryHttp) SearchInventory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.invService.SearchInventory(shopID, q, page, limit)

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

func (h InventoryHttp) InfoMongoInventory(ctx microservice.IContext) error {

	id := ctx.Param("id")

	doc, err := h.invService.InfoMongoInventory(id)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success: true,
			Data:    doc,
		},
	)

	return nil
}

// Fetch Update Inventory By Date godoc
// @Description Fetch Update Inventory By Date
// @Tags		Inventory
// @Param		lastUpdate query string true "DateTime YYYY-MM-DDTHH:mm"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept		json
// @Success		200 {object} models.InventoryFetchUpdateResponse
// @Failure		401 {object} models.AuthResponseFailed
// @Security	AccessToken
// @Router		/inventory/fetchupdate [get]
func (h InventoryHttp) LastActivityInventory(ctx microservice.IContext) error {
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

	docList, pagination, err := h.invService.LastActivityInventory(shopID, lastUpdate, page, limit)

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

// Update Inventory Category List godoc
// @Description Update Inventory Category List
// @Tags		Inventory
// @Param		catid  path      string  true  "Category GUID"
// @Param		Inventory  body      []string  true  "Inventory"
// @Accept 		json
// @Success		201	{object}	models.ResponseSuccess
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /inventory/categoryupdate/{catid} [post]
func (h InventoryHttp) UpdateProductCategory(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	catid := ctx.Param("catid")

	input := ctx.ReadInput()

	var inventoryReq []string
	err := json.Unmarshal([]byte(input), &inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.inventoryCategoryService.UpdateInventoryCategoryBulk(shopID, authUsername, catid, inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		models.ApiResponse{
			Success: true,
		})

	return nil
}
