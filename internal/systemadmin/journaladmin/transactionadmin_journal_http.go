package journal

import (
	"encoding/json"
	"net/http"
	"smlaicloudplatform/internal/config"
	common "smlaicloudplatform/internal/models"
	adminModels "smlaicloudplatform/internal/systemadmin/models"
	"smlaicloudplatform/pkg/microservice"
)

type IJournalTransactionAdminHttp interface {
	RegisterHttp(ms *microservice.Microservice, prefix string)
	ReSyncJournalTransaction(ms microservice.IContext) error
	ReSyncJournalDeleteTransaction(ms microservice.IContext) error
}

type JournalTransactionAdminHttp struct {
	svc IJournalTransactionAdminService
}

func NewJournalTransactionAdminHttp(ms *microservice.Microservice, cfg config.IConfig) IJournalTransactionAdminHttp {

	producer := microservice.NewProducer(cfg.MQConfig().URI(), ms.Logger)
	mongoPersister := microservice.NewPersisterMongo(cfg.MongoPersisterConfig())

	svc := NewJournalTransactionAdminService(mongoPersister, producer)

	return &JournalTransactionAdminHttp{
		svc: svc,
	}
}

func (s *JournalTransactionAdminHttp) RegisterHttp(ms *microservice.Microservice, prefix string) {
	ms.POST(prefix+"/transactionadmin/journal/resynctransaction", s.ReSyncJournalTransaction)
	ms.POST(prefix+"/transactionadmin/journal/resyncdeletetransaction", s.ReSyncJournalDeleteTransaction)
}

func (s *JournalTransactionAdminHttp) ReSyncJournalTransaction(ctx microservice.IContext) error {

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

	err = s.svc.ReSyncJournalTransactionDoc(req.ShopID)
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

func (h *JournalTransactionAdminHttp) ReSyncJournalDeleteTransaction(ctx microservice.IContext) error {

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

	err = h.svc.ReSyncJournalDeleteTransactionDoc(req.ShopID)
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
