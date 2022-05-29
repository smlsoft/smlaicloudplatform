package journalreport_test

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/vfgl/journalreport"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDataFromTrialBalanceReportRepository(t *testing.T) {

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

func TestGetDataFromProfitAndLossReportRepository(t *testing.T) {
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
