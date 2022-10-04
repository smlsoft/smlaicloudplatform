package journalreport_test

import (
	"fmt"
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/mock"
	"smlcloudplatform/pkg/vfgl/journalreport"
	"smlcloudplatform/pkg/vfgl/journalreport/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDataTrialBalanceReportRepository(t *testing.T) {

	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	assert := assert.New(t)

	pstConfig := microservice.NewPersisterConfig()
	pst := microservice.NewPersister(pstConfig)

	assert.NotNil(pst)
	repo := journalreport.NewJournalReportRepository(pst)

	shopId := "27dcEdktOoaSBYFmnN6G6ett4Jb"
	accGroup := "01"
	startDate := time.Date(2022, 05, 01, 00, 00, 00, 0, time.UTC)
	endDate := time.Date(2022, 05, 31, 00, 00, 00, 0, time.UTC)

	get, err := repo.GetDataTrialBalance(shopId, accGroup, false, startDate, endDate)
	assert.Nil(err)
	assert.NotNil(get)
}

func TestGetDataProfitAndLossReportRepository(t *testing.T) {
	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	assert := assert.New(t)

	pstConfig := microservice.NewPersisterConfig()
	pst := microservice.NewPersister(pstConfig)

	assert.NotNil(pst)
	repo := journalreport.NewJournalReportRepository(pst)

	shopId := "27dcEdktOoaSBYFmnN6G6ett4Jb"
	accGroup := "01"
	startDate := time.Date(2022, 05, 01, 00, 00, 00, 0, time.UTC)
	endDate := time.Date(2022, 05, 31, 00, 00, 00, 0, time.UTC)

	get, err := repo.GetDataProfitAndLoss(shopId, accGroup, false, startDate, endDate)
	assert.Nil(err)
	assert.NotNil(get)
}

func TestGetDataBalanceSheetReportRepository(t *testing.T) {
	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}
	assert := assert.New(t)

	pstConfig := microservice.NewPersisterConfig()
	pst := microservice.NewPersister(pstConfig)

	assert.NotNil(pst)
	repo := journalreport.NewJournalReportRepository(pst)

	shopId := "27dcEdktOoaSBYFmnN6G6ett4Jb"
	accGroup := "01"
	endDate := time.Date(2022, 05, 31, 00, 00, 00, 0, time.UTC)

	get, err := repo.GetDataBalanceSheet(shopId, accGroup, false, endDate)
	assert.Nil(err)
	assert.NotNil(get)
}

func TestGetDataLedgerAccount(t *testing.T) {
	if os.Getenv("SERVERLESS") == "serverless" {
		t.Skip()
	}

	pstConfig := mock.NewPersisterPostgresqlConfig()
	pst := microservice.NewPersister(pstConfig)

	repo := journalreport.NewJournalReportRepository(pst)
	results, err := repo.GetDataLedgerAccount("27dcEdktOoaSBYFmnN6G6ett4Jb", []models.LedgerAccountCodeRange{
		{
			Start: "100000",
			End:   "150000",
		},
	}, time.Date(2022, 9, 1, 00, 00, 00, 0, time.UTC), time.Date(2022, 9, 30, 00, 00, 00, 0, time.UTC))

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotEqual(0, len(results))

	fmt.Println(results)
}
