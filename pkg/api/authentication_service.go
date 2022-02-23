package api

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	micromodel "smlcloudplatform/internal/microservice/models"
)

type AuthenticationService struct {
	ms          *microservice.Microservice
	cfg         microservice.IConfig
	authService *microservice.AuthService
}

func NewAuthenticationService(ms *microservice.Microservice, cfg microservice.IConfig) *AuthenticationService {
	// signKey, verifyKey, err := utils.LoadKey(cfg.SignKeyPath(), cfg.VerifyKeyPath())

	// if err != nil {
	// 	fmt.Println("jwt key error :: " + err.Error())
	// }

	// jwtService := microservice.NewJwtService(signKey, verifyKey, 60*24*10)

	authService := microservice.NewAuthService(ms.Cacher(cfg.CacherConfig()), 60*24*10)

	return &AuthenticationService{
		ms:          ms,
		cfg:         cfg,
		authService: authService,
	}
}

func (svc *AuthenticationService) RouteSetup() {

	cacher := svc.ms.Cacher(svc.cfg.CacherConfig())
	svc.ms.POST("/login", svc.Login)
	svc.ms.POST("/register", svc.Register)
	svc.ms.POST("/logout", svc.Logout)
	svc.ms.GET("/profile", svc.Profile, svc.authService.MWFuncWithRedis(cacher))
	svc.ms.POST("/select-merchant", svc.SelectMerchant, svc.authService.MWFuncWithMerchant(cacher))
}

func (svc *AuthenticationService) Login(ctx microservice.IServiceContext) error {

	input := ctx.ReadInput()

	userReq := &models.UserRequest{}
	err := json.Unmarshal([]byte(input), &userReq)

	if err != nil {
		ctx.ResponseError(400, "user payload invalid")
		return err
	}

	svc.ms.Logger.Info("Do Login Event")

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

	tokenString, err := svc.authService.GenerateTokenWithRedis(micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		svc.ms.Log("Authentication service", err.Error())
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "token": tokenString})

	return nil
}

func (svc *AuthenticationService) Register(ctx microservice.IServiceContext) error {
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
		ctx.ResponseError(400, err.Error())
		return err
	}

	if len(findUser.Username) > 0 {
		ctx.ResponseError(400, "username is exists.")
		return err
	}

	hashPassword, err := utils.HashPassword(userReq.Password)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	user := &models.User{
		Username:  userReq.Username,
		Password:  hashPassword,
		Name:      userReq.Name,
		CreatedAt: time.Now(),
	}

	idx, err := pst.Create(&models.User{}, user)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "id": idx})
	return nil
}

func (svc *AuthenticationService) Profile(ctx microservice.IServiceContext) error {

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true, "user": ctx.UserInfo()})

	return nil
}

func (svc *AuthenticationService) SelectMerchant(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	input := ctx.ReadInput()

	merchantSelectReq := &models.MerchantSelectRequest{}
	err := json.Unmarshal([]byte(input), &merchantSelectReq)

	if err != nil {
		ctx.ResponseError(400, "merchant payload invalid.")
		return err
	}

	tokenStr, err := svc.authService.GetTokenFromAuthorizationHeader(ctx.Header("Authorization"))

	if len(tokenStr) < 1 {
		ctx.ResponseError(400, "token invalid.")
		return err
	}

	pst := svc.ms.MongoPersister(svc.cfg.MongoPersisterConfig())

	merchantMember := &models.MerchantMember{}
	err = pst.FindOne(&models.MerchantMember{}, bson.M{"username": userInfo.Username, "merchantId": merchantSelectReq.MerchantId}, merchantMember)

	if merchantMember.Id == primitive.NilObjectID {
		ctx.ResponseError(400, "merchant invalid.")
		return err
	}

	err = svc.authService.SelectMerchant(tokenStr, merchantSelectReq.MerchantId)

	if err != nil {
		ctx.ResponseError(400, "failed merchant select.")
		return err
	}

	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})

	return nil
}

func (svc *AuthenticationService) Logout(ctx microservice.IServiceContext) error {

	svc.authService.ExpireToken(ctx.Header("Authorization"))
	ctx.Response(http.StatusOK, map[string]interface{}{"success": true})

	return nil
}
