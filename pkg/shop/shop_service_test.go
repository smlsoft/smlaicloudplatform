package shop_test

import (
	utilmock "smlcloudplatform/mock"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/shop"
	"smlcloudplatform/pkg/shop/models"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/userplant/mongopagination"
)

func TestShop_Create(t *testing.T) {
	shopRepo := new(ShopRepositoryMock)
	shopUserRepo := new(ShopUserRepositoryMock)

	shopRepo.On("Create", models.ShopDoc{

		ShopInfo: models.ShopInfo{
			DocIdentity: common.DocIdentity{
				GuidFixed: utilmock.MockGUID(),
			},
			Shop: models.Shop{
				Name1:     "shop_name",
				Telephone: "0000000000",
			},
		},
		ActivityDoc: common.ActivityDoc{
			CreatedBy: "user_create",
			CreatedAt: utilmock.MockTime(),
		},
	}).Return("", nil)

	shopUserRepo.On("Save", utilmock.MockGUID(), "user_create", models.ROLE_OWNER).Return(nil)

	type args struct {
		username string
		shop     models.Shop
	}

	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData string
	}{
		{
			name: "success create shop",
			args: args{
				username: "user_create",
				shop: models.Shop{
					Name1:     "shop_name",
					Telephone: "0000000000",
				},
			},
			wantErr:  false,
			wantData: utilmock.MockGUID(),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			shopSvc := shop.NewShopService(shopRepo, shopUserRepo, utilmock.MockGUID, utilmock.MockTime)

			shopGUID, err := shopSvc.CreateShop(tt.args.username, tt.args.shop)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Empty(t, shopGUID)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, shopGUID)
				assert.EqualValues(t, tt.wantData, shopGUID)
			}
		})
	}

}

type ShopRepositoryMock struct {
	mock.Mock
}

func (m *ShopRepositoryMock) Create(shop models.ShopDoc) (string, error) {
	args := m.Called(shop)

	return args.String(0), args.Error(1)
}

func (m *ShopRepositoryMock) Update(guid string, shop models.ShopDoc) error {
	args := m.Called(guid, shop)

	return args.Error(0)
}
func (m *ShopRepositoryMock) FindByGuid(guid string) (models.ShopDoc, error) {
	args := m.Called(guid)
	return args.Get(0).(models.ShopDoc), args.Error(0)
}
func (m *ShopRepositoryMock) FindPage(q string, page int, limit int) ([]models.ShopInfo, mongopagination.PaginationData, error) {
	args := m.Called(q, page, limit)

	return args.Get(0).([]models.ShopInfo), args.Get(1).(mongopagination.PaginationData), args.Error(2)
}
func (m *ShopRepositoryMock) Delete(guid string, username string) error {
	args := m.Called(guid, username)

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

func (m *ShopUserRepositoryMock) FindByUsernamePage(username string, q string, page int, limit int) ([]models.ShopUserInfo, mongopagination.PaginationData, error) {
	args := m.Called(username, q, page, limit)
	return args.Get(0).([]models.ShopUserInfo), args.Get(1).(mongopagination.PaginationData), args.Error(2)
}

func (m *ShopUserRepositoryMock) FindByUserInShopPage(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ShopUser, mongopagination.PaginationData, error) {
	args := m.Called(shopID, q, page, limit, sort)
	return args.Get(0).([]models.ShopUser), args.Get(1).(mongopagination.PaginationData), args.Error(2)
}
