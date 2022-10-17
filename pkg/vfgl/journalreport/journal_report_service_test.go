package journalreport_test

import (
	chartofaccountModel "smlcloudplatform/pkg/vfgl/chartofaccount/models"
	"smlcloudplatform/pkg/vfgl/journalreport"
	"smlcloudplatform/pkg/vfgl/journalreport/models"
	"testing"
	"time"

	mocktest "smlcloudplatform/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJournalreportRepository struct {
	mock.Mock
}

func (m *MockJournalreportRepository) GetDataTrialBalance(shopId string, accountGroup string, includeCloseAccountMode bool, startDate time.Time, endDate time.Time) ([]models.TrialBalanceSheetAccountDetail, error) {
	ret := m.Called(shopId, accountGroup, includeCloseAccountMode, startDate, endDate)
	return ret.Get(0).([]models.TrialBalanceSheetAccountDetail), ret.Error(1)
}

func (m *MockJournalreportRepository) GetDataProfitAndLoss(shopId string, accountGroup string, includeCloseAccountMode bool, startDate time.Time, endDate time.Time) ([]models.ProfitAndLossSheetAccountDetail, error) {
	ret := m.Called(shopId, accountGroup, includeCloseAccountMode, startDate, endDate)
	return ret.Get(0).([]models.ProfitAndLossSheetAccountDetail), ret.Error(1)
}

func (m *MockJournalreportRepository) GetDataBalanceSheet(shopId string, accountGroup string, includeCloseAccountMode bool, endDate time.Time) ([]models.BalanceSheetAccountDetail, error) {
	ret := m.Called(shopId, accountGroup, includeCloseAccountMode, endDate)
	return ret.Get(0).([]models.BalanceSheetAccountDetail), ret.Error(1)
}

func (m *MockJournalreportRepository) GetDataLedgerAccount(shopId string, accountGroup string, consolidateAccountCode string, accountCodeRanges []models.LedgerAccountCodeRange, startDate time.Time, endDate time.Time) ([]models.LedgerAccountRaw, error) {
	ret := m.Called(shopId, accountGroup, consolidateAccountCode, accountCodeRanges, startDate, endDate)
	return ret.Get(0).([]models.LedgerAccountRaw), ret.Error(1)
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
	get, err := service.ProcessBalanceSheetReport("TESTSHOP", "01", false, endDate)
	get.ReportDate = fixReportDate
	assert.Nil(t, err, "Error should be nil")

	assert.Equal(t, get, want, "Process BalanceSheet Report Not Match")
}

func TestLedgerAccount(t *testing.T) {
	repo := new(MockJournalreportRepository)
	repo.On("GetDataLedgerAccount", "TESTSHOP", []models.LedgerAccountCodeRange{
		{
			Start: "100000",
			End:   "150000",
		},
	}, mocktest.MockTime(), mocktest.MockTime()).Return([]models.LedgerAccountRaw{
		{
			RowMode:      -1,
			DocNo:        "",
			DocDate:      mocktest.MockTime(),
			AccountCode:  "AC001",
			AccountName:  "AC Name 1",
			DebitAmount:  0,
			CreditAmount: 0,
			Amount:       75,
		},
		{
			RowMode:      0,
			DocNo:        "DOC001",
			DocDate:      mocktest.MockTime(),
			AccountCode:  "AC001",
			AccountName:  "AC Name 1",
			DebitAmount:  50,
			CreditAmount: 0,
			Amount:       0,
		},
		{
			RowMode:      0,
			DocNo:        "DOC002",
			DocDate:      mocktest.MockTime(),
			AccountCode:  "AC001",
			AccountName:  "AC Name 1",
			DebitAmount:  50,
			CreditAmount: 0,
			Amount:       0,
		},
		{
			RowMode:      -1,
			DocNo:        "",
			DocDate:      mocktest.MockTime(),
			AccountCode:  "AC002",
			AccountName:  "AC Name 2",
			DebitAmount:  0,
			CreditAmount: 0,
			Amount:       200,
		},
		{
			RowMode:      0,
			DocNo:        "DOC003",
			DocDate:      mocktest.MockTime(),
			AccountCode:  "AC002",
			AccountName:  "AC Name 2",
			DebitAmount:  0,
			CreditAmount: 250,
			Amount:       0,
		},
		{
			RowMode:      -1,
			DocNo:        "",
			DocDate:      mocktest.MockTime(),
			AccountCode:  "AC003",
			AccountName:  "AC Name 3",
			DebitAmount:  0,
			CreditAmount: 0,
			Amount:       -50,
		},
		{
			RowMode:      0,
			DocNo:        "",
			DocDate:      mocktest.MockTime(),
			AccountCode:  "AC003",
			AccountName:  "AC Name 3",
			DebitAmount:  100,
			CreditAmount: 0,
			Amount:       0,
		},
		{
			RowMode:      -1,
			DocNo:        "",
			DocDate:      mocktest.MockTime(),
			AccountCode:  "AC004",
			AccountName:  "AC Name 4",
			DebitAmount:  0,
			CreditAmount: 0,
			Amount:       -50,
		},
		{
			RowMode:      0,
			DocNo:        "",
			DocDate:      mocktest.MockTime(),
			AccountCode:  "AC004",
			AccountName:  "AC Name 4",
			DebitAmount:  0,
			CreditAmount: 100,
			Amount:       0,
		},
	}, nil)

	service := journalreport.NewJournalReportService(repo)
	docList, err := service.ProcessLedgerAccount("TESTSHOP", "", "", []models.LedgerAccountCodeRange{
		{
			Start: "100000",
			End:   "150000",
		},
	}, mocktest.MockTime(), mocktest.MockTime())

	assert.Nil(t, err)
	assert.Equal(t, 4, len(docList))

	assert.Equal(t, 175.0, docList[0].NextBalance)
	assert.Equal(t, -50.0, docList[1].NextBalance)
	assert.Equal(t, 50.0, docList[2].NextBalance)
	assert.Equal(t, -150.0, docList[3].NextBalance)

}
