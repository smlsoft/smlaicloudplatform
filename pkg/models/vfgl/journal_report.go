package vfgl

import "time"

type TrialBalanceSheetReport struct {
	// วันที่ทำรายการ
	ReportDate time.Time `json:"reportdate"`
	// วันที่เริ่มต้น
	StartDate time.Time `json:"startdate"`
	// วันที่สิ้นสุด
	EndDate time.Time `json:"enddate"`
	// เล่มบัญชี
	AccountGroup string `json:"accountgroup"`
	// รายละเอียดบัญชี
	AccountDetails *[]TrialBalanceSheetAccountDetail `json:"accountdetails"`
	// รวมยอดยกมาเดบิต
	TotalBalanceDebit float64 `json:"totalbalancedebit"`
	// รวมยอดยกมาเครดิต
	TotalBalanceCredit float64 `json:"totalbalancecredit"`
	// รวมเดบิต
	TotalAmountDebit float64 `json:"totalamountdebit"`
	// รวมเครดิต
	TotalAmountCredit float64 `json:"totalamountcredit"`
	// รวมยอดสะสมเดบิต
	TotalNextBalanceDebit float64 `json:"totalnextbalancedebit"`
	// รวมยอดสะสมเครดิต
	TotalNextBalanceCredit float64 `json:"totalnextbalancecredit"`
}

type TrialBalanceSheetAccountDetail struct {
	ChartOfAccount
	// ยอดคงเหลือ(ประจำงวด)
	Amount float64 `json:"amount"`
	// ยอดคงเหลือยกมา
	BalanceAmount float64 `json:"balanceamount"`
	// ยอดคงเหลือสะสม
	NextBalanceAmount float64 `json:"nextbalanceamount"`
	// ยอดเดบิต
	DebitAmount float64 `json:"debitamount"`
	// ยอดเครดิต
	CreditAmount float64 `json:"creditamount"`
	// ยอดยกมาเดบิต
	BalanceDebitAmount float64 `json:"balancedebitamount"`
	// ยอดยกมาเครดิต
	BalanceCreditAmount float64 `json:"balancecreditamount"`
	// ยอดสะสมเดบิต
	NextBalanceDebitAmount float64 `json:"nextbalancedebitamount"`
	// ยอดสะสมเครดิต
	NextBalanceCreditAmount float64 `json:"nextbalancecreditamount"`
}

type BalanceSheetReport struct {
	// วันที่ทำรายการ
	ReportDate time.Time `json:"reportdate"`
	// วันที่สิ้นสุด
	EndDate time.Time `json:"enddate"`
	// เล่มบัญชี
	AccountGroup string `json:"accountgroup"`
	// สินทรัพย์
	Assets *[]BalanceSheetAccountDetail `json:"assets"`
	// หนี้สิน
	Liabilities *[]BalanceSheetAccountDetail `json:"liabilities"`
	// ทุนและส่วนของเจ้าของ
	OwnesEquities *[]BalanceSheetAccountDetail `json:"ownesequities"`
	// รวมสินทรัพย์
	TotalAssetAmount float64 `json:"totalassetamount"`
	// รวมหนี้สิน
	TotalLiabilityAmount float64 `json:"totalliabilityamount"`
	// รวมทุนและส่วนของเจ้าของ
	TotalOwnersEquityAmount float64 `json:"totalownersequityamount"`
	// รวมหนี้สิน ทุน และส่วนของเจ้าของ
	TotalLiabilityAndOwnersEquityAmount float64 `json:"totalliabilityandownersequityamount"`
}

type BalanceSheetAccountDetail struct {
	ChartOfAccount
	// มูลค่า
	Amount float64 `json:"amount"`
}

type ProfitAndLossSheetReport struct {
	// วันที่ทำรายการ
	ReportDate time.Time `json:"reportdate"`
	// วันที่เริ่มต้น
	StartDate time.Time `json:"startdate"`
	// วันที่สิ้นสุด
	EndDate time.Time `json:"enddate"`
	// เล่มบัญชี
	AccountGroup string `json:"accountgroup"`
	// รายการรายได้
	Incomes *[]ProfitAndLossSheetAccountDetail `json:"incomes"`
	// รายการค่าใช้จ่าย
	Expenses *[]ProfitAndLossSheetAccountDetail `json:"expenses"`
	// รวมรายได้
	TotalIncomeAmount float64 `json:"totalincomeamount"`
	// รวมค่าใช้จ่าย
	TotalExpenseAmount float64 `json:"totalexpenseamount"`
	// กำไรขาดทุน
	ProfitAndLossAmount float64 `json:"profitandlossamount"`
}

type ProfitAndLossSheetAccountDetail struct {
	ChartOfAccount
	// มูลค่า
	Amount float64 `json:"amount"`
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
