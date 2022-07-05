package journalreport

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/vfgl/journalreport/models"
	"time"
)

type IJournalReportRepository interface {
	GetDataTrialBalance(shopId string, accountGroup string, startDate time.Time, endDate time.Time) ([]models.TrialBalanceSheetAccountDetail, error)
	GetDataProfitAndLoss(shopId string, accountGroup string, startDate time.Time, endDate time.Time) ([]models.ProfitAndLossSheetAccountDetail, error)
	GetDataBalanceSheet(shopId string, accountGroup string, endDate time.Time) ([]models.BalanceSheetAccountDetail, error)
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

func (repo JournalReportRepository) GetDataTrialBalance(shopId string, accountGroup string, startDate time.Time, endDate time.Time) ([]models.TrialBalanceSheetAccountDetail, error) {

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
			, case when(accountcategory = 1 or accountcategory = 5) then coalesce(bal.debitamount, 0)-coalesce(bal.creditamount, 0)
				else coalesce(bal.creditamount, 0)-coalesce(bal.debitamount, 0)
				end as balanceamount
			, case when(accountcategory = 1 or accountcategory = 5) then coalesce(prd.debitamount, 0)-coalesce(prd.creditamount, 0)
				else coalesce(prd.creditamount, 0)-coalesce(prd.debitamount, 0)
				end as amount
			, case when(accountcategory = 1 or accountcategory = 5) then coalesce(nex.debitamount, 0)-coalesce(nex.creditamount, 0)
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

	var details []models.TrialBalanceSheetAccountDetail

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

	return details, nil
}

func (repo JournalReportRepository) GetDataProfitAndLoss(shopId string, accountGroup string, startDate time.Time, endDate time.Time) ([]models.ProfitAndLossSheetAccountDetail, error) {

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
			, case when(accountcategory = 5) then coalesce(prd.debitamount, 0)-coalesce(prd.creditamount, 0)
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

	var details []models.ProfitAndLossSheetAccountDetail

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

	return details, nil
}

func (repo JournalReportRepository) GetDataBalanceSheet(shopId string, accountGroup string, endDate time.Time) ([]models.BalanceSheetAccountDetail, error) {
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
			, case when(accountcategory = 1 or accountcategory = 5) then coalesce(nex.debitamount, 0)-coalesce(nex.creditamount, 0)
				else coalesce(nex.creditamount, 0)-coalesce(nex.debitamount, 0)
				end as amount		
			from chartofaccounts as chart
			left join nex on nex.accountcode = chart.accountcode
			where chart.shopid= @shopid 
		)
		select shopid, parid
        , accountcode, accountname, accountcategory, accountbalancetype
        , accountgroup, accountlevel, consolidateaccountcode
		, amount
		from journal_sheet_sum
		where amount <> 0 
		order by accountcode
	`

	var details []models.BalanceSheetAccountDetail

	condition := map[string]interface{}{
		"shopid":       shopId,
		"accountgroup": accountGroup,
		"enddate":      endDate,
	}

	_, err := repo.pst.Raw(query, condition, &details)
	if err != nil {
		return nil, err
	}

	return details, nil
}
