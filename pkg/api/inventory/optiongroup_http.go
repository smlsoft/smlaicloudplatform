package inventory

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

func (h *InventoryHttp) CreateOptionGroup(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopId := ctx.UserInfo().ShopId
	input := ctx.ReadInput()

	docReq := &models.InventoryOptionGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.optGroupService.CreateOptionGroup(shopId, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		Id:      idx,
	})
	return nil
}

func (h *InventoryHttp) UpdateOptionGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopId := userInfo.ShopId

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.InventoryOptionGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.optGroupService.UpdateOptionGroup(id, shopId, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		Id:      id,
	})

	return nil
}

func (h *InventoryHttp) DeleteOptionGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopId := userInfo.ShopId

	id := ctx.Param("id")

	err := h.optGroupService.DeleteOptionGroup(id, shopId)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		Id:      id,
	})

	return nil
}

func (h *InventoryHttp) InfoOptionGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopId := userInfo.ShopId

	id := ctx.Param("id")

	doc, err := h.optGroupService.InfoOptionGroup(id, shopId)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

func (h *InventoryHttp) SearchOptionGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopId := userInfo.ShopId

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}
	docList, pagination, err := h.optGroupService.SearchOptionGroup(shopId, q, page, limit)

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
