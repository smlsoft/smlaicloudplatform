package journalreport

import (
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/vfgl/journalreport/models"
	"time"
)

type JournalReportHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc IJournalReportService
}

func NewJournalReportHttp(ms *microservice.Microservice, cfg microservice.IConfig) JournalReportHttp {

	pstConfig := microservice.NewPersisterConfig()
	pst := microservice.NewPersister(pstConfig)
	repo := NewJournalReportRepository(pst)
	jouralReportService := NewJournalReportService(repo)
	return JournalReportHttp{
		ms:  ms,
		cfg: cfg,
		svc: jouralReportService,
	}
}

func (h JournalReportHttp) RouteSetup() {
	h.ms.GET("/gl/report/trialbalancesheet", h.ProcessReportTrialBalanceSheet)
	h.ms.GET("/gl/report/balancesheet", h.ProcessBalanceSheetReport)
	h.ms.GET("/gl/report/profitandloss", h.ProcessProfitAndLossReport)
	h.ms.GET("/gl/report/ledgeraccount", h.ProcessReportLedgerAccount)
}

// Fetch TrialBalanceSheet godoc
// @Summary		ดึงข้อมูลงบทดลอง
// @Description ดึงข้อมูลงบทดลอง
// @Tags		GLReport
// @Param		accountgroup query string true "กลุ่มบัญชี"
// @Param		startdate query string true "จากวันที่ (Date Format: YYYY-MM-DD)"
// @Param		enddate query string true "ถึงวันที่ (Date Format: YYYY-MM-DD)"
// @Param		ica query int true "รวมรายการปิดปัญชี"
// @Param		timezone query string false "TimeZone"
// @Accept		json
// @Success		200 {object} models.TrialBalanceSheetReportResponse
// @Failure		400 {object} common.AuthResponseFailed
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/gl/report/trialbalancesheet [get]
func (r JournalReportHttp) ProcessReportTrialBalanceSheet(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02 -0700" //
	accountGroup := ctx.QueryParam("accountgroup")
	startDateStr := ctx.QueryParam("startdate")
	endDateStr := ctx.QueryParam("enddate")
	timeZone := ctx.QueryParam("timezone") // +07

	if timeZone != "" {
		// phase with timezone
		startDateStr = startDateStr + " " + timeZone + "00"
		endDateStr = endDateStr + " " + timeZone + "00"
	} else {
		startDateStr = startDateStr + " +0000"
		endDateStr = endDateStr + " +0000"
	}

	includeCloseAccountMode := ctx.QueryParam("ica") == "1"

	//lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(accountGroup) < 1 || len(startDateStr) < 1 || len(endDateStr) < 1 {
		ctx.ResponseError(400, "Invalid Payload.")
		return nil
	}

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		ctx.ResponseError(400, "startdate format invalid.")
		return err
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		ctx.ResponseError(400, "enddate format invalid.")
		return err
	}
	endDate = endDate.AddDate(0, 0, 1).Add(time.Second * -1)

	r.ms.Logger.Debugf("Start Process TrialBalanceSheet %v:%v, includecloseaccount: %v", startDate, endDate, includeCloseAccountMode)
	reportData, err := r.svc.ProcessTrialBalanceSheetReport(shopID, accountGroup, includeCloseAccountMode, startDate.UTC(), endDate.UTC())
	if err != nil {
		ctx.ResponseError(500, fmt.Sprintf("Failed on Process Report : %v.", err.Error()))
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.TrialBalanceSheetReportResponse{
			Success: true,
			Data:    *reportData,
		})
	return nil
}

// Fetch Data Balance Sheet godoc
// @Summary		ดึงข้อมูลงบดุล
// @Description ดึงข้อมูลงบดุล
// @Tags		GLReport
// @Param		accountgroup query string true "กลุ่มบัญชี"
// @Param		enddate query string true "ณ วันที่ (Date Format: YYYY-MM-DD)"
// @Param		ica query int true "รวมรายการปิดปัญชี"
// @Param		timezone query string false "TimeZone"
// @Accept		json
// @Success		200 {object} models.BalanceSheetReportResponse
// @Failure		400 {object} common.AuthResponseFailed
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/gl/report/balancesheet [get]
func (r JournalReportHttp) ProcessBalanceSheetReport(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02 -0700" //
	endDateStr := ctx.QueryParam("enddate")
	accountGroup := ctx.QueryParam("accountgroup")
	includeCloseAccountMode := ctx.QueryParam("ica") == "1"

	timeZone := ctx.QueryParam("timezone") // +07

	if timeZone != "" {
		// phase with timezone
		endDateStr = endDateStr + " " + timeZone + "00"
	} else {
		endDateStr = endDateStr + " +0000"
	}

	//lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(accountGroup) < 1 || len(endDateStr) < 1 {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return nil
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		ctx.ResponseError(400, "enddate format invalid.")
		return err
	}
	endDate = endDate.AddDate(0, 0, 1).Add(time.Second * -1)

	r.ms.Logger.Debugf("Start Process BalanceSheet at %v", endDate)
	reportData, err := r.svc.ProcessBalanceSheetReport(shopID, accountGroup, includeCloseAccountMode, endDate.UTC())
	if err != nil {
		ctx.ResponseError(500, fmt.Sprintf("Failed on Process Report : %v.", err.Error()))
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success: true,
			Data:    reportData,
		})
	return nil
}

// Fetch ProfitAndLossSheet godoc
// @Summary		ดึงข้อมูลงบกำไรขาดทุน
// @Description ดึงข้อมูลงบกำไรขาดทุน
// @Tags		GLReport
// @Param		accountgroup query string true "กลุ่มบัญชี"
// @Param		startdate query string true "จากวันที่ (Date Format: YYYY-MM-DD)"
// @Param		enddate query string true "ถึงวันที่ (Date Format: YYYY-MM-DD)"
// @Param		ica query int true "รวมรายการปิดปัญชี"
// @Param		timezone query string false "TimeZone"
// @Accept		json
// @Success		200 {object} models.LostAndProfitSheetReportResponse
// @Failure		400 {object} common.AuthResponseFailed
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/gl/report/profitandloss [get]
func (r JournalReportHttp) ProcessProfitAndLossReport(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02 -0700" //
	accountGroup := ctx.QueryParam("accountgroup")
	startDateStr := ctx.QueryParam("startdate")
	endDateStr := ctx.QueryParam("enddate")
	includeCloseAccountMode := ctx.QueryParam("ica") == "1"

	timeZone := ctx.QueryParam("timezone") // +07

	if timeZone != "" {
		// phase with timezone
		startDateStr = startDateStr + " " + timeZone + "00"
		endDateStr = endDateStr + " " + timeZone + "00"
	} else {
		startDateStr = startDateStr + " +0000"
		endDateStr = endDateStr + " +0000"
	}

	//lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(accountGroup) < 1 || len(startDateStr) < 1 || len(endDateStr) < 1 {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return nil
	}

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		ctx.ResponseError(400, "startdate format invalid.")
		return err
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		ctx.ResponseError(400, "enddate format invalid.")
		return err
	}
	endDate = endDate.AddDate(0, 0, 1).Add(time.Second * -1)

	r.ms.Logger.Debugf("Start Process ProfitAndLoss %v:%v", startDate, endDate)
	reportData, err := r.svc.ProcessProfitAndLossSheetReport(shopID, accountGroup, includeCloseAccountMode, startDate.UTC(), endDate.UTC())
	if err != nil {
		ctx.ResponseError(500, fmt.Sprintf("Failed on Process Report : %v.", err.Error()))
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success: true,
			Data:    reportData,
		})
	return nil
}

// Fetch Ledger Account godoc
// @Summary		บัญชีแยกประเภท
// @Description บัญชีแยกประเภท
// @Tags		GLReport
// @Param		startdate query string true "จากวันที่ (Date Format: YYYY-MM-DD)"
// @Param		enddate query string true "ถึงวันที่ (Date Format: YYYY-MM-DD)"
// @Param		timezone query string false "TimeZone"
// @Accept		json
// @Success		200 {object} models.TrialBalanceSheetReportResponse
// @Failure		400 {object} common.AuthResponseFailed
// @Failure		401 {object} common.AuthResponseFailed
// @Security	AccessToken
// @Router		/gl/report/ledgeraccount [get]
func (r JournalReportHttp) ProcessReportLedgerAccount(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02 -0700" //
	startDateStr := ctx.QueryParam("startdate")
	endDateStr := ctx.QueryParam("enddate")
	timeZone := ctx.QueryParam("timezone") // +07

	if timeZone != "" {
		// phase with timezone
		startDateStr = startDateStr + " " + timeZone + "00"
		endDateStr = endDateStr + " " + timeZone + "00"
	} else {
		startDateStr = startDateStr + " +0000"
		endDateStr = endDateStr + " +0000"
	}

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		ctx.ResponseError(400, "startdate format invalid.")
		return err
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		ctx.ResponseError(400, "enddate format invalid.")
		return err
	}
	endDate = endDate.AddDate(0, 0, 1).Add(time.Second * -1)

	r.ms.Logger.Debugf("Start Process Ledger Account %v:%v", startDate, endDate)
	reportData, err := r.svc.ProcessLedgerAccount(shopID, startDate.UTC(), endDate.UTC())
	if err != nil {
		ctx.ResponseError(500, fmt.Sprintf("Failed on Process Report : %v.", err.Error()))
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success: true,
			Data:    reportData,
		})
	return nil
}
