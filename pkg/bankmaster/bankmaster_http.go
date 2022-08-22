package bankmaster

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/bankmaster/models"
	"smlcloudplatform/pkg/bankmaster/repositories"
	"smlcloudplatform/pkg/bankmaster/services"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
)

type IBankMasterHttp interface{}

type BankMasterHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IBankMasterHttpService
}

func NewBankMasterHttp(ms *microservice.Microservice, cfg microservice.IConfig) BankMasterHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := repositories.NewBankMasterRepository(pst)
	svc := services.NewBankMasterHttpService(repo)

	return BankMasterHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h BankMasterHttp) RouteSetup() {

	h.ms.POST("/bankmaster/bulk", h.SaveBulk)

	h.ms.GET("/bankmaster", h.SearchBankMaster)
	h.ms.POST("/bankmaster", h.CreateBankMaster)
	h.ms.GET("/bankmaster/:id", h.InfoBankMaster)
	h.ms.PUT("/bankmaster/:id", h.UpdateBankMaster)
	h.ms.DELETE("/bankmaster/:id", h.DeleteBankMaster)
}

// Create Bank Master godoc
// @Summary		สร้าง bank master
// @Description สร้าง bank master
// @Tags		GL
// @Param		BankMaster  body      models.BankMaster  true  "bankmaster"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /bankmaster [post]
func (h BankMasterHttp) CreateBankMaster(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.BankMaster{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateBankMaster(shopID, authUsername, *docReq)

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

// Update Bank Master godoc
// @Summary		แก้ไข bank master
// @Description แก้ไข bank master
// @Tags		GL
// @Param		id  path      string  true  "Bank Master ID"
// @Param		BankMaster  body      models.BankMaster  true  "bankmaster"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /bankmaster/{id} [put]
func (h BankMasterHttp) UpdateBankMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.BankMaster{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateBankMaster(id, shopID, authUsername, *docReq)

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

// Delete Bank Master godoc
// @Summary		ลบ bank master
// @Description ลบ bank master
// @Tags		GL
// @Param		id  path      string  true  "Journal ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /bankmaster/{id} [delete]
func (h BankMasterHttp) DeleteBankMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteBankMaster(id, shopID, authUsername)

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

// Get Bank Master godoc
// @Summary		แสดงรายละเอียด bank master
// @Description แสดงรายละเอียด bank master
// @Tags		GL
// @Param		id  path      string  true  "Journal Id"
// @Accept 		json
// @Success		200	{object}	models.BankMasterInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /bankmaster/{id} [get]
func (h BankMasterHttp) InfoBankMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get BankMaster %v", id)
	doc, err := h.svc.InfoBankMaster(id, shopID)

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

// List Bank Master godoc
// @Summary		แสดงรายการ bank master
// @Description แสดงรายการ bank master
// @Tags		GL
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.BankMasterPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /bankmaster [get]
func (h BankMasterHttp) SearchBankMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchBankMaster(shopID, q, page, limit, sort)

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

// Create Bank Master godoc
// @Summary		นำเข้าข้อมูล bank master
// @Description นำเข้าข้อมูล bank master
// @Tags		GL
// @Param		Journal  body      []models.Journal  true  "bankmaster"
// @Accept 		json
// @Success		201	{object}	common.BulkInsertResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /bankmaster/bulk [post]
func (h BankMasterHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.BankMaster{}
	err := json.Unmarshal([]byte(input), &dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	bulkResponse, err := h.svc.SaveInBatch(shopID, authUsername, dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		common.BulkReponse{
			Success:    true,
			BulkImport: bulkResponse,
		},
	)

	return nil
}
