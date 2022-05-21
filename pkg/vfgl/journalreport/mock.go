package journalreport

import (
	"smlcloudplatform/pkg/models/vfgl"
	"time"
)

func MockTrialBalanceSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) *vfgl.TrialBalanceSheetReport {

	acc12101 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccount: vfgl.ChartOfAccount{
			AccountCode: "12101",
			AccountName: "เงินฝากธนาคาร",
		},
		Amount: 2200000,
	}
	acc32010 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccount: vfgl.ChartOfAccount{
			AccountCode: "32010",
			AccountName: "ทุน - เงินล้าน",
		},
		Amount: 2300000,
	}
	acc13010 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccount: vfgl.ChartOfAccount{
			AccountCode: "13010",
			AccountName: "ลูกหนี้เงินกู้ - บัญชี 1 (เงินล้าน)",
		},
		Amount: 100000,
	}

	var accountDetail []vfgl.TrialBalanceSheetAccountDetail

	accountDetail = append(accountDetail, *acc12101)
	accountDetail = append(accountDetail, *acc13010)
	accountDetail = append(accountDetail, *acc32010)

	reportMock := &vfgl.TrialBalanceSheetReport{
		AccountGroup:  accountGroup,
		ReportDate:    time.Now(),
		StartDate:     startDate,
		EndDate:       endDate,
		AccountDetail: &accountDetail,
	}

	return reportMock
}

func MockBalanceSheetReport(shopId string, accountGroup string, endDate time.Time) *vfgl.BalanceSheetReport {

	reportMock := &vfgl.BalanceSheetReport{
		AccountGroup: accountGroup,
		ReportDate:   time.Now(),
		EndDate:      endDate,
	}
	return reportMock
}

func MockProfitAndLossSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) *vfgl.ProfitAndLossSheetReport {

	reportMock := &vfgl.ProfitAndLossSheetReport{
		AccountGroup: accountGroup,
		ReportDate:   time.Now(),
		StartDate:    startDate,
		EndDate:      endDate,
	}
	return reportMock
}
