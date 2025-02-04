package stockreturnproduct

import (
	"encoding/json"
	"net/http"
	"smlaicloudplatform/internal/config"
	common "smlaicloudplatform/internal/models"
	adminModels "smlaicloudplatform/internal/systemadmin/models"
	"smlaicloudplatform/pkg/microservice"
)

type IStockReturnProductTransactionAdminHttp interface {
	ReSyncStockReturnTransaction(ctx microservice.IContext) error
	ReSyncStockReturnDeleteTransaction(ctx microservice.IContext) error
	RegisterHttp(ms *microservice.Microservice, prefix string)
}

type StockReturnProductTransactionAdminHttp struct {
	svc IStockReturnProductTransactionAdminService
}

func NewStockReturnProductTransactionAdminHttp(ms *microservice.Microservice, cfg config.IConfig) IStockReturnProductTransactionAdminHttp {

	producer := microservice.NewProducer(cfg.MQConfig().URI(), ms.Logger)
	mongoPersister := microservice.NewPersisterMongo(cfg.MongoPersisterConfig())

	svc := NewStockReturnProductTransactionAdminService(mongoPersister, producer)

	return &StockReturnProductTransactionAdminHttp{
		svc: svc,
	}
}

func (s *StockReturnProductTransactionAdminHttp) RegisterHttp(ms *microservice.Microservice, prefix string) {
	ms.POST(prefix+"/transactionadmin/stockreturnproduct/resynctransaction", s.ReSyncStockReturnTransaction)
	ms.POST(prefix+"/transactionadmin/stockreturnproduct/resyncdeletetransaction", s.ReSyncStockReturnDeleteTransaction)
}

func (s *StockReturnProductTransactionAdminHttp) ReSyncStockReturnTransaction(ctx microservice.IContext) error {

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

	err = s.svc.ReSyncStockReturnProductDoc(req.ShopID)
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

func (s *StockReturnProductTransactionAdminHttp) ReSyncStockReturnDeleteTransaction(ctx microservice.IContext) error {
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

	err = s.svc.ReSyncStockReturnProductDeleteDoc(req.ShopID)
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
