package journalreport_test

import (
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/vfgl/journalreport"
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

	get, err := repo.GetDataTrialBalance(shopId, accGroup, startDate, endDate)
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

	get, err := repo.GetDataProfitAndLoss(shopId, accGroup, startDate, endDate)
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

	get, err := repo.GetDataBalanceSheet(shopId, accGroup, endDate)
	assert.Nil(err)
	assert.NotNil(get)
}
