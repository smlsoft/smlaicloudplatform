package smstransaction

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/smstransaction/models"
	"smlcloudplatform/pkg/smstransaction/repositories"
	"smlcloudplatform/pkg/smstransaction/services"
	"smlcloudplatform/pkg/utils"
)

type ISmsTransactionHttp interface{}

type SmsTransactionHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.ISmsTransactionHttpService
}

func NewSmsTransactionHttp(ms *microservice.Microservice, cfg microservice.IConfig) SmsTransactionHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := repositories.NewSmsTransactionRepository(pst)
	svc := services.NewSmsTransactionHttpService(repo)

	return SmsTransactionHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h SmsTransactionHttp) RouteSetup() {
	h.ms.GET("/smstransaction", h.SearchSmsTransaction)
	h.ms.POST("/smstransaction", h.CreateSmsTransaction)
	h.ms.GET("/smstransaction/:id", h.InfoSmsTransaction)
	h.ms.PUT("/smstransaction/:id", h.UpdateSmsTransaction)
	h.ms.DELETE("/smstransaction/:id", h.DeleteSmsTransaction)
}

// Create SMS Transaction godoc
// @Summary		รับข้อมูล sms
// @Description รับข้อมูล sms
// @Tags		SMS
// @Param		SMS Transaction  body      models.SmsTransaction  true  "sms data"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /smstransaction [post]
func (h SmsTransactionHttp) CreateSmsTransaction(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.SmsTransaction{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateSmsTransaction(shopID, authUsername, *docReq)

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

// Update SMS Transaction godoc
// @Summary		รับข้อมูล sms
// @Description รับข้อมูล sms
// @Tags		SMS
// @Param		id  path      string  true  "GIUDFIXED"
// @Param		Journal  body      models.SmsTransaction  true  "sms data"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /smstransaction/{id} [put]
func (h SmsTransactionHttp) UpdateSmsTransaction(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.SmsTransaction{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateSmsTransaction(id, shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

// Delete SMS Transaction godoc
// @Summary		รับข้อมูล sms
// @Description รับข้อมูล sms
// @Tags		SMS
// @Param		id  path      string  true  "GIUDFIXED"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /smstransaction/{id} [delete]
func (h SmsTransactionHttp) DeleteSmsTransaction(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteSmsTransaction(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

// GET SMS Transaction godoc
// @Summary		รับข้อมูล sms
// @Description รับข้อมูล sms
// @Tags		SMS
// @Param		id  path      string  true  "GIUDFIXED"
// @Accept 		json
// @Success		200	{object}	models.SmsTransactionInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /smstransaction/{id} [get]
func (h SmsTransactionHttp) InfoSmsTransaction(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get SmsTransaction %v", id)
	doc, err := h.svc.InfoSmsTransaction(id, shopID)

	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v: %v", id, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// LIST SMS Transaction godoc
// @Summary		รับข้อมูล sms
// @Description รับข้อมูล sms
// @Tags		SMS
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.SmsTransactionPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/journal [get]
func (h SmsTransactionHttp) SearchSmsTransaction(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchSmsTransaction(shopID, q, page, limit, sort)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Data:       docList,
		Pagination: pagination,
	})
	return nil
}
