package accountgroup

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/vfgl/accountgroup/models"
	"smlcloudplatform/pkg/vfgl/accountgroup/repositories"
	"smlcloudplatform/pkg/vfgl/accountgroup/services"

	common "smlcloudplatform/pkg/models"
)

type AccountGroupHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IAccountGroupHttpService
}

func NewAccountGroupHttp(ms *microservice.Microservice, cfg microservice.IConfig) AccountGroupHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	mongoRepo := repositories.NewAccountGroupMongoRepository(pst)
	mqRepo := repositories.NewAccountGroupMqRepository(prod)
	svc := services.NewAccountGroupHttpService(mongoRepo, mqRepo)

	return AccountGroupHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h AccountGroupHttp) RouteSetup() {

	h.ms.POST("/gl/accountgroup/bulk", h.SaveBulk)

	h.ms.GET("/gl/accountgroup", h.SearchAccountGroup)
	h.ms.POST("/gl/accountgroup", h.CreateAccountGroup)
	h.ms.GET("/gl/accountgroup/:id", h.InfoAccountGroup)
	h.ms.PUT("/gl/accountgroup/:id", h.UpdateAccountGroup)
	h.ms.DELETE("/gl/accountgroup/:id", h.DeleteAccountGroup)
}

// Create Account Group godoc
// @Summary		สร้างกลุ่มบัญชี
// @Description สร้างกลุ่มบัญชี
// @Tags		GL
// @Param		AccountGroup  body      models.AccountGroup  true  "กลุ่มบัญชี"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/accountgroup [post]
func (h AccountGroupHttp) CreateAccountGroup(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.AccountGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.Create(shopID, authUsername, *docReq)

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

// Update Account Group godoc
// @Summary		แก้ไขกลุ่มบัญชี
// @Description แก้ไขกลุ่มบัญชี
// @Tags		GL
// @Param		id  path      string  true  "ID"
// @Param		Journal  body      models.AccountGroup  true  "กลุ่มบัญชี"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/accountgroup/{id} [put]
func (h AccountGroupHttp) UpdateAccountGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.AccountGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.Update(id, shopID, authUsername, *docReq)

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

// Delete CAccount Group godoc
// @Summary		ลบกลุ่มบัญชี
// @Description ลบกลุ่มบัญชี
// @Tags		GL
// @Param		id  path      string  true  "ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/accountgroup/{id} [delete]
func (h AccountGroupHttp) DeleteAccountGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.Delete(id, shopID, authUsername)

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

// Get Account Group Infomation godoc
// @Summary		แสดงรายละเอียดกลุ่มบัญชี
// @Description แสดงรายละเอียดกลุ่มบัญชี
// @Tags		GL
// @Param		id  path      string  true  "Id"
// @Accept 		json
// @Success		200	{object}	models.AccountGroupInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/accountgroup/{id} [get]
func (h AccountGroupHttp) InfoAccountGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get doc %v", id)
	doc, err := h.svc.Info(id, shopID)

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

// List Account Group godoc
// @Summary		แสดงรายการกลุ่มบัญชี
// @Description แสดงรายการกลุ่มบัญชี
// @Tags		GL
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.AccountGroupPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/accountgroup [get]
func (h AccountGroupHttp) SearchAccountGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.svc.Search(shopID, q, page, limit, sort)

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

// Create Account Group Bulk godoc
// @Summary		นำเข้ากลุ่มบัญชี
// @Description นำเข้ากลุ่มบัญชี
// @Tags		GL
// @Param		AccountGroup  body      []models.AccountGroup  true  "กลุ่มบัญชี"
// @Accept 		json
// @Success		201	{object}	common.BulkInsertResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/chartofaccount/bulk [post]
func (h AccountGroupHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.AccountGroup{}
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
