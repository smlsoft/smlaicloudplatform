package zonedesign

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shopdesign/zonedesign/models"
	"smlcloudplatform/pkg/shopdesign/zonedesign/repositories"
	"smlcloudplatform/pkg/shopdesign/zonedesign/services"
	"smlcloudplatform/pkg/utils"

	common "smlcloudplatform/pkg/models"
)

type ZoneDesignHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IZoneDesignService
}

func NewZoneDesignHttp(ms *microservice.Microservice, cfg microservice.IConfig) ZoneDesignHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := repositories.NewZoneDesignRepository(pst)
	svc := services.NewZoneDesignService(repo)

	return ZoneDesignHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ZoneDesignHttp) RouteSetup() {

	h.ms.GET("/zonedesign", h.SearchZoneDesign)
	h.ms.POST("/zonedesign", h.CreateZoneDesign)
	h.ms.GET("/zonedesign/:id", h.InfoZoneDesign)
	h.ms.PUT("/zonedesign/:id", h.UpdateZoneDesign)
	h.ms.DELETE("/zonedesign/:id", h.DeleteZoneDesign)
}

func (h ZoneDesignHttp) CreateZoneDesign(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.ZoneDesign{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateZoneDesign(shopID, authUsername, *docReq)

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

func (h ZoneDesignHttp) UpdateZoneDesign(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.ZoneDesign{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateZoneDesign(shopID, id, authUsername, *docReq)

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

func (h ZoneDesignHttp) DeleteZoneDesign(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteZoneDesign(shopID, id, authUsername)

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

func (h ZoneDesignHttp) InfoZoneDesign(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get ZoneDesign %v", id)
	doc, err := h.svc.InfoZoneDesign(shopID, id)

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

func (h ZoneDesignHttp) SearchZoneDesign(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchZoneDesign(shopID, q, page, limit)

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
