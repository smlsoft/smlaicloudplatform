package journalreport

import (
	"smlcloudplatform/pkg/models/vfgl"
	"time"
)

func MockTrialBalanceSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) *vfgl.TrialBalanceSheetReport {

	acc12101 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode:  "12101",
			AccountName:  "เงินฝากธนาคาร บัญชี 1 (เงินล้าน)",
			AccountGroup: "12000",
		},
		Amount: 311026.03,
	}

	acc13010 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode:  "13010",
			AccountName:  "ลูกหนี้เงินกู้ - บัญชี 1 (เงินล้าน)",
			AccountGroup: "13000",
		},
		Amount: 2600000.0,
	}

	acc32010 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode:  "32010",
			AccountName:  "ทุน - เงินล้าน",
			AccountGroup: "32000",
		},
		Amount: 2300000.0,
	}

	acc33070 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33070",
			AccountName: "เงินสมทบกองทุน",
		},
		Amount: 156834.0,
	}

	acc33060 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33060",
			AccountName: "เงินประกันความเส่ียง",
		},
		Amount: 126605.0,
	}

	acc33090 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33090",
			AccountName: "ค่าดำเนินงาน/ค่าบริหารจัดการ",
		},
		Amount: 1140.0,
	}

	acc33080 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33080",
			AccountName: "เงินสวัสดิการ",
		},
		Amount: 42180.0,
	}

	acc33050 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33050",
			AccountName: "สาธารณะประโยชน์",
		},
		Amount: 2170.0,
	}

	acc32050 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "32050",
			AccountName: "ทุน - อื่น",
		},
		Amount: 100301.0,
	}

	acc34010 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "34010",
			AccountName: "กำไรสะสม (ขาดทุน) สะสม บัญชี 1",
		},
		Amount: 25016.87,
	}

	acc35010 := &vfgl.TrialBalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "35010",
			AccountName: "กำไร ( ขาดทุน ) บัญชี 1",
		},
		Amount: 156679.16,
	}

	var accountDetail []vfgl.TrialBalanceSheetAccountDetail
	accountDetail = append(accountDetail, *acc12101)
	accountDetail = append(accountDetail, *acc13010)
	accountDetail = append(accountDetail, *acc32010)
	accountDetail = append(accountDetail, *acc33070)
	accountDetail = append(accountDetail, *acc33060)
	accountDetail = append(accountDetail, *acc33090)
	accountDetail = append(accountDetail, *acc33080)
	accountDetail = append(accountDetail, *acc33050)
	accountDetail = append(accountDetail, *acc32050)
	accountDetail = append(accountDetail, *acc34010)
	accountDetail = append(accountDetail, *acc35010)

	reportMock := &vfgl.TrialBalanceSheetReport{
		AccountGroup:           accountGroup,
		ReportDate:             time.Now(),
		StartDate:              startDate,
		EndDate:                endDate,
		AccountDetails:         &accountDetail,
		TotalAmountDebit:       2911026.03,
		TotalAmountCredit:      2911026.03,
		TotalNextBalanceDebit:  2911026.03,
		TotalNextBalanceCredit: 2911026.03,
	}

	return reportMock
}

func MockBalanceSheetReport(shopId string, accountGroup string, endDate time.Time) *vfgl.BalanceSheetReport {

	acc11010 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "11010",
			AccountName: "เงินสด - บัญชี 1",
		},
		Amount: 0,
	}

	acc12101 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode:  "12101",
			AccountName:  "เงินฝากธนาคาร บัญชี 1 (เงินล้าน)",
			AccountGroup: "12000",
		},
		Amount: 311026.03,
	}

	acc13010 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode:  "13010",
			AccountName:  "ลูกหนี้เงินกู้ - บัญชี 1 (เงินล้าน)",
			AccountGroup: "13000",
		},
		Amount: 2600000.0,
	}
	var assets []vfgl.BalanceSheetAccountDetail
	assets = append(assets, *acc11010)
	assets = append(assets, *acc12101)
	assets = append(assets, *acc13010)

	acc32010 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode:  "32010",
			AccountName:  "ทุน - เงินล้าน",
			AccountGroup: "32000",
		},
		Amount: 2300000.0,
	}

	acc33070 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33070",
			AccountName: "เงินสมทบกองทุน",
		},
		Amount: 156834.0,
	}

	acc33060 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33060",
			AccountName: "เงินประกันความเส่ียง",
		},
		Amount: 126605.0,
	}

	acc33090 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33090",
			AccountName: "ค่าดำเนินงาน/ค่าบริหารจัดการ",
		},
		Amount: 1140.0,
	}

	acc33080 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33080",
			AccountName: "เงินสวัสดิการ",
		},
		Amount: 42180.0,
	}

	acc33050 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "33050",
			AccountName: "สาธารณะประโยชน์",
		},
		Amount: 2170.0,
	}

	acc32050 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "32050",
			AccountName: "ทุน - อื่น",
		},
		Amount: 100301.0,
	}

	acc34010 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "34010",
			AccountName: "กำไรสะสม (ขาดทุน) สะสม บัญชี 1",
		},
		Amount: 25016.87,
	}

	acc35010 := &vfgl.BalanceSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "35010",
			AccountName: "กำไร ( ขาดทุน ) บัญชี 1",
		},
		Amount: 156679.16,
	}

	var ownersEqutities []vfgl.BalanceSheetAccountDetail
	ownersEqutities = append(assets, *acc32010)
	ownersEqutities = append(assets, *acc33070)
	ownersEqutities = append(assets, *acc33060)
	ownersEqutities = append(assets, *acc33090)
	ownersEqutities = append(assets, *acc33080)
	ownersEqutities = append(assets, *acc33050)
	ownersEqutities = append(assets, *acc32050)
	ownersEqutities = append(assets, *acc34010)
	ownersEqutities = append(assets, *acc35010)

	reportMock := &vfgl.BalanceSheetReport{
		AccountGroup:                        accountGroup,
		ReportDate:                          time.Now(),
		EndDate:                             endDate,
		Assets:                              &assets,
		TotalAssetAmount:                    2911026.03,
		OwnesEquities:                       &ownersEqutities,
		TotalLiabilityAmount:                0,
		TotalOwnersEquityAmount:             2911026.03,
		TotalLiabilityAndOwnersEquityAmount: 2911026.03,
	}
	return reportMock
}

func MockProfitAndLossSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) *vfgl.ProfitAndLossSheetReport {

	acc41010 := &vfgl.ProfitAndLossSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "41010",
			AccountName: "รายได้ - ดอกเบี้ยเงินกู้ - บัญชี 1",
		},
		Amount: 156000.0,
	}

	acc45010 := &vfgl.ProfitAndLossSheetAccountDetail{
		ChartOfAccountPG: vfgl.ChartOfAccountPG{
			AccountCode: "45010",
			AccountName: "รายได้ - ดอกเบี้ยเงินฝากธนาคาร-บัญชี 1",
		},
		Amount: 679.16,
	}

	var incomes []vfgl.ProfitAndLossSheetAccountDetail
	incomes = append(incomes, *acc41010)
	incomes = append(incomes, *acc45010)

	var expenses []vfgl.ProfitAndLossSheetAccountDetail

	reportMock := &vfgl.ProfitAndLossSheetReport{
		AccountGroup:        accountGroup,
		ReportDate:          time.Now(),
		StartDate:           startDate,
		EndDate:             endDate,
		Incomes:             &incomes,
		TotalIncomeAmount:   156679.16,
		Expenses:            &expenses,
		TotalExpenseAmount:  0,
		ProfitAndLossAmount: 156679.16,
	}
	return reportMock
}
