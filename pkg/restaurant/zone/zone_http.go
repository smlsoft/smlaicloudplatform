package zone

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/zone/models"
	"smlcloudplatform/pkg/utils"

	mastersync "smlcloudplatform/pkg/mastersync/repositories"
)

type IZoneHttp interface{}

type ZoneHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc IZoneService
}

func NewZoneHttp(ms *microservice.Microservice, cfg config.IConfig) ZoneHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := NewZoneRepository(pst)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := NewZoneService(repo, masterSyncCacheRepo)

	return ZoneHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ZoneHttp) RouteSetup() {

	h.ms.POST("/restaurant/zone/bulk", h.SaveBulk)

	h.ms.GET("/restaurant/zone", h.SearchZone)
	h.ms.POST("/restaurant/zone", h.CreateZone)
	h.ms.GET("/restaurant/zone/:id", h.InfoZone)
	h.ms.PUT("/restaurant/zone/:id", h.UpdateZone)
	h.ms.DELETE("/restaurant/zone/:id", h.DeleteZone)
	h.ms.DELETE("/restaurant/zone", h.DeleteByGUIDs)

}

// Create Restaurant  Zone godoc
// @Description Restaurant  Zone
// @Tags		Restaurant
// @Param		Zone  body      models.Zone  true  "Zone"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone [post]
func (h ZoneHttp) CreateZone(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.Zone{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateZone(shopID, authUsername, *docReq)

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

// Update Restaurant  Zone godoc
// @Description Restaurant  Zone
// @Tags		Restaurant
// @Param		id  path      string  true  "Zone ID"
// @Param		Zone  body      models.Zone  true  "Zone"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone/{id} [put]
func (h ZoneHttp) UpdateZone(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Zone{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateZone(shopID, id, authUsername, *docReq)

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

// Delete Restaurant  Zone godoc
// @Description Restaurant  Zone
// @Tags		Restaurant
// @Param		id  path      string  true  "Zone ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone/{id} [delete]
func (h ZoneHttp) DeleteZone(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteZone(shopID, id, authUsername)

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

// Delete Zone godoc
// @Description Delete Zone
// @Tags		Zone
// @Param		Zone  body      []string  true  "Zone GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone [delete]
func (h ZoneHttp) DeleteByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get Restaurant  Zone Infomation godoc
// @Description Get Restaurant  Zone
// @Tags		Restaurant
// @Param		id  path      string  true  "Zone Id"
// @Accept 		json
// @Success		200	{object}	models.ZoneInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone/{id} [get]
func (h ZoneHttp) InfoZone(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Zone %v", id)
	doc, err := h.svc.InfoZone(shopID, id)

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

// List Restaurant  Zone godoc
// @Description List Restaurant  Zone Category
// @Tags		Restaurant
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.ZonePageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone [get]
func (h ZoneHttp) SearchZone(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.svc.SearchZone(shopID, pageable)

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

// Create Zone Bulk godoc
// @Description Create Zone
// @Tags		Restaurant
// @Param		Zone  body      []models.Zone  true  "Zone"
// @Accept 		json
// @Success		201	{object}	common.BulkInsertResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/zone/bulk [post]
func (h ZoneHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Zone{}
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
