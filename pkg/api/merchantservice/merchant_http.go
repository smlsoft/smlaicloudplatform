package merchantservice

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type IMerchantHttp interface {
	CreateMerchant(ctx microservice.IServiceContext) error
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
		})
		return err
	}

	ctx.Response(http.StatusOK, &models.ApiResponse{
		Success: true,
		Id:      idx,
	})

	return nil
}
