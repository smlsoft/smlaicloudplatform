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
	authenticationService *AuthenticationService
}

func NewAuthenticationHttp(ms *microservice.Microservice, cfg microservice.IConfig) *AuthenticationHttp {

	authService := microservice.NewAuthService(ms.Cacher(cfg.CacherConfig()), 24*3)
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	authenticationService := NewAuthenticationService(pst, authService)

	return &AuthenticationHttp{
		ms:                    ms,
		cfg:                   cfg,
		authenticationService: authenticationService,
	}
}

func (h *AuthenticationHttp) RouteSetup() {

	h.ms.POST("/login", h.Login)
	// h.ms.POST("/register", h.Register)
	// h.ms.POST("/logout", h.Logout)
	// h.ms.GET("/profile", h.Profile, h.jwtService.MWFunc())
}

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
