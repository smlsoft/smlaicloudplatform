package restaurantsettings

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/restaurantsettings/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"

	mastersync "smlcloudplatform/pkg/mastersync/repositories"
)

type IRestaurantSettingsHttp interface{}

type RestaurantSettingsHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc IRestaurantSettingsService
}

func NewRestaurantSettingsHttp(ms *microservice.Microservice, cfg microservice.IConfig) RestaurantSettingsHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := NewRestaurantSettingsRepository(pst)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := NewRestaurantSettingsService(repo, masterSyncCacheRepo)

	return RestaurantSettingsHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h RestaurantSettingsHttp) RouteSetup() {

	h.ms.POST("/restaurant/settings/bulk", h.SaveBulk)
	h.ms.GET("/restaurant/settings/fetchupdate", h.FetchUpdate)

	h.ms.GET("/restaurant/settings", h.SearchRestaurantSettings)
	h.ms.POST("/restaurant/settings", h.CreateRestaurantSettings)
	h.ms.GET("/restaurant/settings/:id", h.InfoRestaurantSettings)
	h.ms.PUT("/restaurant/settings/:id", h.UpdateRestaurantSettings)
	h.ms.DELETE("/restaurant/settings/:id", h.DeleteRestaurantSettings)

}

// Create Restaurant Settings godoc
// @Description Restaurant Settings
// @Tags		Restaurant
// @Param		Settings  body      models.RestaurantSettings  true  "Settings"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/settings [post]
func (h RestaurantSettingsHttp) CreateRestaurantSettings(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.RestaurantSettings{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateRestaurantSettings(shopID, authUsername, *docReq)

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

// Update Restaurant Settings godoc
// @Description Restaurant Settings
// @Tags		Restaurant
// @Param		id  path      string  true  "Settings ID"
// @Param		Settings  body      models.RestaurantSettings  true  "Settings"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/settings/{id} [put]
func (h RestaurantSettingsHttp) UpdateRestaurantSettings(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.RestaurantSettings{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateRestaurantSettings(shopID, id, authUsername, *docReq)

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

// Delete Restaurant Settings godoc
// @Description Restaurant Settings
// @Tags		Restaurant
// @Param		id  path      string  true  "RestaurantSettings ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/settings/{id} [delete]
func (h RestaurantSettingsHttp) DeleteRestaurantSettings(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteRestaurantSettings(shopID, id, authUsername)

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

// Get Restaurant Settings Infomation godoc
// @Description Get Restaurant Settings
// @Tags		Restaurant
// @Param		id  path      string  true  "RestaurantSettings Id"
// @Accept 		json
// @Success		200	{object}	models.RestaurantSettingsInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/settings/{id} [get]
func (h RestaurantSettingsHttp) InfoRestaurantSettings(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get RestaurantSettings %v", id)
	doc, err := h.svc.InfoRestaurantSettings(shopID, id)

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

// List Restaurant Settings godoc
// @Description List Restaurant Settings Category
// @Tags		Restaurant
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.RestaurantSettingsPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/settings [get]
func (h RestaurantSettingsHttp) SearchRestaurantSettings(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.svc.SearchRestaurantSettings(shopID, q, page, limit)

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

// Fetch Restaurant RestaurantSettings Update By Date godoc
// @Description Fetch Restaurant RestaurantSettings Update By Date
// @Tags		Restaurant
// @Param		lastUpdate query string true "DateTime YYYY-MM-DDTHH:mm"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept		json
// @Success		200 {object} models.RestaurantSettingsFetchUpdateResponse
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/restaurant/settings/fetchupdate [get]
func (h RestaurantSettingsHttp) FetchUpdate(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02T15:04" //
	lastUpdateStr := ctx.QueryParam("lastUpdate")

	lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(lastUpdateStr) < 1 {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return nil
	}

	lastUpdate, err := time.Parse(layout, lastUpdateStr)

	if err != nil {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return err
	}

	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.svc.LastActivity(shopID, "all", lastUpdate, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})

	return nil
}

// Create RestaurantSettings Bulk godoc
// @Description Create RestaurantSettings
// @Tags		Restaurant
// @Param		RestaurantSettings  body      []models.RestaurantSettings  true  "RestaurantSettings"
// @Accept 		json
// @Success		201	{object}	common.BulkInsertResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /restaurant/settings/bulk [post]
func (h RestaurantSettingsHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.RestaurantSettings{}
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
