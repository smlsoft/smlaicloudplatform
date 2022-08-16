package authentication

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shop"
	"smlcloudplatform/pkg/shop/models"
	"strings"
	"time"

	micromodel "smlcloudplatform/internal/microservice/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationService interface {
	Login(userReq *models.UserLoginRequest) (string, error)
	Register(userRequest models.UserRequest) (string, error)
	Update(username string, userRequest models.UserProfileRequest) error
	UpdatePassword(username string, currentPassword string, newPassword string) error
	Logout(authorizationHeader string) error
	Profile(username string) (models.UserProfile, error)
	AccessShop(shopID string, username string, authorizationHeader string) error
}

type AuthenticationService struct {
	authService       microservice.IAuthService
	authRepo          IAuthenticationRepository
	shopUserRepo      shop.IShopUserRepository
	passwordEncoder   func(string) (string, error)
	checkHashPassword func(password string, hash string) bool
	timeNow           func() time.Time
}

func NewAuthenticationService(authRepo IAuthenticationRepository, shopUserRepo shop.IShopUserRepository, authService microservice.IAuthService, passwordEncoder func(string) (string, error), checkHashPassword func(password string, hash string) bool, timeNow func() time.Time) AuthenticationService {
	return AuthenticationService{
		authRepo:          authRepo,
		authService:       authService,
		shopUserRepo:      shopUserRepo,
		passwordEncoder:   passwordEncoder,
		checkHashPassword: checkHashPassword,
		timeNow:           timeNow,
	}
}

func (svc AuthenticationService) Login(userLoginReq *models.UserLoginRequest) (string, error) {

	userLoginReq.Username = strings.TrimSpace(userLoginReq.Username)
	userLoginReq.ShopID = strings.TrimSpace(userLoginReq.ShopID)

	findUser, err := svc.authRepo.FindUser(userLoginReq.Username)

	if err != nil && err.Error() != "mongo: no documents in result" {
		// svc.ms.Log("Authentication service", err.Error())
		return "", errors.New("auth: database connect error")
	}

	if len(findUser.Username) < 1 {
		return "", errors.New("username is not exists")
	}

	passwordInvalid := !svc.checkHashPassword(userLoginReq.Password, findUser.Password)

	if passwordInvalid {
		return "", errors.New("password is not invalid")
	}

	tokenString, err := svc.authService.GenerateTokenWithRedis(micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		return "", errors.New("generate token error")
	}

	if len(userLoginReq.ShopID) > 0 {
		shopUser, err := svc.shopUserRepo.FindByShopIDAndUsername(userLoginReq.ShopID, userLoginReq.Username)

		if err != nil {
			return "", err
		}

		if shopUser.ID == primitive.NilObjectID {
			return "", errors.New("shop invalid")
		}

		err = svc.authService.SelectShop(tokenString, userLoginReq.ShopID, shopUser.Role)

		if err != nil {
			return "", errors.New("failed shop select")
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

	hashPassword, err := svc.passwordEncoder(userRequest.Password)

	if err != nil {
		return "", err
	}

	user := models.UserDoc{}

	user.Username = userRequest.Username
	user.Password = hashPassword
	user.UserDetail = userRequest.UserDetail
	user.CreatedAt = svc.timeNow()

	idx, err := svc.authRepo.CreateUser(user)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (svc AuthenticationService) Update(username string, userRequest models.UserProfileRequest) error {

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
	userFind.UpdatedAt = svc.timeNow()

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

	passwordInvalid := !svc.checkHashPassword(currentPassword, userFind.Password)

	if passwordInvalid {
		return errors.New("current password invalid")
	}

	hashPassword, err := svc.passwordEncoder(newPassword)

	if err != nil {
		return err
	}

	userFind.Password = hashPassword
	userFind.UpdatedAt = svc.timeNow()

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

func (svc AuthenticationService) AccessShop(shopID string, username string, authorizationHeader string) error {

	if shopID == "" {
		return errors.New("shop invalid")
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

	shopUser, err := svc.shopUserRepo.FindByShopIDAndUsername(shopID, username)

	if err != nil {
		return err
	}

	if shopUser.ID == primitive.NilObjectID {
		return errors.New("shop invalid")
	}

	err = svc.authService.SelectShop(tokenStr, shopID, shopUser.Role)

	if err != nil {
		return errors.New("failed shop select")
	}
	return nil
}
