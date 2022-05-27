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

	h.ms.POST("/gl/account-group/bulk", h.SaveBulk)

	h.ms.GET("/gl/account-group", h.SearchAccountGroup)
	h.ms.POST("/gl/account-group", h.CreateAccountGroup)
	h.ms.GET("/gl/account-group/:id", h.InfoAccountGroup)
	h.ms.PUT("/gl/account-group/:id", h.UpdateAccountGroup)
	h.ms.DELETE("/gl/account-group/:id", h.DeleteAccountGroup)
}

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

func (h AccountGroupHttp) SearchAccountGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	docList, pagination, err := h.svc.Search(shopID, q, page, limit)

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
