package optiongroup

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/optiongroup/models"
	"smlcloudplatform/pkg/utils"
)

type OptionGroupHttp struct {
	ms              *microservice.Microservice
	cfg             microservice.IConfig
	optGroupService IOptionGroupService
}

func NewOptionGroupHttp(ms *microservice.Microservice, cfg microservice.IConfig) *OptionGroupHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	optGroupRepo := NewOptionGroupRepository(pst)
	optGroupService := NewOptionGroupService(optGroupRepo)

	return &OptionGroupHttp{
		ms:              ms,
		cfg:             cfg,
		optGroupService: optGroupService,
	}
}

func (h OptionGroupHttp) RouteSetup() {

	h.ms.GET("/optgroup/:id", h.InfoOptionGroup)
	h.ms.GET("/optgroup", h.SearchOptionGroup)
	h.ms.POST("/optgroup", h.CreateOptionGroup)
	h.ms.PUT("/optgroup/:id", h.UpdateOptionGroup)
	h.ms.DELETE("/optgroup/:id", h.DeleteOptionGroup)

}

// Create Inventory Option Group godoc
// @Description Create Inventory Option Group
// @Tags		Inventory
// @Param		Option  body      models.InventoryOptionGroup  true  "Option Group"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /optgroup [post]
func (h *OptionGroupHttp) CreateOptionGroup(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.InventoryOptionGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.optGroupService.CreateOptionGroup(shopID, authUsername, *docReq)

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

// Update Option Group godoc
// @Description Update Option Group
// @Tags		Inventory
// @Param		id  path      string  true  "Option ID"
// @Param		OptionGroup  body      models.InventoryOptionGroup  true  "Option Group"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /optgroup/{id} [put]
func (h *OptionGroupHttp) UpdateOptionGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.InventoryOptionGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.optGroupService.UpdateOptionGroup(shopID, id, authUsername, *docReq)

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

// Delete OptionGroup godoc
// @Description Delete OptionGroup
// @Tags		Inventory
// @Param		id  path      string  true  "OptionGroup ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /optgroup/{id} [delete]
func (h *OptionGroupHttp) DeleteOptionGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.optGroupService.DeleteOptionGroup(shopID, id, authUsername)

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

// Get Option Group Infomation godoc
// @Description Get Option Group Information
// @Tags		Inventory
// @Param		id  path      string  true  "OptionGroup Id"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /optgroup/{id} [get]
func (h *OptionGroupHttp) InfoOptionGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.optGroupService.InfoOptionGroup(shopID, id)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// List Option Group godoc
// @Description List Inventory Option Group
// @Tags		Inventory
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /optgroup [get]
func (h *OptionGroupHttp) SearchOptionGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)
	docList, pagination, err := h.optGroupService.SearchOptionGroup(shopID, pageable)

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
