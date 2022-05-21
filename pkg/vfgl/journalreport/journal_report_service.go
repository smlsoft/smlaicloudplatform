package journalreport

import (
	"smlcloudplatform/pkg/models/vfgl"
	"time"
)

type IJournalReportService interface {
	ProcessTrialBalanceSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.TrialBalanceSheetReport, error)
	ProcessBalanceSheetReport(shopId string, accountGroup string, endDate time.Time) (*vfgl.BalanceSheetReport, error)
	ProcessProfitAndLossSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.ProfitAndLossSheetReport, error)
}

type JournalReportService struct {
}

func NewJournalReportService() JournalReportService {
	return JournalReportService{}
}

func (svc JournalReportService) ProcessTrialBalanceSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.TrialBalanceSheetReport, error) {
	mock := MockTrialBalanceSheetReport(shopId, accountGroup, startDate, endDate)
	return mock, nil
}

func (svc JournalReportService) ProcessBalanceSheetReport(shopId string, accountGroup string, endDate time.Time) (*vfgl.BalanceSheetReport, error) {
	mock := MockBalanceSheetReport(shopId, accountGroup, endDate)
	return mock, nil
}

func (svc JournalReportService) ProcessProfitAndLossSheetReport(shopId string, accountGroup string, startDate time.Time, endDate time.Time) (*vfgl.ProfitAndLossSheetReport, error) {
	mock := MockProfitAndLossSheetReport(shopId, accountGroup, startDate, endDate)
	return mock, nil
}
