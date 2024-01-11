package shop

import (
	"encoding/json"
	"errors"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/authentication/models"
	"smlcloudplatform/pkg/config"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
)

type IShopMemberHttp interface{}

type ShopMemberHttp struct {
	ms  *microservice.Microservice
	svc IShopUserService
}

func NewShopMemberHttp(ms *microservice.Microservice, cfg config.IConfig) *ShopMemberHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	repo := NewShopUserRepository(pst)
	svc := NewShopUserService(repo)
	return &ShopMemberHttp{
		svc: svc,
		ms:  ms,
	}
}

func (h *ShopMemberHttp) RegisterHttp() {
	h.ms.GET("/user/permissions", h.ListShopUser)
	h.ms.GET("/shop/users", h.ListUserInShop)

	h.ms.PUT("/shop/permission", h.SaveUserPermissionShop)
	h.ms.GET("/shop/permission/:username", h.InfoShopUser)
	h.ms.DELETE("/shop/permission/:username", h.DeleteUserPermissionShop)
}

// List Shop User godoc
// @Description get shopuser
// @Tags		ShopUser
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/users [get]
func (h ShopMemberHttp) ListUserInShop(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	if userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &common.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.svc.ListUserInShop(shopID, pageable)

	if err != nil {
		ctx.ResponseError(400, "find failed")
		h.ms.Logger.Error("HTTP:: SearchShopUser " + err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Pagination: pagination,
		Data:       docList,
	})
	return nil
}

// Get Shop User godoc
// @Description get shopuser info by username
// @Tags		ShopUser
// @Param		username	path     string  true  "username"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/permission/{username} [get]
func (h ShopMemberHttp) InfoShopUser(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	if userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &common.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	username := ctx.Param("username")

	if len(username) < 1 {
		ctx.ResponseError(400, "username invalid")
		return nil
	}

	doc, err := h.svc.InfoShopByUser(shopID, username)

	if err != nil {
		ctx.ResponseError(400, "find failed")
		h.ms.Logger.Error("HTTP:: SearchShopUser " + err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// List Shop By User godoc
// @Description get shopuser
// @Tags		ShopUser
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /user/permissions [get]
func (h ShopMemberHttp) ListShopUser(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username

	if userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &common.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.svc.ListShopByUser(authUsername, pageable)

	if err != nil {
		ctx.ResponseError(400, "find failed")
		h.ms.Logger.Error("HTTP:: SearchShopUser " + err.Error())
		return err
	}

	ctx.Response(http.StatusOK,
		common.ApiResponse{
			Success:    true,
			Pagination: pagination,
			Data:       docList,
		})
	return nil
}

// Save Permission Shop User godoc
// @Description get shopuser info by username
// @Tags		ShopUser
// @Param		UserRoleRequest  body      models.UserRoleRequest  true  "UserRoleRequest"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/permission [put]
func (h ShopMemberHttp) SaveUserPermissionShop(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username

	if userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &common.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

	input := ctx.ReadInput()

	userRoleReq := &models.UserRoleRequest{}
	err := json.Unmarshal([]byte(input), &userRoleReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.SaveUserPermissionShop(userRoleReq.ShopID, authUsername, userRoleReq.EditUsername, userRoleReq.Username, userRoleReq.Role)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}
	ctx.Response(http.StatusOK,
		common.ApiResponse{
			Success: true,
		})

	return nil
}

// Delete Shop User godoc
// @Description get shopuser info by username
// @Tags		ShopUser
// @Param		username	path     string  true  "username"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /shop/permission/{username} [delete]
func (h ShopMemberHttp) DeleteUserPermissionShop(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	if userInfo.Role != models.ROLE_OWNER {
		ctx.Response(http.StatusOK, &common.ApiResponse{
			Success: false,
			Message: "permission denied",
		})

		return errors.New("permission denied")
	}

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
		common.ApiResponse{
			Success: true,
		})

	return nil
}
