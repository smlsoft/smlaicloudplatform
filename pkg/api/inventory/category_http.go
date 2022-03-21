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
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	categoryReq := &models.Category{}
	err := json.Unmarshal([]byte(input), &categoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.cateService.CreateCategory(shopID, authUsername, *categoryReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      idx,
	})
	return nil
}

func (h *InventoryHttp) UpdateCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	categoryReq := &models.Category{}
	err := json.Unmarshal([]byte(input), &categoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.cateService.UpdateCategory(id, shopID, authUsername, *categoryReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

func (h *InventoryHttp) DeleteCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.cateService.DeleteCategory(id, shopID)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

func (h *InventoryHttp) InfoCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.cateService.InfoCategory(id, shopID)

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
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}
	docList, pagination, err := h.cateService.SearchCategory(shopID, q, page, limit)

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
