package journalreport

import (
	"smlcloudplatform/pkg/models/vfgl"
	"time"
)

type IJournalReportService interface {
	ProcessTrialBalanceSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.TrialBalanceSheetReport, error)
	ProcessProfitAndLossSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.ProfitAndLossSheetReport, error)
	ProcessBalanceSheetReport(shopId string, accountGroup string, endDate time.Time) (*vfgl.BalanceSheetReport, error)
}

type JournalReportService struct {
	repo IJournalReportRepository
}

func NewJournalReportService(repo IJournalReportRepository) JournalReportService {
	return JournalReportService{
		repo: repo,
	}
}

func (svc JournalReportService) ProcessTrialBalanceSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.TrialBalanceSheetReport, error) {
	// mock := MockTrialBalanceSheetReport(shopId, accountGroup, startDate, endDate)
	// return mock, nil
	details, err := svc.repo.GetDataTrialBalance(shopId, accountGroup, startDate, endDate)

	var totalBalanceDebit float64
	var totalBalanceCredit float64
	var totalAmountDebit float64
	var totalAmountCredit float64
	var totalNextBalanceDebit float64
	var totalnextBalanceCredit float64

	for _, v := range details {
		totalBalanceDebit += v.BalanceDebitAmount
		totalBalanceCredit += v.BalanceCreditAmount
		totalAmountDebit += v.DebitAmount
		totalAmountCredit += v.CreditAmount
		totalNextBalanceDebit += v.NextBalanceDebitAmount
		totalnextBalanceCredit += v.NextBalanceCreditAmount
	}

	result := &vfgl.TrialBalanceSheetReport{
		ReportDate:             time.Now(),
		StartDate:              startDate,
		EndDate:                endDate,
		AccountGroup:           accountGroup,
		AccountDetails:         &details,
		TotalBalanceDebit:      totalBalanceDebit,
		TotalBalanceCredit:     totalBalanceCredit,
		TotalAmountDebit:       totalAmountDebit,
		TotalAmountCredit:      totalAmountCredit,
		TotalNextBalanceDebit:  totalNextBalanceDebit,
		TotalNextBalanceCredit: totalnextBalanceCredit,
	}
	return result, err
}

func (svc JournalReportService) ProcessProfitAndLossSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.ProfitAndLossSheetReport, error) {
	// mock := MockProfitAndLossSheetReport(shopId, accountGroup, startDate, endDate)
	// return mock, nil
	details, err := svc.repo.GetDataProfitAndLoss(shopId, accountGroup, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// var totalData = len(details)
	//fmt.Printf("rows: %v", rows)
	//fmt.Printf("details: %+v\n", details)
	var incomeAmount float64 = 0
	var expenseAmount float64 = 0
	var profitAndLossAmount float64 = 0

	var incomes []vfgl.ProfitAndLossSheetAccountDetail
	var expenses []vfgl.ProfitAndLossSheetAccountDetail

	for _, v := range details {
		if v.AccountCategory == 4 {
			incomes = append(incomes, v)
			incomeAmount += v.Amount
		} else {
			incomes = append(incomes, v)
			expenseAmount += v.Amount
		}
	}

	profitAndLossAmount = incomeAmount - expenseAmount

	result := &vfgl.ProfitAndLossSheetReport{
		ReportDate:          time.Now(),
		StartDate:           startDate,
		EndDate:             endDate,
		AccountGroup:        accountGroup,
		Incomes:             &incomes,
		Expenses:            &expenses,
		TotalIncomeAmount:   incomeAmount,
		TotalExpenseAmount:  expenseAmount,
		ProfitAndLossAmount: profitAndLossAmount,
	}

	return result, nil
}

func (svc JournalReportService) ProcessBalanceSheetReport(shopId string, accountGroup string, endDate time.Time) (*vfgl.BalanceSheetReport, error) {
	// mock := MockBalanceSheetReport(shopId, accountGroup, endDate)
	// return mock, nil
	details, err := svc.repo.GetDataBalanceSheet(shopId, accountGroup, endDate)
	if err != nil {
		return nil, err
	}

	// var totalData = len(details)
	//fmt.Printf("rows: %v", rows)
	//fmt.Printf("details: %+v\n", details)
	var totalAssetAmount float64 = 0
	var totalLiabilityAmount float64 = 0
	var totalOwnersEquityAmount float64 = 0
	var totalLiabilityAndOwnersEquityAmount float64 = 0

	var totalIncome float64 = 0
	var totalExpense float64 = 0
	var totalProfitAndLoss float64 = 0

	var assets []vfgl.BalanceSheetAccountDetail
	var liabilities []vfgl.BalanceSheetAccountDetail
	var ownesEquities []vfgl.BalanceSheetAccountDetail

	for _, v := range details {

		if v.AccountCategory <= 3 {
			if v.AccountCategory == 1 {
				assets = append(assets, v)
				totalAssetAmount += v.Amount
			} else if v.AccountCategory == 2 {
				liabilities = append(liabilities, v)
				totalLiabilityAmount += v.Amount
			} else {
				totalOwnersEquityAmount += v.Amount
				ownesEquities = append(ownesEquities, v)
			}
		} else {
			if v.AccountCategory == 4 {
				totalIncome += v.Amount
			} else {
				totalExpense += v.Amount
			}
		}
	}

	totalProfitAndLoss = totalIncome - totalExpense
	if totalProfitAndLoss != 0 {
		totalOwnersEquityAmount += totalProfitAndLoss
		ownesEquities = append(ownesEquities, vfgl.BalanceSheetAccountDetail{
			ChartOfAccountPG: vfgl.ChartOfAccountPG{
				AccountName:     "กำไร (ขาดทุน) สุทธิ",
				AccountCategory: 3,
			},
			Amount: totalProfitAndLoss,
		})
	}
	totalLiabilityAndOwnersEquityAmount = totalLiabilityAmount + totalOwnersEquityAmount

	result := &vfgl.BalanceSheetReport{
		ReportDate:                          time.Now(),
		EndDate:                             endDate,
		AccountGroup:                        accountGroup,
		Assets:                              &assets,
		Liabilities:                         &liabilities,
		OwnesEquities:                       &ownesEquities,
		TotalAssetAmount:                    totalAssetAmount,
		TotalLiabilityAmount:                totalLiabilityAmount,
		TotalOwnersEquityAmount:             totalOwnersEquityAmount,
		TotalLiabilityAndOwnersEquityAmount: totalLiabilityAndOwnersEquityAmount,
	}

	return result, nil
}
