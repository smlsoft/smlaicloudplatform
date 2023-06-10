package smltransaction

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/transaction/smltransaction/models"
	"smlcloudplatform/pkg/transaction/smltransaction/repositories"
	"smlcloudplatform/pkg/transaction/smltransaction/services"
)

type ISMLTransactionHttp interface{}

type SMLTransactionHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.ISMLTransactionHttpService
}

func NewSMLTransactionHttp(ms *microservice.Microservice, cfg config.IConfig) SMLTransactionHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	repo := repositories.NewSMLTransactionRepository(pst)

	mqRepo := repositories.NewSMLTransactionMessageQueueRepository(prod)

	svc := services.NewSMLTransactionHttpService(repo, mqRepo)

	return SMLTransactionHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h SMLTransactionHttp) RouteSetup() {

	h.ms.POST("/sml-transaction", h.CreateSMLTransaction)
	h.ms.POST("/sml-transaction/bulk", h.BulkCreateSMLTransaction)
	h.ms.DELETE("/sml-transaction", h.DeleteSMLTransaction)
}

type Data struct {
	DocNo       string                 `json:"docno"`
	DynamicData map[string]interface{} `json:",inline"`
}

// Create SMLTransaction godoc
// @Description Create SMLTransaction
// @Tags		SMLTransaction
// @Param		SMLTransactionRequest  body      models.SMLTransactionRequest  true  "SMLTransactionRequest"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /sml-transaction [post]
func (h SMLTransactionHttp) CreateSMLTransaction(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := models.SMLTransactionRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	_, ok := docReq.Body[docReq.KeyID]
	if !ok {
		ctx.ResponseError(400, "key field not found in body")
		return err
	}

	idx, err := h.svc.CreateSMLTransaction(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})
	return nil
}

// Bulk Create SMLTransaction godoc
// @Description Bulk Create SMLTransaction
// @Tags		SMLTransaction
// @Param		SMLTransactionBulkRequest  body      models.SMLTransactionBulkRequest  true  "SMLTransaction Bulk Request"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /sml-transaction/bulk [post]
func (h SMLTransactionHttp) BulkCreateSMLTransaction(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := models.SMLTransactionBulkRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	_, err = h.svc.SaveInBatch(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete SMLTransaction godoc
// @Description Delete SMLTransaction
// @Tags		SMLTransaction
// @Param		SMLTransactionKeyRequest  body      models.SMLTransactionKeyRequest  true  "SMLTransaction Key Request"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /sml-transaction [delete]
func (h SMLTransactionHttp) DeleteSMLTransaction(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	input := ctx.ReadInput()

	docReq := models.SMLTransactionKeyRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	keys, err := h.svc.DeleteSMLTransaction(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    keys,
	})

	return nil
}
