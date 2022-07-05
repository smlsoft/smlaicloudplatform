package journalreport_test

import (
	chartofaccountModel "smlcloudplatform/pkg/vfgl/chartofaccount/models"
	"smlcloudplatform/pkg/vfgl/journalreport"
	"smlcloudplatform/pkg/vfgl/journalreport/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJournalreportRepository struct {
	mock.Mock
}

func (m *MockJournalreportRepository) GetDataTrialBalance(shopId string, accountGroup string, startDate time.Time, endDate time.Time) ([]models.TrialBalanceSheetAccountDetail, error) {
	ret := m.Called(shopId, accountGroup, startDate, endDate)
	return ret.Get(0).([]models.TrialBalanceSheetAccountDetail), ret.Error(1)
}

func (m *MockJournalreportRepository) GetDataProfitAndLoss(shopId string, accountGroup string, startDate time.Time, endDate time.Time) ([]models.ProfitAndLossSheetAccountDetail, error) {
	ret := m.Called(shopId, accountGroup, startDate, endDate)
	return ret.Get(0).([]models.ProfitAndLossSheetAccountDetail), ret.Error(1)
}

func (m *MockJournalreportRepository) GetDataBalanceSheet(shopId string, accountGroup string, endDate time.Time) ([]models.BalanceSheetAccountDetail, error) {
	ret := m.Called(shopId, accountGroup, endDate)
	return ret.Get(0).([]models.BalanceSheetAccountDetail), ret.Error(1)
}

func TestProcessBalanceSheetReport(t *testing.T) {

	var balances = journalreport.MockBalanceSheetDetailReport()
	endDate := time.Date(2022, 05, 31, 0, 0, 0, 0, time.UTC)

	repo := new(MockJournalreportRepository)
	repo.On("GetDataBalanceSheet", "TESTSHOP", "01", endDate).Return(balances, nil)

	var liabilities []models.BalanceSheetAccountDetail

	fixReportDate := time.Now()
	want := &models.BalanceSheetReport{
		ReportDate:   fixReportDate,
		EndDate:      endDate,
		AccountGroup: "01",
		Assets: &[]models.BalanceSheetAccountDetail{
			{
				ChartOfAccountPG: chartofaccountModel.ChartOfAccountPG{
					AccountCode:     "12101",
					AccountName:     "เงินฝากธนาคาร บัญชี 1 (เงินล้าน)",
					AccountCategory: 1,
				},
				Amount: 10000,
			},
			{
				ChartOfAccountPG: chartofaccountModel.ChartOfAccountPG{
					AccountCode:     "13010",
					AccountName:     "ลูกหนี้เงินกู้ - บัญชี 1 (เงินล้าน)",
					AccountCategory: 1,
				},
				Amount: 20000,
			},
			{
				ChartOfAccountPG: chartofaccountModel.ChartOfAccountPG{
					AccountCode:     "11010",
					AccountName:     "เงินสด - บัญชี 1",
					AccountCategory: 1,
				},
				Amount: 20,
			},
		},
		Liabilities: &liabilities,
		OwnesEquities: &[]models.BalanceSheetAccountDetail{
			{
				ChartOfAccountPG: chartofaccountModel.ChartOfAccountPG{
					AccountCode:     "32010",
					AccountName:     "ทุน - เงินล้าน",
					AccountCategory: 3,
				},
				Amount: 30000,
			},
			{
				ChartOfAccountPG: chartofaccountModel.ChartOfAccountPG{
					AccountName:     "กำไร (ขาดทุน) สุทธิ",
					AccountCategory: 3,
				},
				Amount: 20,
			},
		},
		TotalAssetAmount:                    30020,
		TotalOwnersEquityAmount:             30020,
		TotalLiabilityAndOwnersEquityAmount: 30020,
	}

	service := journalreport.NewJournalReportService(repo)
	get, err := service.ProcessBalanceSheetReport("TESTSHOP", "01", endDate)
	get.ReportDate = fixReportDate
	assert.Nil(t, err, "Error should be nil")

	assert.Equal(t, get, want, "Process BalanceSheet Report Not Match")
}
