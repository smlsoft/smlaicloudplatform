package inventory

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

func (h *InventoryHttp) CreateCategory(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopId := ctx.UserInfo().ShopId
	input := ctx.ReadInput()

	categoryReq := &models.Category{}
	err := json.Unmarshal([]byte(input), &categoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.cateService.CreateCategory(shopId, authUsername, *categoryReq)

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

func (h *InventoryHttp) UpdateCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopId := userInfo.ShopId

	id := ctx.Param("id")
	input := ctx.ReadInput()

	categoryReq := &models.Category{}
	err := json.Unmarshal([]byte(input), &categoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.cateService.UpdateCategory(id, shopId, authUsername, *categoryReq)

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

func (h *InventoryHttp) DeleteCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopId := userInfo.ShopId

	id := ctx.Param("id")

	err := h.cateService.DeleteCategory(id, shopId)

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

func (h *InventoryHttp) InfoCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopId := userInfo.ShopId

	id := ctx.Param("id")

	doc, err := h.cateService.InfoCategory(id, shopId)

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

func (h *InventoryHttp) SearchCategory(ctx microservice.IContext) error {
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
	docList, pagination, err := h.cateService.SearchCategory(shopId, q, page, limit)

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
