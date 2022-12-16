package services_test

import (
	"fmt"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/models"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/services"
	"testing"
	"time"

	"github.com/tj/assert"
)

func TestInDateSpan(t *testing.T) {

	dateFormat := "2006-01-02"

	fromDate, _ := time.Parse(dateFormat, "2020-04-01")
	toDate, _ := time.Parse(dateFormat, "2020-04-30")

	// checkFromDate, _ := time.Parse(dateFormat, "2020-03-01")
	// checkToDate, _ := time.Parse(dateFormat, "2020-04-01")

	gives := []struct {
		fromDate string
		toDate   string
		want     bool
	}{
		{"2020-03-01", "2020-03-31", false},
		{"2020-03-01", "2020-04-01", true},
		{"2020-03-31", "2020-04-01", true},
		{"2020-03-31", "2020-04-05", true},
		{"2020-03-10", "2020-04-20", true},
		{"2020-04-30", "2020-05-01", true},
		{"2020-05-01", "2020-05-10", false},
	}

	fmt.Println("Compare date 2020-04-01 2020-04-30")
	for _, give := range gives {
		checkFromDate, _ := time.Parse(dateFormat, give.fromDate)
		checkToDate, _ := time.Parse(dateFormat, give.toDate)
		result := services.InDateTimeSpan(fromDate, toDate, checkFromDate, checkToDate)
		fmt.Printf("%s %s %t %t\n", give.fromDate, give.toDate, give.want, result)
		assert.Equal(t, give.want, result)
	}

}

func TestValidatePeriod(t *testing.T) {

	dateFormat := "2006-01-02"
	ts1, _ := time.Parse(dateFormat, "2020-01-01")
	te1, _ := time.Parse(dateFormat, "2020-01-31")

	ts2, _ := time.Parse(dateFormat, "2020-03-01")
	te2, _ := time.Parse(dateFormat, "2020-03-31")

	ts3, _ := time.Parse(dateFormat, "2020-03-01")
	te3, _ := time.Parse(dateFormat, "2020-03-31")

	docList := []models.AccountPeriodMaster{
		{
			Period:    1,
			StartDate: ts1,
			EndDate:   te1,
		},
		{
			Period:    2,
			StartDate: ts2,
			EndDate:   te2,
		},
		{
			Period:    3,
			StartDate: ts3,
			EndDate:   te3,
		},
	}

	err := services.ValidatePeriod(docList)

	if err != nil {
		fmt.Println(err)
	}
}
