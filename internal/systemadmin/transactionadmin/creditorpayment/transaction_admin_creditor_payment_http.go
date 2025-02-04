package creditorpayment

import (
	"encoding/json"
	"net/http"
	"smlaicloudplatform/internal/config"
	common "smlaicloudplatform/internal/models"
	adminModels "smlaicloudplatform/internal/systemadmin/models"
	"smlaicloudplatform/pkg/microservice"
)

type ICreditorPaymentTransactionAdminHttp interface {
	RegisterHttp(ms *microservice.Microservice, prefix string)
	ReSyncCreditorPaymentTransaction(ms microservice.IContext) error
}

type CreditorPaymentTransactionAdminHttp struct {
	svc ICreditorPaymentTransactionAdminService
}

func NewCreditorPaymentTransactionAdminHttp(ms *microservice.Microservice, cfg config.IConfig) ICreditorPaymentTransactionAdminHttp {

	producer := microservice.NewProducer(cfg.MQConfig().URI(), ms.Logger)
	mongoPersister := microservice.NewPersisterMongo(cfg.MongoPersisterConfig())

	svc := NewCreditorPaymentTransactionAdminService(mongoPersister, producer)

	return &CreditorPaymentTransactionAdminHttp{
		svc: svc,
	}
}

func (s *CreditorPaymentTransactionAdminHttp) RegisterHttp(ms *microservice.Microservice, prefix string) {
	ms.POST(prefix+"/transactionadmin/creditorpayment/resynctransaction", s.ReSyncCreditorPaymentTransaction)
}

func (s *CreditorPaymentTransactionAdminHttp) ReSyncCreditorPaymentTransaction(ctx microservice.IContext) error {

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

	err = s.svc.ReSyncCreditorPaymentDoc(req.ShopID)
	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	return nil
}
