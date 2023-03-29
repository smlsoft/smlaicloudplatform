package debtor

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/debtaccount/debtor/models"
	"smlcloudplatform/pkg/debtaccount/debtor/repositories"
	"smlcloudplatform/pkg/debtaccount/debtor/services"
	groupRepositories "smlcloudplatform/pkg/debtaccount/debtorgroup/repositories"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
)

type IDebtorHttp interface{}

type DebtorHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IDebtorHttpService
}

func NewDebtorHttp(ms *microservice.Microservice, cfg microservice.IConfig) DebtorHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewDebtorRepository(pst)
	repoGroup := groupRepositories.NewDebtorGroupRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewDebtorHttpService(repo, repoGroup, masterSyncCacheRepo)

	return DebtorHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h DebtorHttp) RouteSetup() {

	h.ms.POST("/debtaccount/debtor/bulk", h.SaveBulk)

	h.ms.GET("/debtaccount/debtor", h.SearchDebtorPage)
	h.ms.GET("/debtaccount/debtor/list", h.SearchDebtorStep)
	h.ms.POST("/debtaccount/debtor", h.CreateDebtor)
	h.ms.GET("/debtaccount/debtor/:id", h.InfoDebtor)
	h.ms.PUT("/debtaccount/debtor/:id", h.UpdateDebtor)
	h.ms.DELETE("/debtaccount/debtor/:id", h.DeleteDebtor)
	h.ms.DELETE("/debtaccount/debtor", h.DeleteDebtorByGUIDs)
}

// Create Debtor godoc
// @Description Create Debtor
// @Tags		Debtor
// @Param		Debtor  body      models.DebtorRequest  true  "Debtor"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /debtaccount/debtor [post]
func (h DebtorHttp) CreateDebtor(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.DebtorRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateDebtor(shopID, authUsername, *docReq)

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

// Update Debtor godoc
// @Description Update Debtor
// @Tags		Debtor
// @Param		id  path      string  true  "Debtor ID"
// @Param		Debtor  body      models.DebtorRequest  true  "Debtor"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /debtaccount/debtor/{id} [put]
func (h DebtorHttp) UpdateDebtor(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.DebtorRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateDebtor(shopID, id, authUsername, *docReq)

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

// Delete Debtor godoc
// @Description Delete Debtor
// @Tags		Debtor
// @Param		id  path      string  true  "Debtor ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /debtaccount/debtor/{id} [delete]
func (h DebtorHttp) DeleteDebtor(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteDebtor(shopID, id, authUsername)

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

// Delete Debtor godoc
// @Description Delete Debtor
// @Tags		Debtor
// @Param		Debtor  body      []string  true  "Debtor GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /debtaccount/debtor [delete]
func (h DebtorHttp) DeleteDebtorByGUIDs(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	input := ctx.ReadInput()

	docReq := []string{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.DeleteDebtorByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get Debtor godoc
// @Description get struct array by ID
// @Tags		Debtor
// @Param		id  path      string  true  "Debtor ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /debtaccount/debtor/{id} [get]
func (h DebtorHttp) InfoDebtor(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Debtor %v", id)
	doc, err := h.svc.InfoDebtor(shopID, id)

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

// List Debtor godoc
// @Description get struct array by ID
// @Tags		Debtor
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /debtaccount/debtor [get]
func (h DebtorHttp) SearchDebtorPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.svc.SearchDebtor(shopID, map[string]interface{}{}, pageable)

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

// List Debtor godoc
// @Description search limit offset
// @Tags		Debtor
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /debtaccount/debtor/list [get]
func (h DebtorHttp) SearchDebtorStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchDebtorStep(shopID, lang, pageableStep)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    docList,
		Total:   total,
	})
	return nil
}

// Create Debtor Bulk godoc
// @Description Create Debtor
// @Tags		Debtor
// @Param		Debtor  body      []models.DebtorRequest  true  "Debtor"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /debtaccount/debtor/bulk [post]
func (h DebtorHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.DebtorRequest{}
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
