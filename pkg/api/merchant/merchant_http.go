package merchant

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IMerchantHttp interface {
	RouteSetup()
	CreateMerchant(ctx microservice.IContext) error
	UpdateMerchant(ctx microservice.IContext) error
	DeleteMerchant(ctx microservice.IContext) error
	InfoMerchant(ctx microservice.IContext) error
	SearchMerchant(ctx microservice.IContext) error
}

type MerchantHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service IMerchantService
}

func NewMerchantHttp(ms *microservice.Microservice, cfg microservice.IConfig) IMerchantHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	repo := NewMerchantRepository(pst)
	merchantUserRepo := NewMerchantUserRepository(pst)
	service := NewMerchantService(repo, merchantUserRepo)

	return &MerchantHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h *MerchantHttp) RouteSetup() {
	h.ms.GET("/merchant/:id", h.InfoMerchant)
	h.ms.GET("/merchant", h.SearchMerchant)

	h.ms.POST("/merchant", h.CreateMerchant)
	h.ms.PUT("/merchant/:id", h.UpdateMerchant)
	h.ms.DELETE("/merchant/:id", h.DeleteMerchant)
}

func (h *MerchantHttp) CreateMerchant(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	if len(authUsername) < 1 {
		ctx.ResponseError(400, "user authentication invalid")
	}

	input := ctx.ReadInput()

	merchantReq := &models.Merchant{}
	err := json.Unmarshal([]byte(input), &merchantReq)

	if err != nil {
		ctx.ResponseError(400, "merchant payload invalid")
		return err
	}

	idx, err := h.service.CreateMerchant(authUsername, *merchantReq)

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

func (h *MerchantHttp) UpdateMerchant(ctx microservice.IContext) error {

	authUsername := ctx.UserInfo().Username
	id := ctx.Param("id")
	input := ctx.ReadInput()

	merchantRequest := &models.Merchant{}
	err := json.Unmarshal([]byte(input), &merchantRequest)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateMerchant(id, authUsername, *merchantRequest)

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

func (h *MerchantHttp) DeleteMerchant(ctx microservice.IContext) error {

	authUsername := ctx.UserInfo().Username
	id := ctx.Param("id")

	err := h.service.DeleteMerchant(id, authUsername)

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

func (h *MerchantHttp) InfoMerchant(ctx microservice.IContext) error {

	authUsername := ctx.UserInfo().Username
	id := ctx.Param("id")

	merchantInfo, err := h.service.InfoMerchant(id, authUsername)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}
	ctx.Response(http.StatusOK, &models.ApiResponse{
		Success: true,
		Data:    merchantInfo,
	})
	return nil
}

func (h *MerchantHttp) SearchMerchant(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	merchantList, pagination, err := h.service.SearchMerchant(authUsername, q, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "pagination": pagination, "data": merchantList})
	return nil
}
