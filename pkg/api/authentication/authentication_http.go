package authentication

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	micromodel "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type AuthenticationHttp struct {
	ms                    *microservice.Microservice
	cfg                   microservice.IConfig
	jwtService            *microservice.JwtService
	authenticationService *AuthenticationService
}

func NewAuthenticationHttp(ms *microservice.Microservice, cfg microservice.IConfig) *AuthenticationHttp {
	// signKey, verifyKey, err := utils.LoadKey(cfg.SignKeyPath(), cfg.VerifyKeyPath())

	// if err != nil {
	// 	fmt.Println("jwt key error :: " + err.Error())
	// }

	// jwtService := microservice.NewJwtService(signKey, verifyKey, 60*24*10)

	jwtService := microservice.NewJwtService(ms.Cacher(cfg.CacherConfig()), cfg.JwtSecretKey(), 60*24*10)

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	authenticationService := NewAuthenticationService(pst, jwtService)

	return &AuthenticationHttp{
		ms:                    ms,
		cfg:                   cfg,
		jwtService:            jwtService,
		authenticationService: authenticationService,
	}
}

// func (svc *AuthenticationHttp) RouteSetup() {

// 	svc.ms.GET("/", svc.Index)
// 	svc.ms.POST("/login", svc.Login)
// 	svc.ms.POST("/register", svc.Register)
// 	svc.ms.POST("/logout", svc.Logout)
// 	svc.ms.GET("/profile", svc.Profile, svc.jwtService.MWFunc())
// }

func (svc *AuthenticationHttp) Login(ctx microservice.IServiceContext) error {

	input := ctx.ReadInput()

	userReq := &models.UserRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	findUser := &models.User{}
	err = pst.FindOne(&models.User{}, bson.M{"username": userReq.Username}, findUser)

	if err != nil && err.Error() != "mongo: no documents in result" {
		svc.ms.Log("Authentication service", err.Error())
		ctx.ResponseError(400, "database error")
		return err
	}

	if len(findUser.Username) < 1 {
		ctx.ResponseError(400, "username is not exists")
		return err
	}

	passwordInvalid := !utils.CheckPasswordHash(userReq.Password, findUser.Password)

	if passwordInvalid {
		ctx.ResponseError(400, "password is not invalid")
		return err
	}

	tokenString, err := svc.jwtService.GenerateTokenWithRedis(micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		svc.ms.Log("auth", err.Error())
		// ctx.ResponseError(http.StatusBadRequest, "can't create token.")
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "token": tokenString})

	return nil
}
