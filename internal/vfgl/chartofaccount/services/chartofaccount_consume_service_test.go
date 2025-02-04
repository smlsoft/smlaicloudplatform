package services_test

import (
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/vfgl/chartofaccount/models"
	"smlaicloudplatform/internal/vfgl/chartofaccount/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockChartOfAccountRepository struct {
	mock.Mock
}

func (m *MockChartOfAccountRepository) CreateInBatch(docList []models.ChartOfAccountPG) error {
	ret := m.Called(docList)
	return ret.Error(1)
}

func (m *MockChartOfAccountRepository) Create(doc models.ChartOfAccountPG) error {
	ret := m.Called(doc)
	return ret.Error(0)
}

func (m *MockChartOfAccountRepository) Update(shopID string, accountCode string, doc models.ChartOfAccountPG) error {
	ret := m.Called(shopID, accountCode, doc)
	return ret.Error(0)
}

func (m *MockChartOfAccountRepository) Delete(shopID string, accountCode string) error {
	ret := m.Called(shopID, accountCode)
	return ret.Error(1)
}

func (m *MockChartOfAccountRepository) Get(shopID string, accountCode string) (*models.ChartOfAccountPG, error) {
	var charts []models.ChartOfAccountPG

	charts = append(charts, models.ChartOfAccountPG{})
	ret := m.Called(shopID, accountCode)
	return ret.Get(0).(*models.ChartOfAccountPG), ret.Error(1)
}

func ListChartOfAccountPG() []models.ChartOfAccountPG {

	var charts []models.ChartOfAccountPG

	charts = append(charts, models.ChartOfAccountPG{})
	charts = append(charts, models.ChartOfAccountPG{})
	charts = append(charts, models.ChartOfAccountPG{})
	return charts
}

func TestChartOfAccountServiceCreate(t *testing.T) {

	want := models.ChartOfAccountPG{
		AccountCode:        "99999",
		AccountName:        "999999",
		AccountCategory:    5,
		AccountBalanceType: 2,
		AccountGroup:       "99999",
		AccountLevel:       9999,
		ShopIdentity: common.ShopIdentity{
			ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
		},
	}

	give := models.ChartOfAccountDoc{
		ChartOfAccountData: models.ChartOfAccountData{
			ChartOfAccountInfo: models.ChartOfAccountInfo{
				ChartOfAccount: models.ChartOfAccount{
					AccountCode:        "99999",
					AccountName:        "999999",
					AccountCategory:    5,
					AccountBalanceType: 2,
					AccountGroup:       "99999",
					AccountLevel:       9999,
				},
			},
			ShopIdentity: common.ShopIdentity{
				ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
			},
		},
	}

	mockRepo := new(MockChartOfAccountRepository)
	mockRepo.On("Create", want).Return(nil)

	chartOfAccountService := services.NewChartOfAccountConsumeService(mockRepo)

	get, err := chartOfAccountService.Create(give)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, *get, want, "Failed After Created Got Data Not Match")
}

func TestChartOfAccountServiceUpdate(t *testing.T) {

	get := models.ChartOfAccountPG{
		AccountCode: "0001",
	}

	give := models.ChartOfAccountDoc{
		ChartOfAccountData: models.ChartOfAccountData{
			ChartOfAccountInfo: models.ChartOfAccountInfo{
				ChartOfAccount: models.ChartOfAccount{
					AccountCode: "0001",
				},
			},
		},
	}

	mockRepo := new(MockChartOfAccountRepository)
	mockRepo.On("Update", "SHOPTEST", "0001", get).Return(nil)

	chartOfAccountService := services.NewChartOfAccountConsumeService(mockRepo)

	err := chartOfAccountService.Update("SHOPTEST", "0001", give)
	assert.Nil(t, err, "Error should be nil")
}

func TestChartOfAccountConsumeServiceUpsertWhenNotFoundDataInsertNew(t *testing.T) {
	give := models.ChartOfAccountDoc{
		ChartOfAccountData: models.ChartOfAccountData{
			ChartOfAccountInfo: models.ChartOfAccountInfo{
				ChartOfAccount: models.ChartOfAccount{
					AccountCode:        "99999",
					AccountName:        "999999",
					AccountCategory:    5,
					AccountBalanceType: 2,
					AccountGroup:       "99999",
					AccountLevel:       9999,
				},
			},
			ShopIdentity: common.ShopIdentity{
				ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
			},
		},
	}

	want := models.ChartOfAccountPG{
		AccountCode:        "99999",
		AccountName:        "999999",
		AccountCategory:    5,
		AccountBalanceType: 2,
		AccountGroup:       "99999",
		AccountLevel:       9999,
		ShopIdentity: common.ShopIdentity{
			ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
		},
	}

	mockRepo := new(MockChartOfAccountRepository)
	mockRepo.On("Get", give.ShopID, give.AccountCode).Return(&models.ChartOfAccountPG{}, gorm.ErrRecordNotFound)
	mockRepo.On("Create", want).Return(nil)
	// mockRepo.On("Update", give.ShopID, give.AccountCode, give).Return(nil)

	svc := services.NewChartOfAccountConsumeService(mockRepo)
	get, err := svc.Upsert(give.ShopID, give)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, get, &want, "After Update Are Not Equal")

}

func TestChartOfAccountConsumeServiceUpsertWhenFoundDataUpdateOld(t *testing.T) {
	give := models.ChartOfAccountDoc{
		ChartOfAccountData: models.ChartOfAccountData{
			ChartOfAccountInfo: models.ChartOfAccountInfo{
				ChartOfAccount: models.ChartOfAccount{
					AccountCode:        "99999",
					AccountName:        "999999",
					AccountCategory:    5,
					AccountBalanceType: 2,
					AccountGroup:       "99999",
					AccountLevel:       9999,
				},
			},
			ShopIdentity: common.ShopIdentity{
				ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
			},
		},
	}

	giveGetFromRepository := models.ChartOfAccountPG{
		AccountCode:        "99999",
		AccountName:        "",
		AccountCategory:    0,
		AccountBalanceType: 0,
		AccountGroup:       "",
		AccountLevel:       0,
		ShopIdentity: common.ShopIdentity{
			ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
		},
	}

	want := models.ChartOfAccountPG{
		AccountCode:        "99999",
		AccountName:        "999999",
		AccountCategory:    5,
		AccountBalanceType: 2,
		AccountGroup:       "99999",
		AccountLevel:       9999,
		ShopIdentity: common.ShopIdentity{
			ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
		},
	}

	mockRepo := new(MockChartOfAccountRepository)
	mockRepo.On("Get", give.ShopID, give.AccountCode).Return(&giveGetFromRepository, nil)
	mockRepo.On("Update", give.ShopID, give.AccountCode, want).Return(nil)

	svc := services.NewChartOfAccountConsumeService(mockRepo)
	get, err := svc.Upsert(give.ShopID, give)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, get, &want, "After Update Are Not Equal")

}
