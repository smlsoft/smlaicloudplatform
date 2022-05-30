package journalreport

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/vfgl"
	"time"
)

type IJournalReportRepository interface {
	GetDataTrialBalance(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.TrialBalanceSheetReport, error)
	GetDataProfitAndLoss(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.ProfitAndLossSheetReport, error)
	GetDataBalanceSheet(shopId string, accountGroup string, endDate time.Time) (*vfgl.BalanceSheetReport, error)
}

type JournalReportRepository struct {
	pst microservice.IPersister
}

func NewJournalReportRepository(pst microservice.IPersister) JournalReportRepository {
	return JournalReportRepository{
		pst: pst,
	}

}

/* Full Query

-- REPORT TRIAL BALANCE SHEET

WITH journal_doc as (
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			--, acc.accountname , acc.accountcategory, acc.accountbalancetype
			, d.debitamount ,d.creditamount
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			where h.shopid= '27dcEdktOoaSBYFmnN6G6ett4Jb' and h.accountgroup = '01'
			--left join chartofaccounts as acc on acc.shopid = d.shopid and acc.accountcode = d.accountcode
		)
		, bal as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount
			from journal_doc where journal_doc.docdate < '2022-05-01'
			group by accountcode
		)
		, prd as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount
			from journal_doc where journal_doc.docdate between '2022-05-01' and '2022-05-31'
			group by accountcode
		)
		, nex as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount
			from journal_doc where journal_doc.docdate <= '2022-05-31'
			group by accountcode
		)
		, journal_sheet_sum as (
			select chart.shopid, chart.par_id
            , chart.accountcode, chart.accountname
			, chart.accountcategory, chart.accountbalancetype, chart.accountgroup, chart.accountlevel, chart.consolidateaccountcode
			, coalesce(bal.debitamount, 0) as balancedebitamount, coalesce(bal.creditamount, 0) as balancecreditamount
			, coalesce(prd.debitamount, 0) as debitamount, coalesce(prd.creditamount, 0) as creditamount
			, coalesce(nex.debitamount, 0) as nextbalancedebitamount, coalesce(nex.creditamount, 0) as nextbalancecreditamount
			, case when(accountbalancetype = 1) then coalesce(bal.debitamount, 0)-coalesce(bal.creditamount, 0)
				else coalesce(bal.creditamount, 0)-coalesce(bal.debitamount, 0)
				end as balanceamount
			, case when(accountbalancetype = 1) then coalesce(prd.debitamount, 0)-coalesce(prd.creditamount, 0)
				else coalesce(prd.creditamount, 0)-coalesce(prd.debitamount, 0)
				end as amount
			, case when(accountbalancetype = 1) then coalesce(nex.debitamount, 0)-coalesce(nex.creditamount, 0)
				else coalesce(nex.creditamount, 0)-coalesce(nex.debitamount, 0)
				end as nextbalanceamount
			from chartofaccounts as chart
			left join bal on bal.accountcode = chart.accountcode
			left join prd on prd.accountcode = chart.accountcode
			left join nex on nex.accountcode = chart.accountcode
			where chart.shopid= '27dcEdktOoaSBYFmnN6G6ett4Jb'
		)
		select shopid, par_id
        , accountcode, accountname, accountcategory, accountbalancetype
        , accountgroup, accountlevel, consolidateaccountcode
		, balancedebitamount, balancecreditamount, debitamount, creditamount, nextbalancedebitamount, nextbalancecreditamount
		, balanceamount, amount, nextbalanceamount
		from journal_sheet_sum
		where balanceamount <> 0 or amount <> 0 or nextbalanceamount <> 0
		order by accountcode

*/

func (repo JournalReportRepository) GetDataTrialBalance(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.TrialBalanceSheetReport, error) {

	query := `
	WITH journal_doc as (
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			--, acc.accountname , acc.accountcategory, acc.accountbalancetype
			, d.debitamount ,d.creditamount
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			where h.shopid= @shopid and h.accountgroup = @accountgroup
			--left join chartofaccounts as acc on acc.shopid = d.shopid and acc.accountcode = d.accountcode
		)
		, bal as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount 
			from journal_doc where journal_doc.docdate < @startdate
			group by accountcode
		)
		, prd as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount 
			from journal_doc where journal_doc.docdate between @startdate and @enddate
			group by accountcode
		)
		, nex as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount 
			from journal_doc where journal_doc.docdate <= @enddate
			group by accountcode
		)
		, journal_sheet_sum as (
			select chart.shopid, chart.parid
            , chart.accountcode, chart.accountname
			, chart.accountcategory, chart.accountbalancetype, chart.accountgroup, chart.accountlevel, chart.consolidateaccountcode
			, coalesce(bal.debitamount, 0) as balancedebitamount, coalesce(bal.creditamount, 0) as balancecreditamount
			, coalesce(prd.debitamount, 0) as debitamount, coalesce(prd.creditamount, 0) as creditamount
			, coalesce(nex.debitamount, 0) as nextbalancedebitamount, coalesce(nex.creditamount, 0) as nextbalancecreditamount
			, case when(accountbalancetype = 1) then coalesce(bal.debitamount, 0)-coalesce(bal.creditamount, 0)
				else coalesce(bal.creditamount, 0)-coalesce(bal.debitamount, 0)
				end as balanceamount
			, case when(accountbalancetype = 1) then coalesce(prd.debitamount, 0)-coalesce(prd.creditamount, 0)
				else coalesce(prd.creditamount, 0)-coalesce(prd.debitamount, 0)
				end as amount
			, case when(accountbalancetype = 1) then coalesce(nex.debitamount, 0)-coalesce(nex.creditamount, 0)
				else coalesce(nex.creditamount, 0)-coalesce(nex.debitamount, 0)
				end as nextbalanceamount
			from chartofaccounts as chart
			left join bal on bal.accountcode = chart.accountcode
			left join prd on prd.accountcode = chart.accountcode
			left join nex on nex.accountcode = chart.accountcode
			where chart.shopid= @shopid
		)
		select shopid, parid
        , accountcode, accountname, accountcategory, accountbalancetype
        , accountgroup, accountlevel, consolidateaccountcode
		, balancedebitamount, balancecreditamount, debitamount, creditamount, nextbalancedebitamount, nextbalancecreditamount
		, balanceamount, amount, nextbalanceamount
		from journal_sheet_sum
		where balanceamount <> 0 or amount <> 0 or nextbalanceamount <> 0
		order by accountcode
	`

	var details []vfgl.TrialBalanceSheetAccountDetail

	condition := map[string]interface{}{
		"shopid":       shopId,
		"accountgroup": accountGroup,
		"startdate":    startDate,
		"enddate":      endDate,
	}

	_, err := repo.pst.Raw(query, condition, &details)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("rows: %v", rows)
	//fmt.Printf("details: %+v\n", details)
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

	return result, nil
}

func (repo JournalReportRepository) GetDataProfitAndLoss(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.ProfitAndLossSheetReport, error) {

	query := `
	WITH journal_doc as (
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			, d.debitamount ,d.creditamount
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			where h.shopid= @shopid and h.accountgroup = @accountgroup
		)
		, prd as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount 
			from journal_doc where journal_doc.docdate between @startdate and @enddate
			group by accountcode
		)
		, journal_sheet_sum as (
			select chart.shopid, chart.parid
            , chart.accountcode, chart.accountname
			, chart.accountcategory, chart.accountbalancetype, chart.accountgroup, chart.accountlevel, chart.consolidateaccountcode
			, coalesce(prd.debitamount, 0) as debitamount, coalesce(prd.creditamount, 0) as creditamount			
			, case when(accountbalancetype = 1) then coalesce(prd.debitamount, 0)-coalesce(prd.creditamount, 0)
				else coalesce(prd.creditamount, 0)-coalesce(prd.debitamount, 0)
				end as amount			
			from chartofaccounts as chart
			left join prd on prd.accountcode = chart.accountcode
			where chart.shopid= @shopid and chart.accountcategory in (4,5)
		)
		select shopid, parid
        , accountcode, accountname, accountcategory, accountbalancetype
        , accountgroup, accountlevel, consolidateaccountcode
		, debitamount, creditamount, amount
		from journal_sheet_sum
		where amount <> 0 
		order by accountcode
	`

	var details []vfgl.ProfitAndLossSheetAccountDetail

	condition := map[string]interface{}{
		"shopid":       shopId,
		"accountgroup": accountGroup,
		"startdate":    startDate,
		"enddate":      endDate,
	}

	_, err := repo.pst.Raw(query, condition, &details)
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

func (repo JournalReportRepository) GetDataBalanceSheet(shopId string, accountGroup string, endDate time.Time) (*vfgl.BalanceSheetReport, error) {
	query := `
	WITH journal_doc as (
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			, d.debitamount ,d.creditamount
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			where h.shopid= @shopid and h.accountgroup = @accountgroup
		)
		, nex as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount 
			from journal_doc where journal_doc.docdate <= @enddate
			group by accountcode
		)
		, journal_sheet_sum as (
			select chart.shopid, chart.parid
            , chart.accountcode, chart.accountname
			, chart.accountcategory, chart.accountbalancetype, chart.accountgroup, chart.accountlevel, chart.consolidateaccountcode
			, case when(accountbalancetype = 1) then coalesce(nex.debitamount, 0)-coalesce(nex.creditamount, 0)
				else coalesce(nex.creditamount, 0)-coalesce(nex.debitamount, 0)
				end as amount		
			from chartofaccounts as chart
			left join nex on nex.accountcode = chart.accountcode
			where chart.shopid= @shopid and chart.accountcategory <= 3 
		)
		select shopid, parid
        , accountcode, accountname, accountcategory, accountbalancetype
        , accountgroup, accountlevel, consolidateaccountcode
		, amount
		from journal_sheet_sum
		where amount <> 0 
		order by accountcode
	`

	var details []vfgl.BalanceSheetAccountDetail

	condition := map[string]interface{}{
		"shopid":       shopId,
		"accountgroup": accountGroup,
		"enddate":      endDate,
	}

	_, err := repo.pst.Raw(query, condition, &details)
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

	var assets []vfgl.BalanceSheetAccountDetail
	var liabilities []vfgl.BalanceSheetAccountDetail
	var ownesEquities []vfgl.BalanceSheetAccountDetail

	for _, v := range details {
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
