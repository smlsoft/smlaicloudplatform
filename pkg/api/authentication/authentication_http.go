package authentication

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/shop"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IAuthenticationHttp interface {
	Login(ctx microservice.IContext) error
	Register(ctx microservice.IContext) error
	Logout(ctx microservice.IContext) error
	Profile(ctx microservice.IContext) error
}
type AuthenticationHttp struct {
	ms                    *microservice.Microservice
	cfg                   microservice.IConfig
	authService           *microservice.AuthService
	authenticationService IAuthenticationService
	shopUserService       shop.IShopUserService
}

func NewAuthenticationHttp(ms *microservice.Microservice, cfg microservice.IConfig) AuthenticationHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	authService := microservice.NewAuthService(ms.Cacher(cfg.CacherConfig()), 24*3)

	shopUserRepo := shop.NewShopUserRepository(pst)
	authRepo := NewAuthenticationRepository(pst)
	authenticationService := NewAuthenticationService(authRepo, shopUserRepo, authService)

	shopUserService := shop.NewShopUserService(shopUserRepo)
	return AuthenticationHttp{
		ms:                    ms,
		cfg:                   cfg,
		authService:           authService,
		authenticationService: authenticationService,
		shopUserService:       shopUserService,
	}
}

func (h *AuthenticationHttp) RouteSetup() {

	h.ms.POST("/login", h.Login)
	h.ms.POST("/logout", h.Logout)

	h.ms.POST("/register", h.Register)
	h.ms.GET("/profile", h.Profile)
	h.ms.PUT("/profile", h.Update)
	h.ms.PUT("/profile/password", h.UpdatePassword)

	h.ms.GET("/list-shop", h.ListShopCanAccess, h.authService.MWFuncWithShop(h.ms.Cacher(h.cfg.CacherConfig())))
	h.ms.POST("/select-shop", h.SelectShop, h.authService.MWFuncWithShop(h.ms.Cacher(h.cfg.CacherConfig())))
}

// Login login
// @Description get struct array by ID
// @Tags		Authentication
// @Param		User  body      models.UserRequest  true  "Add account"
// @Accept 		json
// @Success		200	{object}	models.AuthResponse
// @Failure		400 {object}	models.AuthResponseFailed
// @Router /login [post]
func (h *AuthenticationHttp) Login(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	userReq := &models.UserLoginRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	tokenString, err := h.authenticationService.Login(userReq)

	if err != nil {
		ctx.ResponseError(400, "login failed.")
		return err
	}

	ctx.Response(http.StatusOK, models.AuthResponse{
		Success: true,
		Token:   tokenString,
	})

	return nil
}

// RegisterMember register
// @Summary		Register An Account
// @Description	For User Register Application
// @Tags		Authentication
// @Param		User  body      models.UserRequest  true  "Add account"
// @Success		200	{object}	models.ResponseSuccessWithId
// @Failure		400 {object}	models.AuthResponseFailed
// @Accept 		json
// @Router		/register [post]
func (h *AuthenticationHttp) Register(ctx microservice.IContext) error {
	h.ms.Logger.Debug("Receive Register Data")
	input := ctx.ReadInput()

	userReq := models.UserRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	idx, err := h.authenticationService.Register(userReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		Id:      idx,
	})

	return nil
}

func (h *AuthenticationHttp) Update(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	input := ctx.ReadInput()

	userReq := models.UserRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	err = h.authenticationService.Update(authUsername, userReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
	})

	return nil
}

func (h *AuthenticationHttp) UpdatePassword(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	input := ctx.ReadInput()

	userPwdReq := models.UserPasswordRequest{}
	err := json.Unmarshal([]byte(input), &userPwdReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	err = h.authenticationService.UpdatePassword(authUsername, userPwdReq.CurrentPassword, userPwdReq.NewPassword)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
	})

	return nil
}

func (h *AuthenticationHttp) Logout(ctx microservice.IContext) error {

	authorizationHeader := ctx.Header("Authorization")

	err := h.authenticationService.Logout(authorizationHeader)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})

	return nil
}

func (h *AuthenticationHttp) Profile(ctx microservice.IContext) error {

	userProfile, err := h.authenticationService.Profile(ctx.UserInfo().Username)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    userProfile,
	})
	return nil
}

func (h *AuthenticationHttp) ListShop(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	authorizationHeader := ctx.Header("Authorization")

	input := ctx.ReadInput()

	shopSelectReq := &models.ShopSelectRequest{}
	err := json.Unmarshal([]byte(input), &shopSelectReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	err = h.authenticationService.AccessShop(authorizationHeader, shopSelectReq.ShopId, authUsername)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})

	return nil
}

// Access Shop godoc
// @Description Access to Shop
// @Tags		Authentication
// @Param		User  body      models.ShopSelectRequest  true  "Shop"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
// @Security     AccessToken
// @Router /select-shop [post]
func (h *AuthenticationHttp) SelectShop(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	authorizationHeader := ctx.Header("Authorization")

	input := ctx.ReadInput()

	shopSelectReq := &models.ShopSelectRequest{}
	err := json.Unmarshal([]byte(input), &shopSelectReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	err = h.authenticationService.AccessShop(authorizationHeader, shopSelectReq.ShopId, authUsername)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})

	return nil
}

func (h *AuthenticationHttp) ListShopCanAccess(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	docList, pagination, err := h.shopUserService.ListShopByUser(authUsername, page, limit)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
	}

	ctx.Response(http.StatusOK,
		models.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		},
	)

	return nil
}
