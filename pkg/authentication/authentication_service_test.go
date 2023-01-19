package authentication_test

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/authentication"
	"smlcloudplatform/pkg/shop/models"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mockLoginData(authRepo *AuthenticationRepositoryMock, shopUserRepo *ShopUserRepositoryMock, microAuthServiceMock *AuthServiceMock) {

	shopID := "SHOP_ID_TEST"
	role := uint8(2)

	tokenMock := "TOKEN_MOCK"
	//authRepo FindUser
	userDoc1 := models.UserDoc{}
	userDoc1.Username = "tester1"
	userDoc1.Password = "tester1"
	userDoc1.Name = "tester1"

	authRepo.On("FindUser", userDoc1.Username).Return(&userDoc1, nil)
	authRepo.On("FindUser", "").Return(&models.UserDoc{}, nil)
	authRepo.On("FindUser", "tester2").Return(&models.UserDoc{}, nil)

	authRepo.On("FindUser", "user_register").Return(&models.UserDoc{}, nil)

	//authRepo CreateUser
	userDoc2 := models.UserDoc{}
	userDoc2.Username = "user_register"
	userDoc2.Password = "register_password_success"
	userDoc2.Name = "user_register"
	userDoc2.CreatedAt = MockTime()

	authRepo.On("CreateUser", userDoc2).Return(MockObjectID(), nil)

	//microAuth
	microAuthServiceMock.On("GenerateTokenWithRedis", micromodels.UserInfo{
		Username: userDoc1.Username,
		Name:     userDoc1.Name,
	}).Return(tokenMock, nil)

	microAuthServiceMock.On("SelectShop", tokenMock, shopID, role).Return(nil)

	//shopUser
	shopUserRepo.On("FindByShopIDAndUsername", shopID, userDoc1.Username).Return(models.ShopUser{
		ID:       primitive.NewObjectID(),
		Username: userDoc1.Username,
		ShopID:   shopID,
		Role:     role,
	}, nil)

	shopUserRepo.On("FindByShopIDAndUsername", "SHOP_ID_INVALID", userDoc1.Username).Return(models.ShopUser{}, nil)
}

func TestAuthService_Login(t *testing.T) {
	shopID := "SHOP_ID_TEST"

	authRepo := new(AuthenticationRepositoryMock)
	shopUserRepo := new(ShopUserRepositoryMock)
	shopUserAccessLogRepo := new(ShopUserAccessLogRepositoryMock)
	microAuthServiceMock := &AuthServiceMock{}

	mockLoginData(authRepo, shopUserRepo, microAuthServiceMock)

	type args struct {
		username string
		password string
		shopID   string
	}

	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData string
	}{
		{
			name: "login success",
			args: args{
				shopID:   shopID,
				username: "tester1",
				password: "tester1",
			},
			wantErr:  false,
			wantData: "TOKEN_MOCK",
		},
		{
			name: "login success without shop id",
			args: args{
				username: "tester1",
				password: "tester1",
			},
			wantErr:  false,
			wantData: "TOKEN_MOCK",
		},
		{
			name: "login failure invalid shop id",
			args: args{
				shopID:   "SHOP_ID_INVALID",
				username: "tester1",
				password: "tester1",
			},
			wantErr:  true,
			wantData: "TOKEN_MOCK",
		},
		{
			name: "login failure password invalid",
			args: args{
				shopID:   shopID,
				username: "tester1",
				password: "invalidpassword",
			},
			wantErr:  true,
			wantData: "TOKEN_MOCK",
		},
		{
			name: "login failure username empty",
			args: args{
				shopID:   shopID,
				username: "",
				password: "invalidpassword",
			},
			wantErr:  true,
			wantData: "TOKEN_MOCK",
		},
		{
			name: "login failure username and password empty",
			args: args{
				shopID:   shopID,
				username: "",
				password: "",
			},
			wantErr:  true,
			wantData: "TOKEN_MOCK",
		},
	}

	authService := authentication.NewAuthenticationService(authRepo, shopUserRepo, shopUserAccessLogRepo, microAuthServiceMock, MockHashPassword, MockCheckPasswordHash, MockTime)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			userReq := &models.UserLoginRequest{}
			userReq.Username = tt.args.username
			userReq.Password = tt.args.password
			userReq.ShopID = tt.args.shopID

			authContext := authentication.AuthenticationContext{
				Ip: "localhost",
			}

			tokenResult, err := authService.Login(userReq, authContext)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Empty(t, tokenResult)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, tokenResult)
				assert.EqualValues(t, tt.wantData, tokenResult)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	authRepo := new(AuthenticationRepositoryMock)
	shopUserRepo := new(ShopUserRepositoryMock)
	shopUserAccessLogRepo := new(ShopUserAccessLogRepositoryMock)
	microAuthServiceMock := &AuthServiceMock{}

	mockLoginData(authRepo, shopUserRepo, microAuthServiceMock)

	type args struct {
		username string
		password string
		name     string
	}

	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData string
	}{
		{
			name: "register success",
			args: args{
				username: "user_register",
				password: "register_password_success",
				name:     "user_register",
			},
			wantErr:  false,
			wantData: "62f9cb12c76fd9e83ac1b2ff",
		},
		{
			name: "register failure user is exist",
			args: args{
				username: "tester1",
				password: "register_password_failure",
				name:     "user_register",
			},
			wantErr:  true,
			wantData: "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			authService := authentication.NewAuthenticationService(authRepo, shopUserRepo, shopUserAccessLogRepo, microAuthServiceMock, MockHashPassword, MockCheckPasswordHash, MockTime)

			userReq := models.UserRequest{}
			userReq.Username = tt.args.username
			userReq.Password = tt.args.password
			userReq.Name = tt.args.name

			idx, err := authService.Register(userReq)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Empty(t, "")
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, idx)
				assert.EqualValues(t, tt.wantData, idx)
			}
		})
	}
}

func TestAuthService_Update(t *testing.T) {
	authRepo := new(AuthenticationRepositoryMock)
	shopUserRepo := new(ShopUserRepositoryMock)
	shopUserAccessLogRepo := new(ShopUserAccessLogRepositoryMock)
	microAuthServiceMock := &AuthServiceMock{}

	userDoc := &models.UserDoc{}

	userDoc.Username = "user_update"
	userDoc.Name = "user_update"
	userDoc.UpdatedAt = MockTime()

	authRepo.On("FindUser", "user_update").Return(userDoc, nil)

	userDocUpdate := models.UserDoc{}
	userDocUpdate.Username = "user_update"
	userDocUpdate.Name = "new name"
	userDocUpdate.UpdatedAt = MockTime()

	authRepo.On("UpdateUser", "user_update", userDocUpdate).Return(nil)

	type args struct {
		username string
		name     string
	}

	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update success",
			args: args{
				username: "user_update",
				name:     "new name",
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			authService := authentication.NewAuthenticationService(authRepo, shopUserRepo, shopUserAccessLogRepo, microAuthServiceMock, MockHashPassword, MockCheckPasswordHash, MockTime)

			userReq := models.UserProfileRequest{}
			userReq.Name = tt.args.name

			err := authService.Update(tt.args.username, userReq)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Empty(t, "")
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthService_UpdatePassword(t *testing.T) {
	authRepo := new(AuthenticationRepositoryMock)
	shopUserRepo := new(ShopUserRepositoryMock)
	shopUserAccessLogRepo := new(ShopUserAccessLogRepositoryMock)
	microAuthServiceMock := &AuthServiceMock{}

	userDoc := &models.UserDoc{}

	userDoc.Username = "user_update"
	userDoc.Password = "current_password"
	userDoc.UpdatedAt = MockTime()

	authRepo.On("FindUser", "user_update").Return(userDoc, nil)

	userDocUpdate := models.UserDoc{}
	userDocUpdate.Username = "user_update"
	userDocUpdate.Password = "new_password"
	userDocUpdate.UpdatedAt = MockTime()

	authRepo.On("UpdateUser", "user_update", userDocUpdate).Return(nil)

	type args struct {
		username        string
		currentPassword string
		newPassword     string
	}

	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update password success",
			args: args{
				username:        "user_update",
				currentPassword: "current_password",
				newPassword:     "new_password",
			},
			wantErr: false,
		},
		{
			name: "update password failure",
			args: args{
				username:        "user_update",
				currentPassword: "current_password_invalid",
				newPassword:     "new_password",
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			authService := authentication.NewAuthenticationService(authRepo, shopUserRepo, shopUserAccessLogRepo, microAuthServiceMock, MockHashPassword, MockCheckPasswordHash, MockTime)

			err := authService.UpdatePassword(tt.args.username, tt.args.currentPassword, tt.args.newPassword)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Empty(t, "")
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthService_AccessShop(t *testing.T) {
	authRepo := new(AuthenticationRepositoryMock)
	shopUserRepo := new(ShopUserRepositoryMock)
	shopUserAccessLogRepo := new(ShopUserAccessLogRepositoryMock)
	microAuthServiceMock := &AuthServiceMock{}

	microAuthServiceMock.On("GetTokenFromAuthorizationHeader", "authorization_header_valid").Return("valid_token", nil)
	microAuthServiceMock.On("GetTokenFromAuthorizationHeader", "").Return("", errors.New("authorization is not empty"))

	shopUserRepo.On("FindByShopIDAndUsername", "shop_test", "user_access_shop").Return(models.ShopUser{
		ID:       MockObjectID(),
		Username: "user_access_shop",
		ShopID:   "shop_test",
		Role:     uint8(0),
	}, nil)

	shopUserRepo.On("FindByShopIDAndUsername", "shop_test_invalid", "user_access_shop").Return(models.ShopUser{}, nil)

	microAuthServiceMock.On("SelectShop", "valid_token", "shop_test", uint8(0)).Return(nil)
	microAuthServiceMock.On("SelectShop", "valid_token_invalid", "shop_test_invalid", uint8(0)).Return(errors.New("select shop failed"))

	type args struct {
		shopID              string
		username            string
		authorizationHeader string
	}

	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success access shop ",
			args: args{
				shopID:              "shop_test",
				username:            "user_access_shop",
				authorizationHeader: "authorization_header_valid",
			},
			wantErr: false,
		},
		{
			name: "failure authorization empty",
			args: args{
				shopID:              "shop_test",
				username:            "user_access_shop",
				authorizationHeader: "",
			},
			wantErr: true,
		},
		{
			name: "failure shop invalid",
			args: args{
				shopID:              "shop_test_invalid",
				username:            "user_access_shop",
				authorizationHeader: "authorization_header_valid",
			},
			wantErr: true,
		},
		{
			name: "failure access shop failed",
			args: args{
				shopID:              "shop_test_invalid",
				username:            "user_access_shop",
				authorizationHeader: "authorization_header_valid",
			},
			wantErr: true,
		},
	}

	authContext := authentication.AuthenticationContext{
		Ip: "localhost",
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			authService := authentication.NewAuthenticationService(authRepo, shopUserRepo, shopUserAccessLogRepo, microAuthServiceMock, MockHashPassword, MockCheckPasswordHash, MockTime)
			err := authService.AccessShop(tt.args.shopID, tt.args.username, tt.args.authorizationHeader, authContext)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Empty(t, "")
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

type AuthenticationRepositoryMock struct {
	mock.Mock
}

func (m *AuthenticationRepositoryMock) FindUser(id string) (*models.UserDoc, error) {
	args := m.Called(id)
	return args.Get(0).(*models.UserDoc), args.Error(1)
}

func (m *AuthenticationRepositoryMock) CreateUser(doc models.UserDoc) (primitive.ObjectID, error) {
	args := m.Called(doc)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *AuthenticationRepositoryMock) UpdateUser(username string, userDoc models.UserDoc) error {
	args := m.Called(username, userDoc)
	return args.Error(0)
}

type ShopUserRepositoryMock struct {
	mock.Mock
}

func (m *ShopUserRepositoryMock) Save(shopID string, username string, role models.UserRole) error {
	args := m.Called(shopID, username, role)
	return args.Error(0)
}

func (m *ShopUserRepositoryMock) UpdateLastAccess(shopID string, username string, lastAccessedAt time.Time) error {
	args := m.Called(shopID, username, lastAccessedAt)
	return args.Error(0)
}

func (m *ShopUserRepositoryMock) SaveFavorite(shopID string, username string, isFavorite bool) error {
	args := m.Called(shopID, username, isFavorite)
	return args.Error(0)
}

func (m *ShopUserRepositoryMock) Delete(shopID string, username string) error {
	args := m.Called(shopID, username)
	return args.Error(0)
}

func (m *ShopUserRepositoryMock) FindByShopIDAndUsernameInfo(shopID string, username string) (models.ShopUserInfo, error) {
	args := m.Called(shopID, username)
	return args.Get(0).(models.ShopUserInfo), args.Error(1)
}

func (m *ShopUserRepositoryMock) FindByShopIDAndUsername(shopID string, username string) (models.ShopUser, error) {
	args := m.Called(shopID, username)
	return args.Get(0).(models.ShopUser), args.Error(1)
}

func (m *ShopUserRepositoryMock) FindRole(shopID string, username string) (models.UserRole, error) {
	args := m.Called(shopID, username)
	return args.Get(0).(models.UserRole), args.Error(1)
}
func (m *ShopUserRepositoryMock) FindByShopID(shopID string) (*[]models.ShopUser, error) {
	args := m.Called(shopID)
	return args.Get(0).(*[]models.ShopUser), args.Error(1)
}

func (m *ShopUserRepositoryMock) FindByUsername(username string) (*[]models.ShopUser, error) {
	args := m.Called(username)
	return args.Get(0).(*[]models.ShopUser), args.Error(1)
}
func (m *ShopUserRepositoryMock) FindByUsernamePage(username string, pageable micromodels.Pageable) ([]models.ShopUserInfo, mongopagination.PaginationData, error) {
	args := m.Called(username, pageable)
	return args.Get(0).([]models.ShopUserInfo), args.Get(1).(mongopagination.PaginationData), args.Error(2)
}
func (m *ShopUserRepositoryMock) FindByUserInShopPage(shopID string, pageable micromodels.Pageable) ([]models.ShopUser, mongopagination.PaginationData, error) {
	args := m.Called(shopID, pageable)
	return args.Get(0).([]models.ShopUser), args.Get(1).(mongopagination.PaginationData), args.Error(2)
}

// Shop User Access Log
type ShopUserAccessLogRepositoryMock struct {
	mock.Mock
}

func (m *ShopUserAccessLogRepositoryMock) Create(shopUserAccessLog models.ShopUserAccessLog) error {
	args := m.Called(shopUserAccessLog)
	return args.Error(0)
}

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) MWFuncWithRedisMixShop(cacher microservice.ICacher, shopPath []string, publicPath ...string) echo.MiddlewareFunc {
	args := m.Called(cacher, shopPath, publicPath)
	return args.Get(0).(echo.MiddlewareFunc)
}

func (m *AuthServiceMock) MWFuncWithRedis(cacher microservice.ICacher, publicPath ...string) echo.MiddlewareFunc {
	args := m.Called(cacher, publicPath)
	return args.Get(0).(echo.MiddlewareFunc)
}

func (m *AuthServiceMock) MWFuncWithShop(cacher microservice.ICacher, publicPath ...string) echo.MiddlewareFunc {
	args := m.Called(cacher, publicPath)
	return args.Get(0).(echo.MiddlewareFunc)
}

func (m *AuthServiceMock) GetPrefixCacheKey(tokenType microservice.TokenType) string {
	args := m.Called(tokenType)
	return args.String(0)
}

func (m *AuthServiceMock) GetTokenFromContext(c echo.Context) (*microservice.TokenContext, error) {

	args := m.Called(c)

	return args.Get(0).(*microservice.TokenContext), args.Error(1)
}

func (m *AuthServiceMock) GetTokenFromAuthorizationHeader(tokenType microservice.TokenType, tokenAuthorization string) (string, error) {

	args := m.Called(tokenType, tokenAuthorization)

	return args.String(0), args.Error(1)
}

func (m *AuthServiceMock) GenerateTokenWithRedis(tokenType microservice.TokenType, userInfo micromodels.UserInfo) (string, error) {

	args := m.Called(tokenType, userInfo)
	return args.String(0), args.Error(1)
}

func (m *AuthServiceMock) GenerateTokenWithRedisExpire(tokenType microservice.TokenType, userInfo micromodels.UserInfo, expireTime time.Duration) (string, error) {

	args := m.Called(tokenType, userInfo, expireTime)
	return args.String(0), args.Error(1)
}

func (m *AuthServiceMock) SelectShop(tokenType microservice.TokenType, tokenStr string, shopID string, role uint8) error {

	args := m.Called(tokenType, tokenStr, shopID, role)
	return args.Error(0)
}

func (m *AuthServiceMock) ExpireToken(tokenType microservice.TokenType, tokenAuthorizationHeader string) error {
	args := m.Called(tokenType, tokenAuthorizationHeader)
	return args.Error(0)
}

func MockObjectID() primitive.ObjectID {
	idx, _ := primitive.ObjectIDFromHex("62f9cb12c76fd9e83ac1b2ff")
	return idx
}

func MockHashPassword(password string) (string, error) {
	return password, nil
}

func MockCheckPasswordHash(password string, hash string) bool {
	return password == hash
}

func MockTime() time.Time {
	timeVal, _ := time.Parse("2006-01-02 15:04:05", "2022-08-30 00:00:00")
	return timeVal
}
