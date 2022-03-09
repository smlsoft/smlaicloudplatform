package authentication

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/merchant"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"

	micromodel "smlcloudplatform/internal/microservice/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationService interface {
	Login(userReq *models.UserLoginRequest) (string, error)
	Register(userRequest models.UserRequest) (string, error)
	Update(username string, userRequest models.UserRequest) error
	UpdatePassword(username string, currentPassword string, newPassword string) error
	Logout(authorizationHeader string) error
	Profile(username string) (models.UserProfile, error)
	AccessMerchant(authorizationHeader string, merchantId string, username string) error
}

type AuthenticationService struct {
	authService      *microservice.AuthService
	authRepo         IAuthenticationRepository
	merchantUserRepo merchant.IMerchantUserRepository
}

func NewAuthenticationService(authRepo IAuthenticationRepository, merchantUserRepo merchant.IMerchantUserRepository, authService *microservice.AuthService) IAuthenticationService {
	return AuthenticationService{
		authRepo:         authRepo,
		authService:      authService,
		merchantUserRepo: merchantUserRepo,
	}
}

func (svc AuthenticationService) Login(userLoginReq *models.UserLoginRequest) (string, error) {

	userLoginReq.Username = strings.TrimSpace(userLoginReq.Username)
	userLoginReq.MerchantId = strings.TrimSpace(userLoginReq.MerchantId)

	findUser, err := svc.authRepo.FindUser(userLoginReq.Username)

	if err != nil && err.Error() != "mongo: no documents in result" {
		// svc.ms.Log("Authentication service", err.Error())
		return "", errors.New("auth: database connect error")
	}

	if len(findUser.Username) < 1 {
		return "", errors.New("username is not exists")
	}

	passwordInvalid := !utils.CheckPasswordHash(userLoginReq.Password, findUser.Password)

	if passwordInvalid {
		return "", errors.New("password is not invalid")
	}

	tokenString, err := svc.authService.GenerateTokenWithRedis(micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		return "", errors.New("generate token error")
	}

	if len(userLoginReq.MerchantId) > 0 {
		merchantUser, err := svc.merchantUserRepo.FindByMerchantIdAndUsername(userLoginReq.MerchantId, userLoginReq.Username)

		if err != nil {
			return "", err
		}

		if merchantUser.Id == primitive.NilObjectID {
			return "", errors.New("merchant invalid")
		}

		err = svc.authService.SelectMerchant(tokenString, userLoginReq.MerchantId, string(merchantUser.Role))

		if err != nil {
			return "", errors.New("failed merchant select")
		}
	}

	return tokenString, nil
}

func (svc AuthenticationService) Register(userRequest models.UserRequest) (string, error) {

	userFind, err := svc.authRepo.FindUser(userRequest.Username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return "", err
	}

	if len(userFind.Username) > 0 {
		return "", errors.New("username is exists")
	}

	hashPassword, err := utils.HashPassword(userRequest.Password)

	if err != nil {
		return "", err
	}

	user := models.User{
		Username:  userRequest.Username,
		Password:  hashPassword,
		Name:      userRequest.Name,
		CreatedAt: time.Now(),
	}

	idx, err := svc.authRepo.CreateUser(user)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (svc AuthenticationService) Update(username string, userRequest models.UserRequest) error {

	if username == "" {
		return errors.New("username invalid")
	}

	userFind, err := svc.authRepo.FindUser(username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.Username) < 1 {
		return errors.New("username is not exists")
	}

	userFind.Name = userRequest.Name
	userFind.UpdatedAt = time.Now()

	err = svc.authRepo.UpdateUser(username, *userFind)

	if err != nil {
		return err
	}

	return nil
}

func (svc AuthenticationService) UpdatePassword(username string, currentPassword string, newPassword string) error {

	if username == "" {
		return errors.New("username invalid")
	}

	userFind, err := svc.authRepo.FindUser(username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.Username) < 1 {
		return errors.New("username is not exists")
	}

	passwordNotMatch := !utils.CheckPasswordHash(currentPassword, userFind.Password)

	if passwordNotMatch {
		return errors.New("current password invalid")
	}

	hashPassword, err := utils.HashPassword(newPassword)

	if err != nil {
		return err
	}

	userFind.Password = hashPassword
	userFind.UpdatedAt = time.Now()

	err = svc.authRepo.UpdateUser(username, *userFind)

	if err != nil {
		return err
	}

	return nil
}

func (svc AuthenticationService) Logout(authorizationHeader string) error {
	return svc.authService.ExpireToken(authorizationHeader)
}

func (svc AuthenticationService) Profile(username string) (models.UserProfile, error) {
	userProfile := models.UserProfile{}
	user, err := svc.authRepo.FindUser(username)
	if err != nil {
		return userProfile, err
	}
	userProfile.Username = user.Username
	userProfile.Name = user.Name
	return userProfile, nil
}

func (svc AuthenticationService) AccessMerchant(authorizationHeader string, merchantId string, username string) error {

	if merchantId == "" {
		return errors.New("merchant invalid")
	}

	if username == "" {
		return errors.New("username invalid")
	}

	tokenStr, err := svc.authService.GetTokenFromAuthorizationHeader(authorizationHeader)

	if err != nil {
		return err
	}

	if len(tokenStr) < 1 {
		return errors.New("token invalid")
	}

	merchantUser, err := svc.merchantUserRepo.FindByMerchantIdAndUsername(merchantId, username)

	if err != nil {
		return err
	}

	if merchantUser.Id == primitive.NilObjectID {
		return errors.New("merchant invalid")
	}

	err = svc.authService.SelectMerchant(tokenStr, merchantId, string(merchantUser.Role))

	if err != nil {
		return errors.New("failed merchant select")
	}
	return nil
}
