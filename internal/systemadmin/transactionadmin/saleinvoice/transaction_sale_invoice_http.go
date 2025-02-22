package saleinvoice

import (
	"encoding/json"
	"net/http"
	"smlaicloudplatform/internal/config"
	common "smlaicloudplatform/internal/models"
	adminModels "smlaicloudplatform/internal/systemadmin/models"
	"smlaicloudplatform/pkg/microservice"
)

type ISaleInvoiceTransactionAdminHttp interface {
	ReSyncSaleInvoiceTransaction(ctx microservice.IContext) error
	ReSyncSaleInvoiceDeleteTransaction(ctx microservice.IContext) error
	ReSyncSaleInvoiceTransactionByDate(ctx microservice.IContext) error
	ReSyncSaleInvoiceDeleteTransactionByDate(ctx microservice.IContext) error
	RegisterHttp(ms *microservice.Microservice, prefix string)
}

type SaleInvoiceTransactionAdminHttp struct {
	svc ISaleInvoiceTransactionAdminService
}

func NewSaleInvoiceTransactionAdminHttp(ms *microservice.Microservice, cfg config.IConfig) ISaleInvoiceTransactionAdminHttp {

	producer := microservice.NewProducer(cfg.MQConfig().URI(), ms.Logger)
	mongoPersister := microservice.NewPersisterMongo(cfg.MongoPersisterConfig())

	svc := NewSaleInvoiceTransactionAdminService(mongoPersister, producer)
	return &SaleInvoiceTransactionAdminHttp{
		svc: svc,
	}
}

func (s *SaleInvoiceTransactionAdminHttp) RegisterHttp(ms *microservice.Microservice, prefix string) {
	ms.POST(prefix+"/transactionadmin/saleinvoice/resynctransaction", s.ReSyncSaleInvoiceTransaction)
	ms.POST(prefix+"/transactionadmin/saleinvoice/resynctransactionbydate", s.ReSyncSaleInvoiceTransactionByDate)
	ms.POST(prefix+"/transactionadmin/saleinvoice/resyncdeletetransaction", s.ReSyncSaleInvoiceDeleteTransaction)
	ms.POST(prefix+"/transactionadmin/saleinvoice/resyncdeletetransactionbydate", s.ReSyncSaleInvoiceDeleteTransactionByDate)
}

func (s *SaleInvoiceTransactionAdminHttp) ReSyncSaleInvoiceTransaction(ctx microservice.IContext) error {

	input := ctx.ReadInput()
	var req adminModels.RequestReSyncTenant

	err := json.Unmarshal([]byte(input), &req)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	err = s.svc.ReSyncSaleInvoiceDoc(req.ShopID)
	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, common.ResponseSuccess{
		Success: true,
	})
	return nil
}

func (s *SaleInvoiceTransactionAdminHttp) ReSyncSaleInvoiceDeleteTransaction(ctx microservice.IContext) error {

	input := ctx.ReadInput()
	var req adminModels.RequestReSyncTenant

	err := json.Unmarshal([]byte(input), &req)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	err = s.svc.ReSyncSaleInvoiceDeleteDoc(req.ShopID)
	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, common.ResponseSuccess{
		Success: true,
	})
	return nil
}

func (s *SaleInvoiceTransactionAdminHttp) ReSyncSaleInvoiceTransactionByDate(ctx microservice.IContext) error {

	input := ctx.ReadInput()
	var req adminModels.RequestReSyncTenantByDate

	err := json.Unmarshal([]byte(input), &req)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	err = s.svc.ReSyncSaleInvoiceDocByDate(req.ShopID, req.Date)
	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, common.ResponseSuccess{
		Success: true,
	})
	return nil
}

func (s *SaleInvoiceTransactionAdminHttp) ReSyncSaleInvoiceDeleteTransactionByDate(ctx microservice.IContext) error {

	input := ctx.ReadInput()
	var req adminModels.RequestReSyncTenant

	err := json.Unmarshal([]byte(input), &req)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	err = s.svc.ReSyncSaleInvoiceDeleteDoc(req.ShopID)
	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, common.ResponseSuccess{
		Success: true,
	})
	return nil
}
