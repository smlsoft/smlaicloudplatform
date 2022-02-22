package authentication

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type AuthenticationHttp struct {
	ms                    *microservice.Microservice
	cfg                   microservice.IConfig
	authenticationService IAuthenticationService
}

func NewAuthenticationHttp(ms *microservice.Microservice, cfg microservice.IConfig) AuthenticationHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	authService := microservice.NewAuthService(ms.Cacher(cfg.CacherConfig()), 24*3)
	authRepository := NewAuthenticationRepository(pst)
	authenticationService := NewAuthenticationService(authRepository, authService)

	return AuthenticationHttp{
		ms:                    ms,
		cfg:                   cfg,
		authenticationService: authenticationService,
	}
}

func (h *AuthenticationHttp) RouteSetup() {

	h.ms.POST("/login", h.Login)
	h.ms.POST("/register", h.Register)
	h.ms.POST("/logout", h.Logout)
	h.ms.GET("/profile", h.Profile)
	h.ms.POST("/select-merchant", h.SelectMerchant)
}

// Login login
// @Description get struct array by ID
// @Tags		Authentication
// @Param		User  body      models.UserRequest  true  "Add account"
// @Accept 		json
// @Success		200	{object}	models.ApiResponse
// @Failure		401 {object}	models.ApiResponse
// @Router /authentication/login [post]
func (h *AuthenticationHttp) Login(ctx microservice.IServiceContext) error {

	input := ctx.ReadInput()

	userReq := &models.UserRequest{}
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
// @Success		200	{object}	models.ApiResponse
// @Accept 		json
// @Router		/authentication/register [post]
func (h *AuthenticationHttp) Register(ctx microservice.IServiceContext) error {
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

func (h *AuthenticationHttp) Logout(ctx microservice.IServiceContext) error {

	authorizationHeader := ctx.Header("Authorization")

	err := h.authenticationService.Logout(authorizationHeader)

	if err != nil {
		ctx.Response(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return err
	}

	return nil
}

func (h *AuthenticationHttp) Profile(ctx microservice.IServiceContext) error {

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

func (h *AuthenticationHttp) SelectMerchant(ctx microservice.IServiceContext) error {
	return nil
}
