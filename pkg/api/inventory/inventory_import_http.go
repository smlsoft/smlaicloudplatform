package inventory

import "smlcloudplatform/internal/microservice"

type InventoryImportHttp struct {
	ms              *microservice.Microservice
	cfg             microservice.IConfig
	invService      IInventoryService
	invOptService   IInventoryOptionMainService
	optGroupService IOptionGroupService
}

type IInventoryImportHttp interface {
	RouteSetup()
}

func NewInventoryImportHttp(ms *microservice.Microservice, cfg microservice.IConfig) IInventoryImportHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	invRepo := NewInventoryRepository(pst)
	invService := NewInventoryService(invRepo, nil, nil)

	invOptRepo := NewInventoryOptionMainRepository(pst)
	invOptService := NewInventoryOptionMainService(invOptRepo)

	optGroupRepo := NewOptionGroupRepository(pst)
	optGroupService := NewOptionGroupService(optGroupRepo)

	return &InventoryImportHttp{
		ms:              ms,
		cfg:             cfg,
		invService:      invService,
		invOptService:   invOptService,
		optGroupService: optGroupService,
	}
}

func (h *InventoryImportHttp) RouteSetup() {

	h.ms.GET("/inventoryimport", h.ListInventoryImport)
	h.ms.POST("/inventoryimport", h.CreateInventoryImport)
	h.ms.DELETE("/inventoryimport/:id", h.DeleteInventoryImport)

	h.ms.GET("/optionimport", h.ListInventoryOptionImport)
	h.ms.POST("/optionimport", h.CreateInventoryOptionImport)
	h.ms.DELETE("/optionimport/:id", h.DeleteInventoryOptionImport)

	h.ms.GET("/categoryimport", h.ListCategoryImport)
	h.ms.POST("/categoryimport", h.CreateCategoryImport)
	h.ms.DELETE("/categoryimport/:id", h.DeleteCategoryImport)
}

// List Inventory Import godoc
// @Description get struct array by ID
// @Tags		Import
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}	models.InventoryPageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /inventoryimport [get]
func (h *InventoryImportHttp) ListInventoryImport(ctx microservice.IContext) error {
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
// @Router /inventoryimport [post]
func (h *InventoryImportHttp) CreateInventoryImport(ctx microservice.IContext) error {
	return nil
}

// Delete Inventory Import godoc
// @Description Delete Inventory
// @Tags		Import
// @Param		id  path      string  true  "Inventory Import ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /inventoryimport/{id} [delete]
func (h *InventoryImportHttp) DeleteInventoryImport(ctx microservice.IContext) error {
	return nil
}

func (h *InventoryImportHttp) ListInventoryOptionImport(ctx microservice.IContext) error {
	return nil
}

func (h *InventoryImportHttp) CreateInventoryOptionImport(ctx microservice.IContext) error {
	return nil
}

func (h *InventoryImportHttp) DeleteInventoryOptionImport(ctx microservice.IContext) error {
	return nil
}

func (h *InventoryImportHttp) ListCategoryImport(ctx microservice.IContext) error {
	return nil
}

func (h *InventoryImportHttp) CreateCategoryImport(ctx microservice.IContext) error {
	return nil
}

func (h *InventoryImportHttp) DeleteCategoryImport(ctx microservice.IContext) error {
	return nil
}
