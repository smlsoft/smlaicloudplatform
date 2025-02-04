package purchasereturn

import (
	"encoding/json"
	"net/http"
	"smlaicloudplatform/internal/config"
	common "smlaicloudplatform/internal/models"
	adminModels "smlaicloudplatform/internal/systemadmin/models"
	"smlaicloudplatform/pkg/microservice"
)

type IPurchaseReturnTransactionAdminHttp interface {
	ResyncPurchaseReturnDoc(ctx microservice.IContext) error
	ResyncPurchaseReturnDeleteDoc(ctx microservice.IContext) error
	RegisterHttp(ms *microservice.Microservice, prefix string)
}

type PurchaseReturnTransactionAdminHttp struct {
	svc IPurchaseReturnTransactionAdminService
}

func NewPurchaseReturnTransactionAdminHttp(ms *microservice.Microservice, cfg config.IConfig) IPurchaseReturnTransactionAdminHttp {

	producer := microservice.NewProducer(cfg.MQConfig().URI(), ms.Logger)
	mongoPersister := microservice.NewPersisterMongo(cfg.MongoPersisterConfig())

	svc := NewPurchaseReturnTransactionAdminService(mongoPersister, producer)
	return &PurchaseReturnTransactionAdminHttp{
		svc: svc,
	}
}

func (p *PurchaseReturnTransactionAdminHttp) RegisterHttp(ms *microservice.Microservice, prefix string) {
	ms.POST(prefix+"/transactionadmin/purchasereturn/resyncpurchasereturndoc", p.ResyncPurchaseReturnDoc)
	ms.POST(prefix+"/transactionadmin/purchasereturn/resyncpurchasereturndeletedoc", p.ResyncPurchaseReturnDeleteDoc)
}

func (h *PurchaseReturnTransactionAdminHttp) ResyncPurchaseReturnDoc(ctx microservice.IContext) error {

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

	err = h.svc.ResyncPurchaseReturnDoc(req.ShopID)
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

func (h *PurchaseReturnTransactionAdminHttp) ResyncPurchaseReturnDeleteDoc(ctx microservice.IContext) error {

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

	err = h.svc.ResyncPurchaseReturnDeleteDoc(req.ShopID)
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
