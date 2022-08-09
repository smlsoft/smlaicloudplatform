package shop

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shop/models"
	"smlcloudplatform/pkg/utils"
)

type IShopMemberHttp interface{}

type ShopMemberHttp struct {
	ms  *microservice.Microservice
	svc IShopUserService
}

func NewShopMemberHttp(ms *microservice.Microservice, cfg microservice.IConfig) *ShopMemberHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	repo := NewShopUserRepository(pst)
	svc := NewShopUserService(repo)
	return &ShopMemberHttp{
		svc: svc,
		ms:  ms,
	}
}

func (h *ShopMemberHttp) RouteSetup() {
	h.ms.GET("/user/permissions", h.ListShopUser)
	h.ms.GET("/shop/users", h.ListUserInShop)
	h.ms.PUT("/shop/permission", h.SaveUserPermissionShop)
	h.ms.DELETE("/shop/permission/:username", h.DeleteUserPermissionShop)
}

func (h ShopMemberHttp) ListUserInShop(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)

	docList, pagination, err := h.svc.ListUserInShop(shopID, q, page, limit, sort)

	if err != nil {
		ctx.ResponseError(400, "find failed")
		h.ms.Logger.Error("HTTP:: SearchShopUser " + err.Error())
		return err
	}

	ctx.Response(http.StatusOK,
		map[string]interface{}{
			"success":    true,
			"pagination": pagination,
			"data":       docList,
		})
	return nil
}

func (h ShopMemberHttp) ListShopUser(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.svc.ListShopByUser(authUsername, q, page, limit)

	if err != nil {
		ctx.ResponseError(400, "find failed")
		h.ms.Logger.Error("HTTP:: SearchShopUser " + err.Error())
		return err
	}

	ctx.Response(http.StatusOK,
		map[string]interface{}{
			"success":    true,
			"pagination": pagination,
			"data":       docList,
		})
	return nil
}

func (h ShopMemberHttp) SaveUserPermissionShop(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username

	input := ctx.ReadInput()

	userRoleReq := &models.UserRoleRequest{}
	err := json.Unmarshal([]byte(input), &userRoleReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.SaveUserPermissionShop(userRoleReq.ShopID, authUsername, userRoleReq.Username, userRoleReq.Role)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}
	ctx.Response(http.StatusOK,
		map[string]interface{}{
			"success": true,
		})

	return nil
}

func (h ShopMemberHttp) DeleteUserPermissionShop(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	// input := ctx.ReadInput()

	// userRoleReq := &models.UserRoleRequest{}
	// err := json.Unmarshal([]byte(input), &userRoleReq)

	// if err != nil {
	// 	ctx.ResponseError(400, err.Error())
	// 	return err
	// }

	username := ctx.Param("username")

	if len(username) < 1 {
		ctx.ResponseError(400, "username invalid")
		return nil
	}

	err := h.svc.DeleteUserPermissionShop(shopID, authUsername, username)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK,
		map[string]interface{}{
			"success": true,
		})

	return nil
}
