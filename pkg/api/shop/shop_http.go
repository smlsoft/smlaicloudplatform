package shop

import (
	"encoding/json"
	"errors"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IShopHttp interface {
	RouteSetup()
	CreateShop(ctx microservice.IContext) error
	UpdateShop(ctx microservice.IContext) error
	DeleteShop(ctx microservice.IContext) error
	InfoShop(ctx microservice.IContext) error
	SearchShop(ctx microservice.IContext) error
}

type ShopHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service IShopService
}

func NewShopHttp(ms *microservice.Microservice, cfg microservice.IConfig) IShopHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	repo := NewShopRepository(pst)
	shopUserRepo := NewShopUserRepository(pst)
	service := NewShopService(repo, shopUserRepo)

	return &ShopHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h *ShopHttp) RouteSetup() {
	h.ms.GET("/shop/:id", h.InfoShop)
	h.ms.GET("/shop", h.SearchShop)

	h.ms.POST("/shop", h.CreateShop)
	h.ms.PUT("/shop/:id", h.UpdateShop)
	h.ms.DELETE("/shop/:id", h.DeleteShop)
}

// Create Shop godoc
// @Description Create Shop
// @Tags		Shop
// @Accept 		json
// @Param		Shop  body      models.Shop  true  "Add Shop"
// @Success		200	{array}	models.Shop
// @Failure		401 {object}	models.ResponseSuccessWithId
// @Security     AccessToken
// @Router /shop [post]
func (h *ShopHttp) CreateShop(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	if len(authUsername) < 1 {
		ctx.ResponseError(400, "user authentication invalid")
	}

	input := ctx.ReadInput()

	shopReq := &models.Shop{}
	err := json.Unmarshal([]byte(input), &shopReq)

	if err != nil {
		ctx.ResponseError(400, "shop payload invalid")
		return err
	}

	idx, err := h.service.CreateShop(authUsername, *shopReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, &models.ApiResponse{
		Success: true,
		Id:      idx,
	})

	return nil
}

func (h *ShopHttp) UpdateShop(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	id := ctx.Param("id")
	input := ctx.ReadInput()

	shopRequest := &models.Shop{}
	err := json.Unmarshal([]byte(input), &shopRequest)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if userInfo.Role == "" || userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &models.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	err = h.service.UpdateShop(id, authUsername, *shopRequest)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, &models.ApiResponse{
		Success: true,
		Id:      id,
	})
	return nil
}

func (h *ShopHttp) DeleteShop(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username

	id := ctx.Param("id")

	if userInfo.Role == "" || userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &models.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	err := h.service.DeleteShop(id, authUsername)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}
	ctx.Response(http.StatusOK, &models.ApiResponse{
		Success: true,
		Id:      id,
	})
	return nil
}

func (h *ShopHttp) InfoShop(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	id := ctx.Param("id")

	if userInfo.Role == "" || userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &models.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	shopInfo, err := h.service.InfoShop(id)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, &models.ApiResponse{
		Success: true,
		Data:    shopInfo,
	})
	return nil
}

// List Shop godoc
// @Description Access to Shop
// @Tags		Shop
// @Accept 		json
// @Success		200	{array}	models.ShopInfo
// @Failure		401 {object}	models.ApiResponse
// @Security     AccessToken
// @Router /shop [get]
func (h *ShopHttp) SearchShop(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()

	if userInfo.Role == "" || userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &models.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	shopList, pagination, err := h.service.SearchShop(q, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": shopList})
	return nil
}
