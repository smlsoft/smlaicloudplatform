package authentication

import (
	"context"
	"errors"
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/firebase"
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
	LoginWithFirebaseToken(token string) (string, error)
}

type AuthenticationService struct {
	authService           microservice.IAuthService
	authRepo              IAuthenticationMongoCacheRepository
	shopUserRepo          shop.IShopUserRepository
	shopUserAccessLogRepo shop.IShopUserAccessLogRepository
	passwordEncoder       func(string) (string, error)
	checkHashPassword     func(password string, hash string) bool
	timeNow               func() time.Time
	firebaseAdapter       firebase.IFirebaseAdapter
}

func NewAuthenticationService(
	authRepo IAuthenticationRepository,
	shopUserRepo shop.IShopUserRepository,
	shopUserAccessLogRepo shop.IShopUserAccessLogRepository,
	authService microservice.IAuthService,
	passwordEncoder func(string) (string, error),
	checkHashPassword func(password string, hash string) bool, timeNow func() time.Time,
	firebaseAdapter firebase.IFirebaseAdapter) AuthenticationService {
	return AuthenticationService{
		authRepo:              authRepo,
		authService:           authService,
		shopUserRepo:          shopUserRepo,
		shopUserAccessLogRepo: shopUserAccessLogRepo,
		passwordEncoder:       passwordEncoder,
		checkHashPassword:     checkHashPassword,
		timeNow:               timeNow,
		firebaseAdapter:       firebaseAdapter,
	}
}

func (svc AuthenticationService) normalizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func (svc AuthenticationService) Login(userLoginReq *models.UserLoginRequest, authContext AuthenticationContext) (string, error) {

	userLoginReq.Username = svc.normalizeUsername(userLoginReq.Username)

	userLoginReq.Username = strings.TrimSpace(userLoginReq.Username)
	userLoginReq.ShopID = strings.TrimSpace(userLoginReq.ShopID)

	findUser, err := svc.authRepo.FindUser(context.Background(), userLoginReq.Username)

	if err != nil && err.Error() != "mongo: no documents in result" {
		// svc.ms.Log("Authentication service", err.Error())
		return "", errors.New("auth: database connect error")
	}

	if len(findUser.Username) < 1 {
		return "", errors.New("username or password is invalid")
	}

	passwordInvalid := !svc.checkHashPassword(userLoginReq.Password, findUser.Password)

	if passwordInvalid {
		return "", errors.New("username or password is invalid")
	}

	tokenString, err := svc.authService.GenerateTokenWithRedis(microservice.AUTHTYPE_BEARER, micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		return "", errors.New("login failed")
	}

	if len(userLoginReq.ShopID) > 0 {
		shopUser, err := svc.shopUserRepo.FindByShopIDAndUsername(context.Background(), userLoginReq.ShopID, userLoginReq.Username)

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

		err = svc.shopUserRepo.UpdateLastAccess(context.Background(), userLoginReq.ShopID, userLoginReq.Username, lastAccessedAt)
		if err != nil {
			// implement log
			fmt.Println("error :: ", err.Error())
		}

		err = svc.shopUserAccessLogRepo.Create(context.Background(), models.ShopUserAccessLog{
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

	userRequest.Username = svc.normalizeUsername(userRequest.Username)

	userFind, err := svc.authRepo.FindUser(context.Background(), userRequest.Username)
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

	idx, err := svc.authRepo.CreateUser(context.Background(), user)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (svc AuthenticationService) Update(username string, userRequest models.UserProfileRequest) error {

	if username == "" {
		return errors.New("username invalid")
	}

	userFind, err := svc.authRepo.FindUser(context.Background(), username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.Username) < 1 {
		return errors.New("username is not exists")
	}

	userFind.Name = userRequest.Name
	userFind.UpdatedAt = svc.timeNow()

	err = svc.authRepo.UpdateUser(context.Background(), username, *userFind)

	if err != nil {
		return err
	}

	return nil
}

func (svc AuthenticationService) UpdatePassword(username string, currentPassword string, newPassword string) error {

	if username == "" {
		return errors.New("username invalid")
	}

	userFind, err := svc.authRepo.FindUser(context.Background(), username)
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

	err = svc.authRepo.UpdateUser(context.Background(), username, *userFind)

	if err != nil {
		return err
	}

	return nil
}

func (svc AuthenticationService) Logout(authorizationHeader string) error {
	return svc.authService.ExpireToken(microservice.AUTHTYPE_BEARER, authorizationHeader)
}

func (svc AuthenticationService) Profile(username string) (models.UserProfile, error) {
	stime := time.Now()
	userProfile := models.UserProfile{}
	user, err := svc.authRepo.FindUser(context.Background(), username)
	if err != nil {
		return userProfile, err
	}
	userProfile.Username = user.Username
	userProfile.Name = user.Name
	fmt.Println("total time :: ", time.Since(stime))
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

	shopUser, err := svc.shopUserRepo.FindByShopIDAndUsername(context.Background(), shopID, username)

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
	err = svc.shopUserRepo.UpdateLastAccess(context.Background(), shopID, username, lastAccessedAt)
	if err != nil {
		// implement log
		fmt.Println("error :: ", err.Error())
	}

	err = svc.shopUserAccessLogRepo.Create(
		context.Background(),
		models.ShopUserAccessLog{
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

	shopUser, err := svc.shopUserRepo.FindByShopIDAndUsername(context.Background(), shopID, username)

	if err != nil {
		return err
	}

	if shopUser.ID == primitive.NilObjectID {
		return errors.New("shop invalid")
	}

	err = svc.shopUserRepo.SaveFavorite(context.Background(), shopID, username, isFavorite)
	if err != nil {
		return errors.New("favorite failed")
	}

	return nil
}

func (svc AuthenticationService) LoginWithFirebaseToken(token string) (string, error) {

	userInfo, err := svc.firebaseAdapter.ValidateToken(token)
	if err != nil {
		return "", err
	}

	// find
	userFind, err := svc.authRepo.FindUser(context.Background(), userInfo.Email)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return "", err
	}

	if len(userFind.Username) == 0 {
		// register
		user := models.UserDoc{}

		user.Username = userInfo.Email
		user.Password = ""
		user.UserDetail.Name = userInfo.Name
		user.CreatedAt = svc.timeNow()

		_, err := svc.authRepo.CreateUser(context.Background(), user)
		if err != nil {
			return "", err
		}
		userFind, err = svc.authRepo.FindUser(context.Background(), userInfo.Email)
		if err != nil && err.Error() != "mongo: no documents in result" {
			return "", err
		}
	}

	tokenString, err := svc.authService.GenerateTokenWithRedis(microservice.AUTHTYPE_BEARER, micromodel.UserInfo{Username: userFind.Username, Name: userFind.Name})

	if err != nil {
		return "", errors.New("generate token error")
	}

	return tokenString, nil
}
