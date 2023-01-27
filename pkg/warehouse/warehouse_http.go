package warehouse

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/warehouse/models"
	"smlcloudplatform/pkg/warehouse/repositories"
	"smlcloudplatform/pkg/warehouse/services"
)

type IWarehouseHttp interface{}

type WarehouseHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IWarehouseHttpService
}

func NewWarehouseHttp(ms *microservice.Microservice, cfg microservice.IConfig) WarehouseHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := repositories.NewWarehouseRepository(pst)

	svc := services.NewWarehouseHttpService(repo)

	return WarehouseHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h WarehouseHttp) RouteSetup() {

	h.ms.POST("/warehouse/bulk", h.SaveBulk)

	h.ms.GET("/warehouse", h.SearchWarehouse)
	h.ms.POST("/warehouse", h.CreateWarehouse)
	h.ms.GET("/warehouse/:id", h.InfoWarehouse)
	h.ms.PUT("/warehouse/:id", h.UpdateWarehouse)
	h.ms.DELETE("/warehouse/:id", h.DeleteWarehouse)
}

// Create Warehouse godoc
// @Description Create Warehouse
// @Tags		Warehouse
// @Param		Warehouse  body      models.Warehouse  true  "Warehouse"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /warehouse [post]
func (h WarehouseHttp) CreateWarehouse(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.Warehouse{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateWarehouse(shopID, authUsername, *docReq)

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

// Update Warehouse godoc
// @Description Update Warehouse
// @Tags		Warehouse
// @Param		id  path      string  true  "Warehouse ID"
// @Param		Warehouse  body      models.Warehouse  true  "Warehouse"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /warehouse/{id} [put]
func (h WarehouseHttp) UpdateWarehouse(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Warehouse{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateWarehouse(shopID, id, authUsername, *docReq)

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

// Delete Warehouse godoc
// @Description Delete Warehouse
// @Tags		Warehouse
// @Param		id  path      string  true  "Warehouse ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /warehouse/{id} [delete]
func (h WarehouseHttp) DeleteWarehouse(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteWarehouse(shopID, id, authUsername)

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

// Get Warehouse godoc
// @Description get struct array by ID
// @Tags		Warehouse
// @Param		id  path      string  true  "Warehouse ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /warehouse/{id} [get]
func (h WarehouseHttp) InfoWarehouse(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Warehouse %v", id)
	doc, err := h.svc.InfoWarehouse(shopID, id)

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

// List Warehouse godoc
// @Description get struct array by ID
// @Tags		Warehouse
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /warehouse [get]
func (h WarehouseHttp) SearchWarehouse(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchWarehouse(shopID, pageable)

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

// Create Warehouse Bulk godoc
// @Description Create Warehouse
// @Tags		Warehouse
// @Param		Warehouse  body      []models.Warehouse  true  "Warehouse"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /warehouse/bulk [post]
func (h WarehouseHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Warehouse{}
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
