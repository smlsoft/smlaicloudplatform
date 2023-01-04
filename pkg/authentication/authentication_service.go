package authentication

import (
	"errors"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shop"
	"smlcloudplatform/pkg/shop/models"
	"strings"
	"time"

	micromodel "smlcloudplatform/internal/microservice/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationService interface {
	Login(userReq *models.UserLoginRequest, authContext AuthenticationContext) (string, error)
	Register(userRequest models.UserRequest) (string, error)
	Update(username string, userRequest models.UserProfileRequest) error
	UpdatePassword(username string, currentPassword string, newPassword string) error
	Logout(authorizationHeader string) error
	Profile(username string) (models.UserProfile, error)
	AccessShop(shopID string, username string, authorizationHeader string, authContext AuthenticationContext) error
	UpdateFavoriteShop(shopID string, username string, isFavorite bool) error
}

type AuthenticationService struct {
	authService           microservice.IAuthService
	authRepo              IAuthenticationRepository
	shopUserRepo          shop.IShopUserRepository
	shopUserAccessLogRepo shop.IShopUserAccessLogRepository
	passwordEncoder       func(string) (string, error)
	checkHashPassword     func(password string, hash string) bool
	timeNow               func() time.Time
}

func NewAuthenticationService(authRepo IAuthenticationRepository, shopUserRepo shop.IShopUserRepository, shopUserAccessLogRepo shop.IShopUserAccessLogRepository, authService microservice.IAuthService, passwordEncoder func(string) (string, error), checkHashPassword func(password string, hash string) bool, timeNow func() time.Time) AuthenticationService {
	return AuthenticationService{
		authRepo:              authRepo,
		authService:           authService,
		shopUserRepo:          shopUserRepo,
		shopUserAccessLogRepo: shopUserAccessLogRepo,
		passwordEncoder:       passwordEncoder,
		checkHashPassword:     checkHashPassword,
		timeNow:               timeNow,
	}
}

func (svc AuthenticationService) Login(userLoginReq *models.UserLoginRequest, authContext AuthenticationContext) (string, error) {

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

	tokenString, err := svc.authService.GenerateTokenWithRedis(microservice.AUTHTYPE_BEARER, micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

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

		err = svc.authService.SelectShop(microservice.AUTHTYPE_BEARER, tokenString, userLoginReq.ShopID, shopUser.Role)

		if err != nil {
			return "", errors.New("failed shop select")
		}

		lastAccessedAt := svc.timeNow()

		err = svc.shopUserRepo.UpdateLastAccess(userLoginReq.ShopID, userLoginReq.Username, lastAccessedAt)
		if err != nil {
			// implement log
			fmt.Println("error :: ", err.Error())
		}

		err = svc.shopUserAccessLogRepo.Create(models.ShopUserAccessLog{
			ShopID:         userLoginReq.ShopID,
			Username:       userLoginReq.Username,
			Ip:             authContext.Ip,
			LastAccessedAt: lastAccessedAt,
		})

		if err != nil {
			// implement log
			fmt.Println("error :: ", err.Error())
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
	return svc.authService.ExpireToken(microservice.AUTHTYPE_BEARER, authorizationHeader)
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

func (svc AuthenticationService) AccessShop(shopID string, username string, authorizationHeader string, authContext AuthenticationContext) error {

	if shopID == "" {
		return errors.New("shop invalid")
	}

	if username == "" {
		return errors.New("username invalid")
	}

	tokenStr, err := svc.authService.GetTokenFromAuthorizationHeader(microservice.AUTHTYPE_BEARER, authorizationHeader)

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

	err = svc.authService.SelectShop(microservice.AUTHTYPE_BEARER, tokenStr, shopID, shopUser.Role)

	if err != nil {
		return errors.New("failed shop select")
	}

	lastAccessedAt := svc.timeNow()
	err = svc.shopUserRepo.UpdateLastAccess(shopID, username, lastAccessedAt)
	if err != nil {
		// implement log
		fmt.Println("error :: ", err.Error())
	}

	err = svc.shopUserAccessLogRepo.Create(models.ShopUserAccessLog{
		ShopID:         shopID,
		Username:       username,
		Ip:             authContext.Ip,
		LastAccessedAt: lastAccessedAt,
	})

	if err != nil {
		// implement log
		fmt.Println("error :: ", err.Error())
	}

	return nil
}

func (svc AuthenticationService) UpdateFavoriteShop(shopID string, username string, isFavorite bool) error {

	if shopID == "" {
		return errors.New("shop invalid")
	}

	if username == "" {
		return errors.New("username invalid")
	}

	shopUser, err := svc.shopUserRepo.FindByShopIDAndUsername(shopID, username)

	if err != nil {
		return err
	}

	if shopUser.ID == primitive.NilObjectID {
		return errors.New("shop invalid")
	}

	err = svc.shopUserRepo.SaveFavorite(shopID, username, isFavorite)
	if err != nil {
		return errors.New("favorite failed")
	}

	return nil
}
