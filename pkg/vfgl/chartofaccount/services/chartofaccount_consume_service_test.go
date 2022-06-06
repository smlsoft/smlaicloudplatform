package services_test

import (
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/vfgl/chartofaccount/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockChartOfAccountRepository struct {
	mock.Mock
}

func (m *MockChartOfAccountRepository) CreateInBatch(docList []vfgl.ChartOfAccountPG) error {
	ret := m.Called(docList)
	return ret.Error(1)
}

func (m *MockChartOfAccountRepository) Create(doc vfgl.ChartOfAccountPG) error {
	ret := m.Called(doc)
	return ret.Error(0)
}

func (m *MockChartOfAccountRepository) Update(shopID string, accountCode string, doc vfgl.ChartOfAccountPG) error {
	ret := m.Called(shopID, accountCode, doc)
	return ret.Error(0)
}

func (m *MockChartOfAccountRepository) Delete(shopID string, accountCode string) error {
	ret := m.Called(shopID, accountCode)
	return ret.Error(1)
}

func (m *MockChartOfAccountRepository) Get(shopID string, accountCode string) (*vfgl.ChartOfAccountPG, error) {
	var charts []vfgl.ChartOfAccountPG

	charts = append(charts, vfgl.ChartOfAccountPG{})
	ret := m.Called(shopID, accountCode)
	return ret.Get(0).(*vfgl.ChartOfAccountPG), ret.Error(1)
}

func ListChartOfAccountPG() []vfgl.ChartOfAccountPG {

	var charts []vfgl.ChartOfAccountPG

	charts = append(charts, vfgl.ChartOfAccountPG{})
	charts = append(charts, vfgl.ChartOfAccountPG{})
	charts = append(charts, vfgl.ChartOfAccountPG{})
	return charts
}

func TestChartOfAccountServiceCreate(t *testing.T) {

	want := vfgl.ChartOfAccountPG{
		AccountCode:        "99999",
		AccountName:        "999999",
		AccountCategory:    5,
		AccountBalanceType: 2,
		AccountGroup:       "99999",
		AccountLevel:       9999,
		ShopIdentity: models.ShopIdentity{
			ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
		},
	}

	give := vfgl.ChartOfAccountDoc{
		ChartOfAccountData: vfgl.ChartOfAccountData{
			ChartOfAccountInfo: vfgl.ChartOfAccountInfo{
				ChartOfAccount: vfgl.ChartOfAccount{
					AccountCode:        "99999",
					AccountName:        "999999",
					AccountCategory:    5,
					AccountBalanceType: 2,
					AccountGroup:       "99999",
					AccountLevel:       9999,
				},
			},
			ShopIdentity: models.ShopIdentity{
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

	get := vfgl.ChartOfAccountPG{
		AccountCode: "0001",
	}

	give := vfgl.ChartOfAccountDoc{
		ChartOfAccountData: vfgl.ChartOfAccountData{
			ChartOfAccountInfo: vfgl.ChartOfAccountInfo{
				ChartOfAccount: vfgl.ChartOfAccount{
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
	give := vfgl.ChartOfAccountDoc{
		ChartOfAccountData: vfgl.ChartOfAccountData{
			ChartOfAccountInfo: vfgl.ChartOfAccountInfo{
				ChartOfAccount: vfgl.ChartOfAccount{
					AccountCode:        "99999",
					AccountName:        "999999",
					AccountCategory:    5,
					AccountBalanceType: 2,
					AccountGroup:       "99999",
					AccountLevel:       9999,
				},
			},
			ShopIdentity: models.ShopIdentity{
				ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
			},
		},
	}

	want := vfgl.ChartOfAccountPG{
		AccountCode:        "99999",
		AccountName:        "999999",
		AccountCategory:    5,
		AccountBalanceType: 2,
		AccountGroup:       "99999",
		AccountLevel:       9999,
		ShopIdentity: models.ShopIdentity{
			ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
		},
	}

	mockRepo := new(MockChartOfAccountRepository)
	mockRepo.On("Get", give.ShopID, give.AccountCode).Return(&vfgl.ChartOfAccountPG{}, gorm.ErrRecordNotFound)
	mockRepo.On("Create", want).Return(nil)
	// mockRepo.On("Update", give.ShopID, give.AccountCode, give).Return(nil)

	svc := services.NewChartOfAccountConsumeService(mockRepo)
	get, err := svc.Upsert(give.ShopID, give)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, get, &want, "After Update Are Not Equal")

}

func TestChartOfAccountConsumeServiceUpsertWhenFoundDataUpdateOld(t *testing.T) {
	give := vfgl.ChartOfAccountDoc{
		ChartOfAccountData: vfgl.ChartOfAccountData{
			ChartOfAccountInfo: vfgl.ChartOfAccountInfo{
				ChartOfAccount: vfgl.ChartOfAccount{
					AccountCode:        "99999",
					AccountName:        "999999",
					AccountCategory:    5,
					AccountBalanceType: 2,
					AccountGroup:       "99999",
					AccountLevel:       9999,
				},
			},
			ShopIdentity: models.ShopIdentity{
				ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
			},
		},
	}

	giveGetFromRepository := vfgl.ChartOfAccountPG{
		AccountCode:        "99999",
		AccountName:        "",
		AccountCategory:    0,
		AccountBalanceType: 0,
		AccountGroup:       "",
		AccountLevel:       0,
		ShopIdentity: models.ShopIdentity{
			ShopID: "27dcEdktOoaSBYFmnN6G6ett4Jb",
		},
	}

	want := vfgl.ChartOfAccountPG{
		AccountCode:        "99999",
		AccountName:        "999999",
		AccountCategory:    5,
		AccountBalanceType: 2,
		AccountGroup:       "99999",
		AccountLevel:       9999,
		ShopIdentity: models.ShopIdentity{
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
