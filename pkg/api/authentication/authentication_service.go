package authentication

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"

	micromodel "smlcloudplatform/internal/microservice/models"
)

type AuthenticationService struct {
	pst        microservice.IPersisterMongo
	jwtService *microservice.JwtService
}

func NewAuthenticationService(pst microservice.IPersisterMongo, jwtService *microservice.JwtService) *AuthenticationService {
	return &AuthenticationService{
		pst:        pst,
		jwtService: jwtService,
	}
}

func (svc *AuthenticationService) Login(userReq models.UserRequest) (string, error) {

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

	tokenString, err := svc.jwtService.GenerateTokenWithRedis(micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		// svc.ms.Log("Authentication service", err.Error())
		return "", errors.New("generate token error")
	}

	return tokenString, nil
}
