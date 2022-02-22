package authentication

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"

	micromodel "smlcloudplatform/internal/microservice/models"
)

type IAuthenticationService interface {
	Login(userReq *models.UserRequest) (string, error)
	Register() error
	Logout() error
	Profile() error
	AccessMerchant() error
}

type AuthenticationService struct {
	pst         microservice.IPersisterMongo
	authService *microservice.AuthService
}

func NewAuthenticationService(pst microservice.IPersisterMongo, authService *microservice.AuthService) AuthenticationService {
	return AuthenticationService{
		pst:         pst,
		authService: authService,
	}
}

/* imprement a service */

func (svc AuthenticationService) Login(userReq *models.UserRequest) (string, error) {

	findUser := &models.User{}
	err := svc.pst.FindOne(&models.User{}, bson.M{"username": userReq.Username}, findUser)

	if err != nil && err.Error() != "mongo: no documents in result" {
		// svc.ms.Log("Authentication service", err.Error())
		return "", errors.New("auth: database connect error")
	}

	if len(findUser.Username) < 1 {
		return "", errors.New("username is not exists")
	}

	passwordInvalid := !utils.CheckPasswordHash(userReq.Password, findUser.Password)

	if passwordInvalid {
		return "", errors.New("password is not invalid")
	}

	tokenString, err := svc.authService.GenerateTokenWithRedis(micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		// svc.ms.Log("Authentication service", err.Error())
		return "", errors.New("generate token error")
	}

	return tokenString, nil
}

func (svc AuthenticationService) Register() error {
	return nil
}

func (svc AuthenticationService) Logout() error {
	return nil
}

func (svc AuthenticationService) Profile() error {
	return nil
}

func (svc AuthenticationService) AccessMerchant() error {
	return nil
}
