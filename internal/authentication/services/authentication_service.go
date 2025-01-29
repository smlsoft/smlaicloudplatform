package services

import (
	"context"
	"errors"
	"fmt"
	"smlcloudplatform/internal/authentication/models"
	auth_models "smlcloudplatform/internal/authentication/models"
	"smlcloudplatform/internal/authentication/repositories"
	"smlcloudplatform/internal/firebase"
	"smlcloudplatform/internal/logger"
	"smlcloudplatform/internal/shop"
	"smlcloudplatform/internal/utils"
	"smlcloudplatform/pkg/microservice"
	"strings"
	"time"

	micromodel "smlcloudplatform/pkg/microservice/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationService interface {
	LoginWithPhoneNumber(userLoginReq *auth_models.UserLoginPhoneNumberRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error)
	LoginWithPhoneNumberOTP(userLoginReq *auth_models.PhoneNumberOTPRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error)
	Login(userReq *auth_models.UserLoginRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error)
	Poslogin(userReq *auth_models.PosLoginRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error)
	LoginEmail(userReq *auth_models.PosLoginRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error)
	Register(userRequest auth_models.RegisterEmailRequest) (string, error)
	ForgotPasswordByPhonenumber(userRequest auth_models.ForgotPasswordPhoneNumberRequest) error
	Update(username string, userRequest auth_models.UserProfileRequest) error
	UpdatePassword(username string, currentPassword string, newPassword string) error
	Logout(authorizationHeader string) error
	Profile(username string) (auth_models.UserProfile, error)
	AccessShop(shopID string, username string, authorizationHeader string, authContext models.AuthenticationContext) error
	UpdateFavoriteShop(shopID string, username string, isFavorite bool) error
	LoginWithFirebaseToken(token string) (string, error)
	RefreshToken(tokenRequest models.TokenLoginRequest) (models.TokenLoginResponse, error)

	CheckExistsUsername(username string) (bool, error)
	CheckExistsPhonenumber(phoneNumber string) (bool, error)
	SendPhonenumberOTP(otpRequest auth_models.OTPRequest) (auth_models.OTPResponse, error)
	RegisterByPhonenumber(userRequest auth_models.RegisterPhoneNumberRequest) (string, error)
}

type AuthenticationService struct {
	authService           microservice.IAuthService
	authRepo              repositories.IAuthenticationMongoCacheRepository
	shopUserRepo          shop.IShopUserRepository
	shopUserAccessLogRepo shop.IShopUserAccessLogRepository
	smsRepo               repositories.IAuthenticationSMSRepository
	randdomString         func(int) string
	randdomNumber         func(int) string
	generateGUID          func() string
	passwordEncoder       func(string) (string, error)
	checkHashPassword     func(password string, hash string) bool
	timeNow               func() time.Time
	firebaseAdapter       firebase.IFirebaseAdapter
}

func NewAuthenticationService(
	authRepo repositories.IAuthenticationRepository,
	shopUserRepo shop.IShopUserRepository,
	shopUserAccessLogRepo shop.IShopUserAccessLogRepository,
	smsRepo repositories.IAuthenticationSMSRepository,
	authService microservice.IAuthService,
	randdomString func(int) string,
	randdomNumber func(int) string,
	generateGUID func() string,
	passwordEncoder func(string) (string, error),
	checkHashPassword func(password string, hash string) bool,
	timeNow func() time.Time,
	firebaseAdapter firebase.IFirebaseAdapter) AuthenticationService {
	return AuthenticationService{
		authRepo:              authRepo,
		authService:           authService,
		shopUserRepo:          shopUserRepo,
		shopUserAccessLogRepo: shopUserAccessLogRepo,
		smsRepo:               smsRepo,
		randdomString:         randdomString,
		randdomNumber:         randdomNumber,
		generateGUID:          generateGUID,
		passwordEncoder:       passwordEncoder,
		checkHashPassword:     checkHashPassword,
		timeNow:               timeNow,
		firebaseAdapter:       firebaseAdapter,
	}
}

func (svc AuthenticationService) ValidateOTP(refCode, OTP string) (bool, error) {
	return svc.smsRepo.VerifyOTP(refCode, OTP)
}

func (svc AuthenticationService) LoginWithPhoneNumberOTP(userLoginReq *auth_models.PhoneNumberOTPRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error) {

	isOTPPassed, err := svc.ValidateOTP(userLoginReq.RefCode, userLoginReq.OTP)

	if err != nil {
		return models.TokenLoginResponse{}, errors.New("OTP invalid")
	}

	if !isOTPPassed {
		return models.TokenLoginResponse{}, errors.New("OTP invalid")
	}

	findUser, err := svc.authRepo.FindByIdentity(context.Background(), "phonenumber", userLoginReq.PhoneNumber)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.TokenLoginResponse{}, errors.New("auth: database connect error")
	}

	if len(findUser.PhoneNumber) < 1 {
		return models.TokenLoginResponse{}, errors.New("username or password is invalid")
	}

	tokenString, err := svc.authService.GenerateTokenWithRedis(microservice.AUTHTYPE_BEARER, micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		return models.TokenLoginResponse{}, errors.New("login failed")
	}

	refreshTokenString, err := svc.authService.GenerateTokenWithRedis(microservice.AUTHTYPE_REFRESH, micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		svc.authService.DeleteToken(microservice.AUTHTYPE_BEARER, tokenString)
		return models.TokenLoginResponse{}, errors.New("login failed")
	}

	return models.TokenLoginResponse{
		Token:   tokenString,
		Refresh: refreshTokenString,
	}, nil
}

func (svc AuthenticationService) LoginWithPhoneNumber(userLoginReq *auth_models.UserLoginPhoneNumberRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error) {

	findUser, err := svc.authRepo.FindByPhonenumber(context.Background(), userLoginReq.PhoneNumberField)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.TokenLoginResponse{}, errors.New("auth: database connect error")
	}

	if len(findUser.PhoneNumber) < 1 {
		return models.TokenLoginResponse{}, errors.New("username or password is invalid")
	}

	passwordInvalid := !svc.checkHashPassword(userLoginReq.Password, findUser.Password)

	if passwordInvalid {
		return models.TokenLoginResponse{}, errors.New("username or password is invalid")
	}

	resultLogin, err := svc.processUserLogin(*findUser, userLoginReq.ShopID, authContext)

	if err != nil {
		return models.TokenLoginResponse{}, err
	}

	return resultLogin, nil
}

func (svc AuthenticationService) Login(userLoginReq *auth_models.UserLoginRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error) {

	userLoginReq.Username = utils.NormalizeUsername(userLoginReq.Username)

	userLoginReq.Username = strings.TrimSpace(userLoginReq.Username)
	userLoginReq.ShopID = strings.TrimSpace(userLoginReq.ShopID)

	findUser, err := svc.authRepo.FindUser(context.Background(), userLoginReq.Username)

	if err != nil && err.Error() != "mongo: no documents in result" {
		// svc.ms.Log("Authentication service", err.Error())
		return models.TokenLoginResponse{}, errors.New("auth: database connect error")
	}

	if len(findUser.Username) < 1 {
		return models.TokenLoginResponse{}, errors.New("username or password is invalid")
	}

	passwordInvalid := !svc.checkHashPassword(userLoginReq.Password, findUser.Password)

	if passwordInvalid {
		return models.TokenLoginResponse{}, errors.New("username or password is invalid")
	}

	resultLogin, err := svc.processUserLogin(*findUser, userLoginReq.ShopID, authContext)

	if err != nil {
		return models.TokenLoginResponse{}, err
	}

	return resultLogin, nil
}

func (svc AuthenticationService) Poslogin(userLoginReq *auth_models.PosLoginRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error) {

	userLoginReq.Username = utils.NormalizeUsername(userLoginReq.Username)

	userLoginReq.Username = strings.TrimSpace(userLoginReq.Username)
	userLoginReq.ShopID = strings.TrimSpace(userLoginReq.ShopID)

	findUser, err := svc.authRepo.FindUser(context.Background(), userLoginReq.Username)

	if err != nil && err.Error() != "mongo: no documents in result" {
		// svc.ms.Log("Authentication service", err.Error())
		return models.TokenLoginResponse{}, errors.New("auth: database connect error")
	}

	if len(findUser.Username) < 1 {
		return models.TokenLoginResponse{}, errors.New("username or password is invalid")
	}

	// passwordInvalid := !svc.checkHashPassword(userLoginReq.Password, findUser.Password)

	// if passwordInvalid {
	// 	return models.TokenLoginResponse{}, errors.New("username or password is invalid")
	// }

	resultLogin, err := svc.processUserLogin(*findUser, userLoginReq.ShopID, authContext)

	if err != nil {
		return models.TokenLoginResponse{}, err
	}

	return resultLogin, nil
}

func (svc AuthenticationService) LoginEmail(userLoginReq *auth_models.PosLoginRequest, authContext models.AuthenticationContext) (models.TokenLoginResponse, error) {

	userLoginReq.Username = utils.NormalizeUsername(userLoginReq.Username)

	userLoginReq.Username = strings.TrimSpace(userLoginReq.Username)
	userLoginReq.ShopID = strings.TrimSpace(userLoginReq.ShopID)

	findUser, err := svc.authRepo.FindUser(context.Background(), userLoginReq.Username)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.TokenLoginResponse{}, errors.New("auth: database connect error")
	}

	if len(findUser.Username) == 0 {
		// Register user if not found
		user := auth_models.UserDoc{}
		user.Username = userLoginReq.Username
		user.Password = ""
		user.UserDetail.Name = userLoginReq.Username
		user.CreatedAt = svc.timeNow()

		_, err := svc.authRepo.CreateUser(context.Background(), user)
		if err != nil {
			return models.TokenLoginResponse{}, err
		}

		findUser, err = svc.authRepo.FindUser(context.Background(), userLoginReq.Username)
		if err != nil && err.Error() != "mongo: no documents in result" {
			return models.TokenLoginResponse{}, err
		}
	}

	resultLogin, err := svc.processUserLogin(*findUser, userLoginReq.ShopID, authContext)

	if err != nil {
		return models.TokenLoginResponse{}, err
	}

	return resultLogin, nil
}

func (svc *AuthenticationService) processUserLogin(findUser auth_models.UserDoc, shopID string, authContext models.AuthenticationContext) (models.TokenLoginResponse, error) {
	tokenString, err := svc.authService.GenerateTokenWithRedis(microservice.AUTHTYPE_BEARER, micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		return models.TokenLoginResponse{}, errors.New("login failed")
	}

	refreshTokenString, err := svc.authService.GenerateTokenWithRedis(microservice.AUTHTYPE_REFRESH, micromodel.UserInfo{Username: findUser.Username, Name: findUser.Name})

	if err != nil {
		svc.authService.DeleteToken(microservice.AUTHTYPE_BEARER, tokenString)
		return models.TokenLoginResponse{}, errors.New("login failed")
	}

	if len(shopID) > 0 {
		shopUser, err := svc.shopUserRepo.FindByShopIDAndUsername(context.Background(), shopID, findUser.Username)

		if err != nil {
			return models.TokenLoginResponse{}, err
		}

		if shopUser.ID == primitive.NilObjectID {
			return models.TokenLoginResponse{}, errors.New("shop invalid")
		}

		err = svc.authService.SelectShop(microservice.AUTHTYPE_BEARER, tokenString, shopID, shopUser.Role)

		if err != nil {
			return models.TokenLoginResponse{}, errors.New("failed shop select")
		}

		lastAccessedAt := svc.timeNow()

		err = svc.shopUserRepo.UpdateLastAccess(context.Background(), shopID, findUser.Username, lastAccessedAt)
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}

		err = svc.shopUserAccessLogRepo.Create(context.Background(), auth_models.ShopUserAccessLog{
			ShopID:         shopID,
			Username:       findUser.Username,
			Ip:             authContext.Ip,
			LastAccessedAt: lastAccessedAt,
		})

		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}

	return models.TokenLoginResponse{Token: tokenString, Refresh: refreshTokenString}, nil
}

func (svc AuthenticationService) RefreshToken(tokenRequest models.TokenLoginRequest) (models.TokenLoginResponse, error) {

	token, refreshToken, err := svc.authService.RefreshToken(tokenRequest.Token)

	if err != nil {
		return models.TokenLoginResponse{}, err
	}

	return models.TokenLoginResponse{
		Token:   token,
		Refresh: refreshToken,
	}, nil
}

func (svc AuthenticationService) Register(userEmailRequest auth_models.RegisterEmailRequest) (string, error) {

	userEmailRequest.Email = utils.NormalizeEmail(userEmailRequest.Email)

	userFind, err := svc.authRepo.FindByIdentity(context.Background(), "email", userEmailRequest.Email)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return "", err
	}

	if len(userFind.Username) > 0 {
		return "", errors.New("username is exists")
	}

	hashPassword, err := svc.passwordEncoder(userEmailRequest.Password)

	if err != nil {
		return "", err
	}

	user := auth_models.UserDoc{}

	user.UserDetail = userEmailRequest.UserDetail

	user.UID = svc.generateGUID()
	user.Username = userEmailRequest.Email
	user.Email = userEmailRequest.Email
	user.Password = hashPassword
	user.CreatedAt = svc.timeNow()

	idx, err := svc.authRepo.CreateUser(context.Background(), user)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (svc AuthenticationService) CheckExistsUsername(username string) (bool, error) {

	username = utils.NormalizeUsername(username)

	userFind, err := svc.authRepo.FindByIdentity(context.Background(), "username", username)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return true, err
	}

	if len(userFind.Username) > 0 {
		return true, nil
	}

	return false, nil
}

func (svc AuthenticationService) CheckExistsPhonenumber(phoneNumber string) (bool, error) {

	phoneNumber = utils.NormalizePhonenumber(phoneNumber)

	userPhonenumberFind, err := svc.authRepo.FindByIdentity(context.Background(), "phonenumber", phoneNumber)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return true, err
	}

	if len(userPhonenumberFind.PhoneNumber) > 0 {
		return true, nil
	}

	return false, nil
}

func (svc AuthenticationService) SendPhonenumberOTP(otpRequest auth_models.OTPRequest) (auth_models.OTPResponse, error) {

	otpRequest.PhoneNumber = utils.NormalizePhonenumber(otpRequest.PhoneNumber)

	fullPhoneNumber := fmt.Sprintf("%s%s", otpRequest.CountryCode, otpRequest.PhoneNumber)
	result, err := svc.smsRepo.SendOTPViaLink(fullPhoneNumber)

	if err != nil {
		return auth_models.OTPResponse{}, err
	}

	return result, nil
}

func (svc AuthenticationService) RegisterByPhonenumber(userRequest auth_models.RegisterPhoneNumberRequest) (string, error) {

	isOtpPass, err := svc.smsRepo.VerifyOTPViaLink(userRequest.OTPToken, userRequest.OTPRefCode, userRequest.OTPPin)

	if err != nil {
		return "", err
	}

	if !isOtpPass {
		return "", errors.New("otp invalid")
	}

	userRequest.PhoneNumber = utils.NormalizePhonenumber(userRequest.PhoneNumber)

	if exists, err := svc.CheckExistsUsername(userRequest.Username); err != nil {
		return "", err
	} else if exists {
		return "", errors.New("username is exists")
	}

	if exists, err := svc.CheckExistsPhonenumber(userRequest.PhoneNumber); err != nil {
		return "", err
	} else if exists {
		return "", errors.New("phonenumber is exists")
	}

	hashPassword, err := svc.passwordEncoder(userRequest.Password)

	if err != nil {
		return "", err
	}

	user := auth_models.UserDoc{}

	user.UserDetail = userRequest.UserDetail

	user.UID = svc.generateGUID()
	user.Username = userRequest.Username
	user.Email = ""
	user.Password = hashPassword
	user.PhoneNumber = userRequest.PhoneNumber
	user.RegisterType = "phonenumber"

	user.CreatedAt = svc.timeNow()

	idx, err := svc.authRepo.CreateUser(context.Background(), user)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (svc AuthenticationService) ForgotPasswordByPhonenumber(userRequest auth_models.ForgotPasswordPhoneNumberRequest) error {

	isOtpPass, err := svc.smsRepo.VerifyOTPViaLink(userRequest.OTPToken, userRequest.OTPRefCode, userRequest.OTPPin)

	if err != nil {
		return err
	}

	if !isOtpPass {
		return errors.New("otp invalid")
	}

	userRequest.PhoneNumber = utils.NormalizePhonenumber(userRequest.PhoneNumber)

	if exists, err := svc.CheckExistsPhonenumber(userRequest.PhoneNumber); err != nil {
		return err
	} else if exists {
		return errors.New("phonenumber is exists")
	}

	userFind, err := svc.authRepo.FindByPhonenumber(context.Background(), userRequest.PhoneNumberField)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return err
	}

	if len(userFind.PhoneNumber) < 1 {
		return errors.New("phone number is not exists")
	}

	hashPassword, err := svc.passwordEncoder(userRequest.Password)

	if err != nil {
		return err
	}

	userFind.Password = hashPassword

	err = svc.authRepo.UpdateUser(context.Background(), userFind.Username, *userFind)

	if err != nil {
		return err
	}

	return nil
}

func (svc AuthenticationService) Update(username string, userRequest auth_models.UserProfileRequest) error {

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

	userFind.UserDetail = userRequest.UserDetail
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

func (svc AuthenticationService) Profile(username string) (auth_models.UserProfile, error) {

	userProfile := auth_models.UserProfile{}
	user, err := svc.authRepo.FindUser(context.Background(), username)
	if err != nil {
		return userProfile, err
	}
	userProfile.Username = user.Username
	userProfile.UserDetail = user.UserDetail

	return userProfile, nil
}

func (svc AuthenticationService) AccessShop(shopID string, username string, authorizationHeader string, authContext models.AuthenticationContext) error {

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
		logger.GetLogger().Error(err.Error())
	}

	err = svc.shopUserAccessLogRepo.Create(
		context.Background(),
		auth_models.ShopUserAccessLog{
			ShopID:         shopID,
			Username:       username,
			Ip:             authContext.Ip,
			LastAccessedAt: lastAccessedAt,
		})

	if err != nil {
		logger.GetLogger().Error(err.Error())
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
		user := auth_models.UserDoc{}

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
