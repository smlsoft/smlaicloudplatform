package vfgl

import "time"

type TrialBalanceSheetReport struct {
	// วันที่ทำรายการ
	ReportDate time.Time
	// วันที่เริ่มต้น
	StartDate time.Time
	// วันที่สิ้นสุด
	EndDate time.Time
	// เล่มบัญชี
	AccountGroup string
	// รายละเอียดบัญชี
	AccountDetail *[]TrialBalanceSheetAccountDetail
	// รวมยอดยกมาเดบิต
	TotalBalanceDebit float64
	// รวมยอดยกมาเครดิต
	TotalBalanceCredit float64
	// รวมเดบิต
	TotalAmountDebit float64
	// รวมเครดิต
	TotalAmountCredit float64
	// รวมยอดสะสมเดบิต
	TotalNextBalanceDebit float64
	// รวมยอดสะสมเครดิต
	TotalNextBalanceCredit float64
}

type TrialBalanceSheetAccountDetail struct {
	ChartOfAccount
	// ยอดคงเหลือ(ประจำงวด)
	Amount float64 ``
	// ยอดคงเหลือยกมา
	BalanceAmount float64 ``
	// ยอดคงเหลือสะสม
	NextBalanceAmount float64 ``
	// ยอดเดบิต
	DebitAmount float64 ``
	// ยอดเครดิต
	CreditAmount float64 ``
	// ยอดยกมาเดบิต
	BalanceDebitAmount float64 ``
	// ยอดยกมาเครดิต
	BalanceCreditAmount float64 ``
	// ยอดสะสมเดบิต
	NextBalanceDebitAmount float64 ``
	// ยอดสะสมเครดิต
	NextBalanceCreditAmount float64 ``
}

type BalanceSheetReport struct {
	ReportDate                           time.Time
	EndDate                              time.Time
	AccountGroup                         string
	Assets                               *[]BalanceSheetAccountDetail
	Liabilities                          *[]BalanceSheetAccountDetail
	OwnesQutitys                         *[]BalanceSheetAccountDetail
	TotalAssetAmount                     float64
	TotalLiabilityAmount                 float64
	TotalOwnersEqutityAmount             float64
	TotalLiabilityAndOwnersEqutityAmount float64
}

type BalanceSheetAccountDetail struct {
	ChartOfAccount
	// มูลค่า
	Amount float64 ``
}

type ProfitAndLossSheetReport struct {
	// วันที่ทำรายการ
	ReportDate time.Time
	// วันที่เริ่มต้น
	StartDate time.Time
	// วันที่สิ้นสุด
	EndDate time.Time
	// เล่มบัญชี
	AccountGroup string
	// รายการรายได้
	Incomes *[]ProfitAndLossSheetAccountDetail
	// รายการค่าใช้จ่าย
	Expenses *[]ProfitAndLossSheetAccountDetail
	// รวมรายได้
	TotalIncomeAmount float64
	// รวมค่าใช้จ่าย
	TotalExpenseAmount float64
	// กำไรขาดทุน
	ProfitAndLossAmount float64
}

type ProfitAndLossSheetAccountDetail struct {
	ChartOfAccount
	// มูลค่า
	Amount float64 ``
}

type TrialBalanceSheetReportResponse struct {
	Success bool                    `json:"success"`
	Data    TrialBalanceSheetReport `json:"data,omitempty"`
}

type BalanceSheetReportResponse struct {
	Success bool               `json:"success"`
	Data    BalanceSheetReport `json:"data,omitempty"`
}

type LostAndProfitSheetReportResponse struct {
	Success bool                     `json:"success"`
	Data    ProfitAndLossSheetReport `json:"data,omitempty"`
}
