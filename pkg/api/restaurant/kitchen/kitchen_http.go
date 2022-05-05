package kitchen

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"strconv"
)

type IKitchenHttp interface{}

type KitchenHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc IKitchenService
}

func NewKitchenHttp(ms *microservice.Microservice, cfg microservice.IConfig) KitchenHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	crudRepo := repositories.NewCrudRepository[restaurant.KitchenDoc](pst)
	searchRepo := repositories.NewSearchRepository[restaurant.KitchenInfo](pst)
	svc := NewKitchenService(crudRepo, searchRepo)

	return KitchenHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h KitchenHttp) RouteSetup() {

	h.ms.GET("/restaurant/kitchen/:id", h.InfoKitchen)
	h.ms.GET("/restaurant/kitchen", h.SearchKitchen)
	h.ms.POST("/restaurant/kitchen", h.CreateKitchen)
	h.ms.PUT("/restaurant/kitchen/:id", h.UpdateKitchen)
	h.ms.DELETE("/restaurant/kitchen/:id", h.DeleteKitchen)
}

func (h KitchenHttp) CreateKitchen(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &restaurant.Kitchen{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateKitchen(shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      idx,
	})
	return nil
}

func (h KitchenHttp) UpdateKitchen(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &restaurant.Kitchen{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateKitchen(id, shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

func (h KitchenHttp) DeleteKitchen(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteKitchen(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

func (h KitchenHttp) InfoKitchen(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Kitchen %v", id)
	doc, err := h.svc.InfoKitchen(id, shopID)

	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v: %v", id, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

func (h KitchenHttp) SearchKitchen(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}
	docList, pagination, err := h.svc.SearchKitchen(shopID, q, page, limit)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success:    true,
		Data:       docList,
		Pagination: pagination,
	})
	return nil
}
