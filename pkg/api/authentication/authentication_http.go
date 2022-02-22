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
	authService           microservice.AuthService
	authenticationService IAuthenticationService
}

func NewAuthenticationHttp(ms *microservice.Microservice, cfg microservice.IConfig) AuthenticationHttp {

	authService := microservice.NewAuthService(ms.Cacher(cfg.CacherConfig()), 24*3)
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	authenticationService := NewAuthenticationService(pst, authService)

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
// @Router /login [post]
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

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "token": tokenString})

	return nil
}

// RegisterMember register
// @Summary		Register An Account
// @Description	For User Register Application
// @Tags		Authentication
// @Param		User  body      models.UserRequest  true  "Add account"
// @Success		200	{object}	models.ApiResponse
// @Accept 		json
// @Router		/register [post]
func (h *AuthenticationHttp) Register(ctx microservice.IServiceContext) error {
	return nil
}

func (h *AuthenticationHttp) Logout(ctx microservice.IServiceContext) error {
	return nil
}

func (h *AuthenticationHttp) Profile(ctx microservice.IServiceContext) error {
	return nil
}

func (h *AuthenticationHttp) SelectMerchant(ctx microservice.IServiceContext) error {
	return nil
}
