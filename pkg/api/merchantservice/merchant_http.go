package merchantservice

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type IMerchantHttp interface {
	CreateMerchant(ctx microservice.IServiceContext) error
	UpdateMerchant(ctx microservice.IServiceContext) error
	DeleteMerchant(ctx microservice.IServiceContext) error
	InfoMerchant(ctx microservice.IServiceContext) error
}

type MerchantHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service IMerchantService
}

func NewMerchantHttp(ms *microservice.Microservice, cfg microservice.IConfig) IMerchantHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	repo := NewMerchantRepository(pst)
	service := NewMerchantService(repo)

	return &MerchantHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h *MerchantHttp) CreateMerchant(ctx microservice.IServiceContext) error {
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

func (h *MerchantHttp) UpdateMerchant(ctx microservice.IServiceContext) error {

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

func (h *MerchantHttp) DeleteMerchant(ctx microservice.IServiceContext) error {

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

func (h *MerchantHttp) InfoMerchant(ctx microservice.IServiceContext) error {

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
