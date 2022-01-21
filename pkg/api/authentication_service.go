package api

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"

	micromodel "smlcloudplatform/internal/microservice/models"
)

type AuthenticationService struct {
	ms         *microservice.Microservice
	cfg        microservice.IConfig
	jwtService *microservice.JwtService
	signKey    *rsa.PrivateKey
	verifyKey  *rsa.PublicKey
}

func NewAuthenticationService(ms *microservice.Microservice, cfg microservice.IConfig) *AuthenticationService {
	signBytes, err := ioutil.ReadFile("./../../private.key")

	if err != nil {
		ms.Log("auth", err.Error())
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err != nil {
		ms.Log("auth", err.Error())
	}

	verifyBytes, err := ioutil.ReadFile("./../../public.key")

	if err != nil {
		ms.Log("auth", err.Error())
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	if err != nil {
		ms.Log("auth", err.Error())
	}

	jwtService := microservice.NewJwtService(signKey, verifyKey, 60*24*10)
	return &AuthenticationService{
		ms:         ms,
		cfg:        cfg,
		jwtService: jwtService,
		signKey:    signKey,
		verifyKey:  verifyKey,
	}
}

func (svc *AuthenticationService) RouteSetup() {

	svc.ms.GET("/", svc.Index)
	svc.ms.POST("/login", svc.Login)
	svc.ms.POST("/register", svc.Register)
	svc.ms.GET("/profile", svc.Profile, svc.jwtService.MWFunc())
}

func (svc *AuthenticationService) Index(ctx microservice.IServiceContext) error {
	ctx.ResponseS(http.StatusOK, "ok")
	return nil
}

func (svc *AuthenticationService) Login(ctx microservice.IServiceContext) error {

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
		ctx.ResponseError(400, "find user error")
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

	tokenString, err := svc.jwtService.GenerateToken(micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name}, time.Duration(60*24*10)*time.Minute)
	if err != nil {
		svc.ms.Log("auth", err.Error())
		// ctx.ResponseError(http.StatusBadRequest, "can't create token.")
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
