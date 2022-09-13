package main_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatePhase(t *testing.T) {

	layout := "2006-01-02 -0700" //
	//timezone := "7"
	startDateStr := "2022-01-01 +0700"
	timeResult, err := time.Parse(layout, startDateStr)
	assert.Nil(t, err, "error not nil")
	assert.NotNil(t, timeResult, "Time is nil")

	fmt.Println(timeResult)
	fmt.Println(timeResult.UTC())

	timeResult, err = time.Parse(layout, "2022-01-01 +0000")
	assert.Nil(t, err, "error not nil")
	assert.NotNil(t, timeResult, "Time is nil")

	fmt.Println(timeResult)
	fmt.Println(timeResult.UTC())

	dateZoneLayout := "-0700"
	timewithzone, err := time.Parse(dateZoneLayout, "+0700")

	fmt.Println(timewithzone)
	fmt.Println(timewithzone.UTC())
}
