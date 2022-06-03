package services_test

import (
	"smlcloudplatform/pkg/models/vfgl"
	"smlcloudplatform/pkg/vfgl/chartofaccount/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	ret := m.Called(accountCode)
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
	mockRepo.On("Create", get).Return(nil)

	chartOfAccountService := services.NewChartOfAccountConsumeService(mockRepo)

	err := chartOfAccountService.Create(give)
	assert.Nil(t, err, "Error should be nil")
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
