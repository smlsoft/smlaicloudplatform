package smspaymentsettings

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	smspatternsrepo "smlcloudplatform/pkg/smsreceive/smspatterns/repositories"
	"smlcloudplatform/pkg/smsreceive/smspaymentsettings/models"
	"smlcloudplatform/pkg/smsreceive/smspaymentsettings/repositories"
	"smlcloudplatform/pkg/smsreceive/smspaymentsettings/services"
)

type ISmsPaymentSettingsHttp interface{}

type SmsPaymentSettingsHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.ISmsPaymentSettingsHttpService
}

func NewSmsPaymentSettingsHttp(ms *microservice.Microservice, cfg microservice.IConfig) SmsPaymentSettingsHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := repositories.NewSmsPaymentSettingsRepository(pst)

	smspatternRepo := smspatternsrepo.NewSmsPatternsRepository(pst)

	svc := services.NewSmsPaymentSettingsHttpService(repo, smspatternRepo)

	return SmsPaymentSettingsHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h SmsPaymentSettingsHttp) RouteSetup() {

	h.ms.POST("/smspaymentsettings", h.CreateSmsPaymentSettings)
	h.ms.GET("/smspaymentsettings", h.InfoSmsPaymentSettings)
}

func (h SmsPaymentSettingsHttp) CreateSmsPaymentSettings(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.SmsPaymentSettings{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.SaveSmsPaymentSettings(shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

func (h SmsPaymentSettingsHttp) InfoSmsPaymentSettings(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	doc, err := h.svc.InfoSmsPaymentSettings(shopID)

	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v", err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}
