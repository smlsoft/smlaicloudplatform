package shopzone

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/shopzone/models"
	"smlcloudplatform/pkg/utils"

	mastersync "smlcloudplatform/pkg/mastersync/repositories"
)

type IShopZoneHttp interface{}

type ShopZoneHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc IShopZoneService
}

func NewShopZoneHttp(ms *microservice.Microservice, cfg config.IConfig) ShopZoneHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := NewShopZoneRepository(pst)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := NewShopZoneService(repo, masterSyncCacheRepo)

	return ShopZoneHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ShopZoneHttp) RouteSetup() {

	h.ms.POST("/restaurant/zone/bulk", h.SaveBulk)

	h.ms.GET("/restaurant/zone", h.SearchShopZone)
	h.ms.POST("/restaurant/zone", h.CreateShopZone)
	h.ms.GET("/restaurant/zone/:id", h.InfoShopZone)
	h.ms.PUT("/restaurant/zone/:id", h.UpdateShopZone)
	h.ms.DELETE("/restaurant/zone/:id", h.DeleteShopZone)

}

// Create Restaurant Shop Zone godoc
// @Description Restaurant Shop Zone
// @Tags		Restaurant
// @Param		Zone  body      models.ShopZone  true  "Zone"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone [post]
func (h ShopZoneHttp) CreateShopZone(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.ShopZone{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateShopZone(shopID, authUsername, *docReq)

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

// Update Restaurant Shop Zone godoc
// @Description Restaurant Shop Zone
// @Tags		Restaurant
// @Param		id  path      string  true  "Zone ID"
// @Param		Zone  body      models.ShopZone  true  "Zone"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone/{id} [put]
func (h ShopZoneHttp) UpdateShopZone(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.ShopZone{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateShopZone(shopID, id, authUsername, *docReq)

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

// Delete Restaurant Shop Zone godoc
// @Description Restaurant Shop Zone
// @Tags		Restaurant
// @Param		id  path      string  true  "ShopZone ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone/{id} [delete]
func (h ShopZoneHttp) DeleteShopZone(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteShopZone(shopID, id, authUsername)

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

// Get Restaurant Shop Zone Infomation godoc
// @Description Get Restaurant Shop Zone
// @Tags		Restaurant
// @Param		id  path      string  true  "ShopZone Id"
// @Accept 		json
// @Success		200	{object}	models.ShopZoneInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone/{id} [get]
func (h ShopZoneHttp) InfoShopZone(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get ShopZone %v", id)
	doc, err := h.svc.InfoShopZone(shopID, id)

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

// List Restaurant Shop Zone godoc
// @Description List Restaurant Shop Zone Category
// @Tags		Restaurant
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.ShopZonePageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone [get]
func (h ShopZoneHttp) SearchShopZone(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.svc.SearchShopZone(shopID, pageable)

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

// Create ShopZone Bulk godoc
// @Description Create ShopZone
// @Tags		Restaurant
// @Param		ShopZone  body      []models.ShopZone  true  "ShopZone"
// @Accept 		json
// @Success		201	{object}	common.BulkInsertResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone/bulk [post]
func (h ShopZoneHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.ShopZone{}
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
