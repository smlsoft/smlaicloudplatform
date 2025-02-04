package debtoradmin

import (
	"encoding/json"
	"net/http"
	"smlaicloudplatform/internal/config"
	common "smlaicloudplatform/internal/models"
	adminModels "smlaicloudplatform/internal/systemadmin/models"
	"smlaicloudplatform/pkg/microservice"
)

type IDebtorAdminHttp interface {
	ReSyncDebtor(ctx microservice.IContext) error
	ReCalcDebtorBalance(ctx microservice.IContext) error
	RegisterHttp(ms *microservice.Microservice, prefix string)
}

type DebtorAdminHttp struct {
	svc IDebtorAdminService
}

func NewDebtorAdminHttp(ms *microservice.Microservice, cfg config.IConfig) IDebtorAdminHttp {

	producer := microservice.NewProducer(cfg.MQConfig().URI(), ms.Logger)
	mongoPersister := microservice.NewPersisterMongo(cfg.MongoPersisterConfig())

	svc := NewDebtorAdminService(mongoPersister, producer)

	return &DebtorAdminHttp{
		svc: svc,
	}
}

func (s *DebtorAdminHttp) RegisterHttp(ms *microservice.Microservice, prefix string) {
	ms.POST(prefix+"/debtoradmin/resyncdebtor", s.ReSyncDebtor)
	ms.POST(prefix+"/debtoradmin/recalcbalance", s.ReCalcDebtorBalance)
}

func (s *DebtorAdminHttp) ReSyncDebtor(ctx microservice.IContext) error {

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

	err = s.svc.ReSyncDebtor(req.ShopID)
	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	return nil
}

func (s *DebtorAdminHttp) ReCalcDebtorBalance(ctx microservice.IContext) error {

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

	err = s.svc.ReCalcDebtorBalance(req.ShopID)
	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	return nil
}
