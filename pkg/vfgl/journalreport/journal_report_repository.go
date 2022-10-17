package journalreport

import (
	"fmt"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/vfgl/journalreport/models"
	"strconv"
	"time"
)

type IJournalReportRepository interface {
	GetDataTrialBalance(shopId string, accountGroup string, includeCloseAccountMode bool, startDate time.Time, endDate time.Time) ([]models.TrialBalanceSheetAccountDetail, error)
	GetDataProfitAndLoss(shopId string, accountGroup string, includeCloseAccountMode bool, startDate time.Time, endDate time.Time) ([]models.ProfitAndLossSheetAccountDetail, error)
	GetDataBalanceSheet(shopId string, accountGroup string, includeCloseAccountMode bool, endDate time.Time) ([]models.BalanceSheetAccountDetail, error)
	GetDataLedgerAccount(shopId string, accountGroup string, consolidateAccountCode string, accountRanges []models.LedgerAccountCodeRange, startDate time.Time, endDate time.Time) ([]models.LedgerAccountRaw, error)
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

func (repo JournalReportRepository) GetDataTrialBalance(shopId string, accountGroup string, includeCloseAccountMode bool,
	startDate time.Time, endDate time.Time) ([]models.TrialBalanceSheetAccountDetail, error) {

	var closeDocFilter string

	if includeCloseAccountMode == true {
		closeDocFilter = ` 
		union all 
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			, d.debitamount ,d.creditamount
			, acc.accountcategory
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			left join chartofaccounts as acc on acc.shopid = d.shopid and acc.accountcode = d.accountcode 
            where h.shopid= @shopid and ( h.journaltype = 1 and h.docdate between @startdate and @enddate )
	`
	}
	query := `
	WITH journal_doc as (
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			--, acc.accountname , acc.accountbalancetype
			, d.debitamount ,d.creditamount
			, acc.accountcategory
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			left join chartofaccounts as acc on acc.shopid = d.shopid and acc.accountcode = d.accountcode
			where h.shopid= @shopid and h.accountgroup = @accountgroup  and h.docdate < @enddate 
			 and (( h.journaltype = 0) or (h.journaltype=1 and h.docdate < @startdate ))
		
		` + closeDocFilter + `
			
		)
		, bal as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount 
			from journal_doc where journal_doc.docdate < @startdate
			and accountcategory in (1,2,3)
			group by accountcode
		)
		, prd as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount 
			from journal_doc where journal_doc.docdate between @startdate and @enddate
			group by accountcode
		)
		, nex as (
			select accountcode, sum(debitamount) as debitamount, sum(creditamount) as creditamount 
			from journal_doc 
			where 
			 ( accountcategory in (1,2,3) and journal_doc.docdate <= @enddate ) or ( accountcategory in (4,5) and journal_doc.docdate between @startdate and @enddate )
			group by accountcode
		)
		, journal_sheet_sum as (
			select chart.shopid, chart.parid
            , chart.accountcode, chart.accountname
			, chart.accountcategory, chart.accountbalancetype, chart.accountgroup, chart.accountlevel, chart.consolidateaccountcode
			, coalesce(bal.debitamount, 0) as balancedebitamount, coalesce(bal.creditamount, 0) as balancecreditamount
			, coalesce(prd.debitamount, 0) as debitamount, coalesce(prd.creditamount, 0) as creditamount
			, case when(accountcategory = 1 or accountcategory = 5) then coalesce(nex.debitamount, 0)-coalesce(nex.creditamount, 0) else 0 end as nextbalancedebitamount
            , case when(accountcategory = 1 or accountcategory = 5) then 0 else coalesce(nex.creditamount, 0)-coalesce(nex.debitamount, 0) end as nextbalancecreditamount
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

func (repo JournalReportRepository) GetDataProfitAndLoss(shopId string, accountGroup string, includeCloseAccountMode bool, startDate time.Time, endDate time.Time) ([]models.ProfitAndLossSheetAccountDetail, error) {

	var closeDocFilter string

	if includeCloseAccountMode == true {
		closeDocFilter = ` 
		union all 
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			, d.debitamount ,d.creditamount
			, acc.accountcategory
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			left join chartofaccounts as acc on acc.shopid = d.shopid and acc.accountcode = d.accountcode 

            where h.shopid= @shopid and ( h.journaltype = 1 and h.docdate between @startdate and @enddate )
	`
	}

	query := `
	WITH journal_doc as (
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			, d.debitamount ,d.creditamount
			, acc.accountcategory
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			left join chartofaccounts as acc on acc.shopid = d.shopid and acc.accountcode = d.accountcode 

			where h.shopid= @shopid and h.accountgroup = @accountgroup and h.docdate < @enddate 
			and (( h.journaltype = 0) or (h.journaltype=1 and h.docdate < @startdate ))

		` + closeDocFilter + `
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

func (repo JournalReportRepository) GetDataBalanceSheet(shopId string, accountGroup string, includeCloseAccountMode bool, endDate time.Time) ([]models.BalanceSheetAccountDetail, error) {

	reportYear := endDate.Year()
	var closeDocFilter string

	if includeCloseAccountMode == true {
		closeDocFilter = ` 
		union all 
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			, d.debitamount ,d.creditamount
			, acc.accountcategory
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			left join chartofaccounts as acc on acc.shopid = d.shopid and acc.accountcode = d.accountcode 

            where h.shopid= @shopid and ( h.journaltype = 1 and (extract (year from h.docdate)) = @reportyear )
		`
	}

	/*

		and (( h.journaltype = 0) or (h.journaltype=1 and date(h.docdate) < date(@enddate) ))

				and (
					(acc.accountcategory in (1,2,3)) or
					(
						acc.accountcategory in (4,5) and (extract (year from h.docdate)) = @reportyear
					)
				)
	*/

	query := `
	WITH journal_doc as (
		select h.shopid, h.docno, h.docdate, h.accountyear, h.accountperiod, h.accountgroup
			, d.accountcode
			, d.debitamount ,d.creditamount
			, acc.accountcategory
			from journals_detail as d
			join journals as h on h.shopid = d.shopid and h.docno = d.docno
			left join chartofaccounts as acc on acc.shopid = d.shopid and acc.accountcode = d.accountcode 

			where h.shopid= @shopid and h.accountgroup = @accountgroup  and h.docdate < @enddate 
			and ( 
				acc.accountcategory in (1,2,3) or (acc.accountcategory in (4,5) and (extract (year from h.docdate)) = @reportyear)
			)
			and ( (h.journaltype = 0) or (
				h.journaltype=1 and (extract (year from h.docdate)) < @reportyear
			))
			
		` + closeDocFilter + `
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
		"reportyear":   reportYear,
	}

	_, err := repo.pst.Raw(query, condition, &details)
	if err != nil {
		return nil, err
	}

	fmt.Print("query : " + query + "\n")
	fmt.Printf("details: \n %+v\n", details)

	return details, nil
}

func (repo JournalReportRepository) GetDataLedgerAccount(shopID string, accountGroup string, consolidateAccountCode string, accountRanges []models.LedgerAccountCodeRange, startDate time.Time, endDate time.Time) ([]models.LedgerAccountRaw, error) {

	accountCodeQuery := ""
	accountGroupQuery := ""
	consolidateAccountCodeQuery := ""

	values := map[string]interface{}{
		"shopid":    shopID,
		"startdate": startDate,
		"enddate":   endDate,
	}

	if len(accountRanges) > 0 {

		for idx, accRange := range accountRanges {
			if len(accountCodeQuery) > 0 {
				accountCodeQuery = accountCodeQuery + " or "
			}

			idxStr := strconv.Itoa(idx)
			accountCodeQuery = accountCodeQuery + " a.accountcode between @accountcode" + idxStr + "1 and  @accountcode" + idxStr + "2"

			values["accountcode"+idxStr+"1"] = accRange.Start
			values["accountcode"+idxStr+"2"] = accRange.End

		}
	}

	if len(accountGroup) > 0 {
		accountGroupQuery = " AND a.accountgroup = @accountgroup"
		values["accountgroup"] = accountGroup
	}

	if len(consolidateAccountCode) > 0 {
		consolidateAccountCodeQuery = " AND a.consolidateaccountcode = @consolidateaccountcode"
		values["consolidateaccountcode"] = consolidateAccountCode
	}

	if len(accountCodeQuery) > 0 {
		accountCodeQuery = " and ( " + accountCodeQuery + " ) "
	}

	rawQuery := `select * from (
		WITH 
		acc as ( 
			SELECT  a.accountcode,a.accountname
		from chartofaccounts a  WHERE shopid = @shopid ` + accountGroupQuery + consolidateAccountCodeQuery + accountCodeQuery + ` 
		)
		,
		acc_balance as (
		select d.accountcode,  sum(d.debitamount -  d.creditamount) as amount
		from chartofaccounts a 
		left join journals_detail d on d.shopid = a.shopid AND a.accountcode = d.accountcode 
		left join journals j on j.shopid = d.shopid AND j.docno = d.docno
		where j.docdate < @startdate and a.shopid = @shopid ` + accountGroupQuery + consolidateAccountCodeQuery + accountCodeQuery + ` 
		group by d.accountcode
		)
		SELECT -1 as rowmode, '1900-01-01'::date as docdate, '' as docno,acc.accountcode,acc.accountname, '' as accountdescription,
		0 as debitamount, 0 as creditamount, COALESCE(amount, 0) as amount 
		FROM acc left join acc_balance ON acc.accountcode = acc_balance.accountcode
		union all
		select 0 as rowmode, j.docdate, j.docno, d.accountcode,d.accountname, j.accountdescription as accountdescription, d.debitamount, d.creditamount, 0 as amount 
		from journals_detail d 
		join journals j on j.shopid = d.shopid and j.docno = d.docno
		join chartofaccounts a on a.accountcode = d.accountcode
		where j.docdate between @startdate and @enddate and d.shopid = @shopid ` + accountGroupQuery + consolidateAccountCodeQuery + accountCodeQuery + ` 
			) as final_data order by accountcode,rowmode,docdate`

	rawDocList := []models.LedgerAccountRaw{}

	_, err := repo.pst.Raw(rawQuery, values, &rawDocList)
	if err != nil {
		return nil, err
	}

	return rawDocList, nil

}
