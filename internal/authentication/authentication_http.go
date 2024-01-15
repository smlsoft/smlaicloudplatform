package authentication

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/authentication/models"
	"smlcloudplatform/internal/authentication/repositories"
	"smlcloudplatform/internal/authentication/services"
	"smlcloudplatform/internal/config"
	"smlcloudplatform/internal/firebase"
	common "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/shop"
	"smlcloudplatform/internal/utils"
	"smlcloudplatform/pkg/microservice"
	"time"
)

type IAuthenticationHttp interface {
	Login(ctx microservice.IContext) error
	TokenLogin(ctx microservice.IContext) error
	Register(ctx microservice.IContext) error
	Logout(ctx microservice.IContext) error
	Profile(ctx microservice.IContext) error
}
type AuthenticationHttp struct {
	ms                    *microservice.Microservice
	cfg                   config.IConfig
	authService           *microservice.AuthService
	authenticationService services.IAuthenticationService
	shopService           shop.IShopService
	shopUserService       shop.IShopUserService
}

func NewAuthenticationHttp(ms *microservice.Microservice, cfg config.IConfig) AuthenticationHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	authService := microservice.NewAuthService(ms.Cacher(cfg.CacherConfig()), 24*3*time.Hour, 24*30*time.Hour)

	shopRepo := shop.NewShopRepository(pst)
	shopUserRepo := shop.NewShopUserRepository(pst)
	shopUserAccessLogRepo := shop.NewShopUserAccessLogRepository(pst)
	// authRepo := NewAuthenticationRepository(pst)
	authRepo := repositories.NewAuthenticationMongoCacheRepository(pst, cache)
	smsRepo := repositories.NewAuthenticationSMSRepository(cache)
	firebaseAdapter := firebase.NewFirebaseAdapter()
	authenticationService := services.NewAuthenticationService(
		authRepo,
		shopUserRepo,
		shopUserAccessLogRepo,
		smsRepo,
		authService,
		utils.RandStringBytesMaskImprSrcUnsafe,
		utils.RandNumber,
		utils.NewGUID,
		utils.HashPassword,
		utils.CheckHashPassword,
		ms.TimeNow,
		firebaseAdapter)

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

func (h AuthenticationHttp) RegisterHttp() {

	h.ms.POST("/login", h.Login)
	h.ms.POST("/login/phone-number", h.LoginWithPhoneNumber)
	h.ms.POST("/tokenlogin", h.TokenLogin)
	h.ms.POST("/logout", h.Logout)
	h.ms.POST("/refresh", h.RefreshToken)

	h.ms.POST("/register", h.Register)
	h.ms.POST("/send-phonenumber-otp", h.SendPhoneNumberOTP)
	h.ms.POST("/forgot-password-phonenumber", h.ForgotPasswordByPhoneNumber)
	h.ms.POST("/register-phonenumber", h.RegisterByPhoneNumber)
	h.ms.POST("/register/exists-phonenumber", h.RegisterCheckExistPhonenumber)
	h.ms.POST("/register/exists-username", h.RegisterCheckExistUsername)

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

// Login with phone number
// @Description Login with phone number
// @Tags		Authentication
// @Param		UserLoginPhoneNumberRequest  body      models.UserLoginPhoneNumberRequest  true  "User Login PhoneNumber Request"
// @Accept 		json
// @Success		200	{object}	common.AuthResponse
// @Failure		400 {object}	common.AuthResponseFailed
// @Router /login [post]
func (h AuthenticationHttp) LoginWithPhoneNumber(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	userReq := &models.UserLoginPhoneNumberRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(userReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	authContext := models.AuthenticationContext{
		Ip: ctx.RealIp(),
	}

	result, err := h.authenticationService.LoginWithPhoneNumber(userReq, authContext)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{
		"success": true,
		"token":   result.Token,
		"refresh": result.Refresh,
	})

	return nil
}

// Login login
// @Description get struct array by ID
// @Tags		Authentication
// @Param		User  body      models.UserLoginRequest  true  "User Account"
// @Accept 		json
// @Success		200	{object}	common.AuthResponse
// @Failure		400 {object}	common.AuthResponseFailed
// @Router /login [post]
func (h AuthenticationHttp) Login(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	userReq := &models.UserLoginRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(userReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	authContext := models.AuthenticationContext{
		Ip: ctx.RealIp(),
	}

	result, err := h.authenticationService.Login(userReq, authContext)

	if err != nil {
		ctx.ResponseError(400, "login failed.")
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{
		"success": true,
		"token":   result.Token,
		"refresh": result.Refresh,
	})

	return nil
}

// Login refresh
// @Description refresh token
// @Tags		Authentication
// @Param		TokenLoginRequest  body      models.TokenLoginRequest  true  "Reresh Token"
// @Accept 		json
// @Success		200	{object}	common.AuthResponse
// @Failure		400 {object}	common.AuthResponseFailed
// @Router /refresh [post]
func (h AuthenticationHttp) RefreshToken(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	reqBody := models.TokenLoginRequest{}
	err := json.Unmarshal([]byte(input), &reqBody)

	if err != nil {
		ctx.ResponseError(400, "payload invalid")
		return err
	}

	if err = ctx.Validate(reqBody); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	result, err := h.authenticationService.RefreshToken(reqBody)

	if err != nil {
		ctx.ResponseError(400, "login failed.")
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{
		"success": true,
		"token":   result.Token,
		"refresh": result.Refresh,
	})

	return nil
}

// Login login
// @Description get struct array by ID
// @Tags		Authentication
// @Param		TokenLoginRequest  body      models.TokenLoginRequest  true  "User Account"
// @Accept 		json
// @Success		200	{object}	common.AuthResponse
// @Failure		400 {object}	common.AuthResponseFailed
// @Router /tokenlogin [post]
func (h AuthenticationHttp) TokenLogin(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	tokenReq := &models.TokenLoginRequest{}
	err := json.Unmarshal([]byte(input), &tokenReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	tokenString, err := h.authenticationService.LoginWithFirebaseToken(tokenReq.Token)

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
// @Param		RegisterEmailRequest  body      models.RegisterEmailRequest  true  "Register account"
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		400 {object}	common.AuthResponseFailed
// @Accept 		json
// @Router		/register [post]
func (h AuthenticationHttp) Register(ctx microservice.IContext) error {
	h.ms.Logger.Debug("Receive Register Data")
	input := ctx.ReadInput()

	userReq := models.RegisterEmailRequest{}
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

// Send Phonenumber OTP godoc
// @Summary		Send Phonenumber OTP
// @Description	For User Send Phonenumber OTP
// @Tags		Authentication
// @Param		OTPRequest  body      models.OTPRequest  true  "OTP Request"
// @Success		200	{object}	common.ApiResponse
// @Failure		400 {object}	common.AuthResponseFailed
// @Accept 		json
// @Router		/send-phonenumber-otp [post]
func (h AuthenticationHttp) SendPhoneNumberOTP(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	payload := models.OTPRequest{}
	err := json.Unmarshal([]byte(input), &payload)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(payload); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	result, err := h.authenticationService.SendPhonenumberOTP(payload)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Data:    result,
	})

	return nil
}

// Register By Phonenumber  godoc
// @Summary		Register By Phonenumber
// @Description	For User Register Phonenumber
// @Tags		Authentication
// @Param		OTPRequest  body      models.OTPRequest  true  "OTP Request"
// @Success		200	{object}	common.ApiResponse
// @Failure		400 {object}	common.AuthResponseFailed
// @Accept 		json
// @Router		/register-phonenumber [post]
func (h AuthenticationHttp) RegisterByPhoneNumber(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	payload := models.RegisterPhoneNumberRequest{}
	err := json.Unmarshal([]byte(input), &payload)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(payload); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.authenticationService.RegisterByPhonenumber(payload)

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

// Forgot Password By Phonenumber  godoc
// @Summary		Forgot Password By Phonenumber
// @Description	For User Forgot Password Phonenumber
// @Tags		Authentication
// @Param		ForgotPasswordPhoneNumberRequest  body      models.ForgotPasswordPhoneNumberRequest  true  "Forgot Password PhoneNumber Request"
// @Success		200	{object}	common.ApiResponse
// @Failure		400 {object}	common.AuthResponseFailed
// @Accept 		json
// @Router		/forgot-password-phonenumber [post]
func (h AuthenticationHttp) ForgotPasswordByPhoneNumber(ctx microservice.IContext) error {
	input := ctx.ReadInput()

	payload := models.ForgotPasswordPhoneNumberRequest{}
	err := json.Unmarshal([]byte(input), &payload)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(payload); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.authenticationService.ForgotPasswordByPhonenumber(payload)

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

// Register Check Exists Username godoc
// @Summary		Register Check Exists Username
// @Description	Check Exists Username
// @Tags		Authentication
// @Param		Username  body      models.UsernameField  true  "Username"
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		400 {object}	common.AuthResponseFailed
// @Accept 		json
// @Router		/register/exists-username [post]
func (h AuthenticationHttp) RegisterCheckExistUsername(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	payload := models.UsernameField{}
	err := json.Unmarshal([]byte(input), &payload)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(payload); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	isExists, err := h.authenticationService.CheckExistsUsername(payload.Username)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Data:    isExists,
	})

	return nil
}

// Register Check Exists Phone Number godoc
// @Summary		Register Check Exists Phone Number
// @Description	Check Exists Phone Number
// @Tags		Authentication
// @Param		PhoneNumber  body      models.PhoneNumberField  true  "Username"
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		400 {object}	common.AuthResponseFailed
// @Accept 		json
// @Router		/register/exists-phonenumber [post]
func (h AuthenticationHttp) RegisterCheckExistPhonenumber(ctx microservice.IContext) error {

	input := ctx.ReadInput()

	payload := models.PhoneNumberField{}
	err := json.Unmarshal([]byte(input), &payload)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	if err = ctx.Validate(payload); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	isExists, err := h.authenticationService.CheckExistsPhonenumber(payload.PhoneNumber)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Data:    isExists,
	})

	return nil
}

// Update User Profile godoc
// @Summary		Update profile
// @Description	For User Update Profile
// @Tags		Authentication
// @Param		UserProfileRequest  body      models.UserProfileRequest  true  "Update account"
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		400 {object}	common.AuthResponseFailed
// @Accept 		json
// @Router		/profile [put]
func (h AuthenticationHttp) Update(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	input := ctx.ReadInput()

	userReq := models.UserProfileRequest{}
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

	userPwdReq := models.UserPasswordRequest{}
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
// @Success		200	{array}	common.ApiResponse
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
// @Success		200	{array}	models.UserProfileReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /profile [get]
func (h AuthenticationHttp) Profile(ctx microservice.IContext) error {

	// stime := time.Now()
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

	// fmt.Println("Time Profile", time.Since(stime))
	return nil
}

// Get Current Profile
// @Description Get Current Profile
// @Tags		Authentication
// @Accept 		json
// @Success		200	{array}	models.UserProfileReponse
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

// Access Shop godoc
// @Description Access to Shop
// @Tags		Authentication
// @Param		User  body      models.ShopSelectRequest  true  "Shop"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.ApiResponse
// @Security     AccessToken
// @Router /select-shop [post]
func (h AuthenticationHttp) SelectShop(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	authorizationHeader := ctx.Header("Authorization")

	input := ctx.ReadInput()

	shopSelectReq := &models.ShopSelectRequest{}
	err := json.Unmarshal([]byte(input), &shopSelectReq)

	if err != nil {
		ctx.Response(http.StatusBadRequest, common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	authContext := models.AuthenticationContext{
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
// @Success		200	{array}	models.ShopUserInfo
// @Failure		401 {object}	common.ApiResponse
// @Security     AccessToken
// @Router /list-shop [get]
func (h AuthenticationHttp) ListShopCanAccess(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.shopUserService.ListShopByUser(authUsername, pageable)

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
// @Param		ShopFavoriteRequest  body      models.ShopFavoriteRequest  true  "Shop Favorite Request"
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.ApiResponse
// @Security     AccessToken
// @Router /favorite-shop [put]
func (h AuthenticationHttp) UpdateShopFavorite(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username

	input := ctx.ReadInput()

	reqBody := models.ShopFavoriteRequest{}
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
