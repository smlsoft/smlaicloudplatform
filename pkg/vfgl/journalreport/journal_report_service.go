package journalreport

import (
	"fmt"
	chartofaccountModel "smlcloudplatform/pkg/vfgl/chartofaccount/models"
	"smlcloudplatform/pkg/vfgl/journalreport/models"
	"time"
)

type IJournalReportService interface {
	ProcessTrialBalanceSheetReport(shopId string, accountGroup string, includeCloseAccountMode bool, startDate time.Time, endDate time.Time) (*models.TrialBalanceSheetReport, error)
	ProcessProfitAndLossSheetReport(shopId string, accountGroup string, includeCloseAccountMode bool, startDate time.Time, endDate time.Time) (*models.ProfitAndLossSheetReport, error)
	ProcessBalanceSheetReport(shopId string, accountGroup string, includeCloseAccountMode bool, endDate time.Time) (*models.BalanceSheetReport, error)
	ProcessLedgerAccount(shopId string, accountGroup string, consolidateAccountCode string, accountRanges []models.LedgerAccountCodeRange, startDate time.Time, endDate time.Time) ([]models.LedgerAccount, error)
}

type JournalReportService struct {
	repo IJournalReportRepository
}

func NewJournalReportService(repo IJournalReportRepository) JournalReportService {
	return JournalReportService{
		repo: repo,
	}
}

func (svc JournalReportService) ProcessTrialBalanceSheetReport(shopId string, accountGroup string, includeCloseAccountMode bool, startDate time.Time, endDate time.Time) (*models.TrialBalanceSheetReport, error) {
	// mock := MockTrialBalanceSheetReport(shopId, accountGroup, startDate, endDate)
	// return mock, nil
	details, err := svc.repo.GetDataTrialBalance(shopId, accountGroup, includeCloseAccountMode, startDate, endDate)

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

	result := &models.TrialBalanceSheetReport{
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

func (svc JournalReportService) ProcessProfitAndLossSheetReport(shopId string, accountGroup string, includeCloseAccountMode bool, startDate time.Time, endDate time.Time) (*models.ProfitAndLossSheetReport, error) {
	// mock := MockProfitAndLossSheetReport(shopId, accountGroup, startDate, endDate)
	// return mock, nil
	details, err := svc.repo.GetDataProfitAndLoss(shopId, accountGroup, includeCloseAccountMode, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// var totalData = len(details)
	//fmt.Printf("rows: %v", rows)
	//fmt.Printf("details: %+v\n", details)
	var incomeAmount float64 = 0
	var expenseAmount float64 = 0
	var profitAndLossAmount float64 = 0

	var incomes []models.ProfitAndLossSheetAccountDetail
	var expenses []models.ProfitAndLossSheetAccountDetail

	for _, v := range details {
		if v.AccountCategory == 4 {
			incomes = append(incomes, v)
			incomeAmount += v.Amount
		} else {
			expenses = append(expenses, v)
			expenseAmount += v.Amount
		}
	}

	profitAndLossAmount = incomeAmount - expenseAmount

	result := &models.ProfitAndLossSheetReport{
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

func (svc JournalReportService) ProcessBalanceSheetReport(shopId string, accountGroup string, includeCloseAccountMode bool, endDate time.Time) (*models.BalanceSheetReport, error) {
	// mock := MockBalanceSheetReport(shopId, accountGroup, endDate)
	// return mock, nil
	details, err := svc.repo.GetDataBalanceSheet(shopId, accountGroup, includeCloseAccountMode, endDate)
	if err != nil {
		return nil, err
	}

	// var totalData = len(details)
	fmt.Printf("rows: %v", len(details))
	//fmt.Printf("details: %+v\n", details)
	var totalAssetAmount float64 = 0
	var totalLiabilityAmount float64 = 0
	var totalOwnersEquityAmount float64 = 0
	var totalLiabilityAndOwnersEquityAmount float64 = 0

	var totalIncome float64 = 0
	var totalExpense float64 = 0
	var totalProfitAndLoss float64 = 0

	var assets []models.BalanceSheetAccountDetail
	var liabilities []models.BalanceSheetAccountDetail
	var ownesEquities []models.BalanceSheetAccountDetail

	for _, v := range details {

		// fmt.Printf("%+v\n", v)

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
		ownesEquities = append(ownesEquities, models.BalanceSheetAccountDetail{
			ChartOfAccountPG: chartofaccountModel.ChartOfAccountPG{
				AccountName:     "กำไร (ขาดทุน) สุทธิ",
				AccountCategory: 3,
			},
			Amount: totalProfitAndLoss,
		})
	}
	totalLiabilityAndOwnersEquityAmount = totalLiabilityAmount + totalOwnersEquityAmount

	result := &models.BalanceSheetReport{
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

func (svc JournalReportService) ProcessLedgerAccount(shopID string, accountGroup string, consolidateAccountCode string, accountRanges []models.LedgerAccountCodeRange, startDate time.Time, endDate time.Time) ([]models.LedgerAccount, error) {

	rawDocList, err := svc.repo.GetDataLedgerAccount(shopID, accountGroup, consolidateAccountCode, accountRanges, startDate, endDate)

	if err != nil {
		return nil, err
	}

	docList := []models.LedgerAccount{}

	lastAccountCode := ""
	lastAmount := 0.0
	tempDoc := models.LedgerAccount{}

	currentIndexAccount := -1
	for _, doc := range rawDocList {

		if lastAccountCode != doc.AccountCode && doc.RowMode == -1 {
			currentIndexAccount++
			tempDoc = models.LedgerAccount{}
			tempDoc.Details = &[]models.LedgerAccountDetail{}
			tempDoc.AccountCode = doc.AccountCode
			tempDoc.AccountName = doc.AccountName
			tempDoc.AccountGroup = doc.AccountGroup
			tempDoc.ConsolidateAccountCode = doc.ConsolidateAccountCode

			lastAmount = doc.Amount
			tempDoc.Balance = lastAmount
			tempDoc.NextBalance = lastAmount

			docList = append(docList, tempDoc)
		}

		if doc.RowMode == 0 {
			lastAmount = (lastAmount + doc.DebitAmount) - doc.CreditAmount
			docList[currentIndexAccount].NextBalance = lastAmount

			detail := models.LedgerAccountDetail{
				DocNo:              doc.DocNo,
				AccountDescription: doc.AccountDescription,
				DocDate:            doc.DocDate,
				Debit:              doc.DebitAmount,
				Credit:             doc.CreditAmount,
				Amount:             lastAmount,
			}
			*tempDoc.Details = append(*tempDoc.Details, detail)
		}

		lastAccountCode = doc.AccountCode
	}

	return docList, nil
}
