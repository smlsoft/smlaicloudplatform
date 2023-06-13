package api

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	common "smlcloudplatform/pkg/models"
	shopModel "smlcloudplatform/pkg/shop/models"
	journalModels "smlcloudplatform/pkg/vfgl/journal/models"
)

type IMigrationAPI interface {
	ImportJournal(ctx microservice.IContext) error
	ImportShop(ctx microservice.IContext) error
	RegisterHttp(ms *microservice.Microservice)
}

type MigrationAPI struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc IMigrationService
}

func NewMigrationAPI(ms *microservice.Microservice, cfg config.IConfig) *MigrationAPI {

	mqPersister := ms.Producer(cfg.MQConfig())
	svc := NewMigrationService(ms.MongoPersister(cfg.MongoPersisterConfig()), mqPersister)
	return &MigrationAPI{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (m *MigrationAPI) RegisterHttp(ms *microservice.Microservice) {
	ms.POST("/migrationtools/journalimport", m.ImportJournal)
	ms.POST("/migrationtools/shopimport", m.ImportShop)
}

func (m *MigrationAPI) ImportJournal(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	docReq := &[]journalModels.JournalDoc{}
	err := json.Unmarshal([]byte(input), &docReq)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	err = m.svc.ImportJournal(*docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ResponseSuccess{
		Success: true,
	})
	// ctx.Response(http.StatusOK, common.ApiResponse{
	// 	Success: true,
	// 	Data:    docReq,
	// })
	return nil
}

func (m *MigrationAPI) ImportShop(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	docReq := &[]shopModel.ShopDoc{}
	err := json.Unmarshal([]byte(input), &docReq)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	err = m.svc.ImportShop(*docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ResponseSuccess{
		Success: true,
	})

	// ctx.Response(http.StatusOK, common.ApiResponse{
	// 	Success: true,
	// 	Data:    docReq,
	// })
	return nil
}
