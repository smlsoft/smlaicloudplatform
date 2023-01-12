package authentication

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/shop"
	shopmodel "smlcloudplatform/pkg/shop/models"
	"smlcloudplatform/pkg/utils"
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
	shopService           shop.IShopService
	shopUserService       shop.IShopUserService
}

func NewAuthenticationHttp(ms *microservice.Microservice, cfg microservice.IConfig) AuthenticationHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	authService := microservice.NewAuthService(ms.Cacher(cfg.CacherConfig()), 24*3)

	shopRepo := shop.NewShopRepository(pst)
	shopUserRepo := shop.NewShopUserRepository(pst)
	shopUserAccessLogRepo := shop.NewShopUserAccessLogRepository(pst)
	authRepo := NewAuthenticationRepository(pst)
	authenticationService := NewAuthenticationService(authRepo, shopUserRepo, shopUserAccessLogRepo, authService, utils.HashPassword, utils.CheckHashPassword, ms.TimeNow)

	shopService := shop.NewShopService(shopRepo, shopUserRepo, utils.NewGUID, ms.TimeNow)
	shopUserService := shop.NewShopUserService(shopUserRepo)
	return AuthenticationHttp{
		ms:                    ms,
		cfg:                   cfg,
		authService:           authService,
		authenticationService: authenticationService,
		shopUserService:       shopUserService,
		shopService:           shopService,
	}
}

func (h AuthenticationHttp) RouteSetup() {

	h.ms.POST("/login", h.Login)
	h.ms.POST("/logout", h.Logout)

	h.ms.POST("/register", h.Register)
	h.ms.GET("/profile", h.Profile)
	h.ms.GET("/profileshop", h.ProfileShop)

	h.ms.PUT("/profile", h.Update)
	h.ms.PUT("/profile/password", h.UpdatePassword)

	middlewareShop := h.authService.MWFuncWithShop(h.ms.Cacher(h.cfg.CacherConfig()))
	h.ms.GET("/list-shop", h.ListShopCanAccess, middlewareShop)
	h.ms.POST("/select-shop", h.SelectShop, middlewareShop)
	h.ms.PUT("/favorite-shop", h.UpdateShopFavorite, middlewareShop)

	shopHttp := shop.NewShopHttp(h.ms, h.cfg)
	h.ms.POST("/create-shop", shopHttp.CreateShop, middlewareShop)
}

// Login login
// @Description get struct array by ID
// @Tags		Authentication
// @Param		User  body      shopmodel.UserLoginRequest  true  "User Account"
// @Accept 		json
// @Success		200	{object}	common.AuthResponse
// @Failure		400 {object}	common.AuthResponseFailed
// @Router /login [post]
func (h AuthenticationHttp) Login(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	userReq := &shopmodel.UserLoginRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(userReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	authContext := AuthenticationContext{
		Ip: ctx.RealIp(),
	}

	tokenString, err := h.authenticationService.Login(userReq, authContext)

	if err != nil {
		ctx.ResponseError(400, "login failed.")
		return err
	}

	ctx.Response(http.StatusOK, common.AuthResponse{
		Success: true,
		Token:   tokenString,
	})

	return nil
}

// Register Member godoc
// @Summary		Register An Account
// @Description	For User Register Application
// @Tags		Authentication
// @Param		User  body      shopmodel.UserRequest  true  "Register account"
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		400 {object}	common.AuthResponseFailed
// @Accept 		json
// @Router		/register [post]
func (h AuthenticationHttp) Register(ctx microservice.IContext) error {
	h.ms.Logger.Debug("Receive Register Data")
	input := ctx.ReadInput()

	userReq := shopmodel.UserRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(userReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.authenticationService.Register(userReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

func (h AuthenticationHttp) Update(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	input := ctx.ReadInput()

	userReq := shopmodel.UserProfileRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(userReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.authenticationService.Update(authUsername, userReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})

	return nil
}

func (h AuthenticationHttp) UpdatePassword(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	input := ctx.ReadInput()

	userPwdReq := shopmodel.UserPasswordRequest{}
	err := json.Unmarshal([]byte(input), &userPwdReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(userPwdReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.authenticationService.UpdatePassword(authUsername, userPwdReq.CurrentPassword, userPwdReq.NewPassword)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Logout
// @Description Logout Current Profile
// @Tags		Authentication
// @Accept 		json
// @Success		200	{array}	shopmodel.UserProfileReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /logout [post]
func (h AuthenticationHttp) Logout(ctx microservice.IContext) error {

	authorizationHeader := ctx.Header("Authorization")

	err := h.authenticationService.Logout(authorizationHeader)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get Current Profile
// @Description Get Current Profile
// @Tags		Authentication
// @Accept 		json
// @Success		200	{array}	shopmodel.UserProfileReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /profile [get]
func (h AuthenticationHttp) Profile(ctx microservice.IContext) error {

	userProfile, err := h.authenticationService.Profile(ctx.UserInfo().Username)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    userProfile,
	})
	return nil
}

// Get Current Profile
// @Description Get Current Profile
// @Tags		Authentication
// @Accept 		json
// @Success		200	{array}	shopmodel.UserProfileReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /profileshop [get]
func (h AuthenticationHttp) ProfileShop(ctx microservice.IContext) error {

	userProfile, err := h.shopService.InfoShop(ctx.UserInfo().ShopID)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    userProfile,
	})
	return nil
}

// func (h AuthenticationHttp) ListShop(ctx microservice.IContext) error {
// 	authUsername := ctx.UserInfo().Username
// 	authorizationHeader := ctx.Header("Authorization")

// 	input := ctx.ReadInput()

// 	shopSelectReq := &shopmodel.ShopSelectRequest{}
// 	err := json.Unmarshal([]byte(input), &shopSelectReq)

// 	if err != nil {
// 		ctx.Response(http.StatusBadRequest, common.ApiResponse{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 		return err
// 	}

// 	err = h.authenticationService.AccessShop(shopSelectReq.ShopID, authUsername, authorizationHeader)

// 	if err != nil {
// 		ctx.Response(http.StatusBadRequest, common.ApiResponse{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 		return err
// 	}

// 	ctx.Response(http.StatusOK, common.ApiResponse{
// 		Success: true,
// 	})

// 	return nil
// }

// Access Shop godoc
// @Description Access to Shop
// @Tags		Authentication
// @Param		User  body      shopmodel.ShopSelectRequest  true  "Shop"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.ApiResponse
// @Security     AccessToken
// @Router /select-shop [post]
func (h AuthenticationHttp) SelectShop(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	authorizationHeader := ctx.Header("Authorization")

	input := ctx.ReadInput()

	shopSelectReq := &shopmodel.ShopSelectRequest{}
	err := json.Unmarshal([]byte(input), &shopSelectReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	authContext := AuthenticationContext{
		Ip: ctx.RealIp(),
	}

	err = h.authenticationService.AccessShop(shopSelectReq.ShopID, authUsername, authorizationHeader, authContext)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// List Shop godoc
// @Description List Merchant In My Account
// @Tags		Authentication
// @Accept 		json
// @Success		200	{array}	shopmodel.ShopUserInfo
// @Failure		401 {object}	common.ApiResponse
// @Security     AccessToken
// @Router /list-shop [get]
func (h AuthenticationHttp) ListShopCanAccess(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	docList, pagination, err := h.shopUserService.ListShopByUser(authUsername, q, page, limit)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return nil
	}

	ctx.Response(http.StatusOK,
		common.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		},
	)

	return nil
}

// Favorite Shop godoc
// @Description Favorite Shop In Account
// @Tags		Authentication
// @Accept 		json
// @Param		ShopFavoriteRequest  body      ShopFavoriteRequest  true  "Shop Favorite Request"
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.ApiResponse
// @Security     AccessToken
// @Router /favorite-shop [put]
func (h AuthenticationHttp) UpdateShopFavorite(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username

	input := ctx.ReadInput()

	reqBody := ShopFavoriteRequest{}
	err := json.Unmarshal([]byte(input), &reqBody)

	if err != nil {
		ctx.ResponseError(400, "request payload invalid")
		return err
	}

	err = h.authenticationService.UpdateFavoriteShop(reqBody.ShopID, authUsername, reqBody.IsFavorite)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return nil
	}

	ctx.Response(http.StatusOK,
		common.ApiResponse{
			Success: true,
		},
	)

	return nil
}
