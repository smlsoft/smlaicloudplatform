package inventory

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

func (h *InventoryHttp) CreateOptionGroup(ctx microservice.IServiceContext) error {
	authUsername := ctx.UserInfo().Username
	merchantId := ctx.UserInfo().MerchantId
	input := ctx.ReadInput()

	docReq := &models.InventoryOptionGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.optGroupService.CreateOptionGroup(merchantId, authUsername, *docReq)

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

func (h *InventoryHttp) UpdateOptionGroup(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.InventoryOptionGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.optGroupService.UpdateOptionGroup(id, merchantId, authUsername, *docReq)

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

func (h *InventoryHttp) DeleteOptionGroup(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	err := h.optGroupService.DeleteOptionGroup(id, merchantId)

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

func (h *InventoryHttp) InfoOptionGroup(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	doc, err := h.optGroupService.InfoOptionGroup(id, merchantId)

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

func (h *InventoryHttp) SearchOptionGroup(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	merchantId := userInfo.MerchantId

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}
	docList, pagination, err := h.optGroupService.SearchOptionGroup(merchantId, q, page, limit)

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
