package main_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatePhase(t *testing.T) {

	layout := "2006-01-02" //
	//timezone := "+0700"
	startDateStr := "2022-01-01"
	timeResult, err := time.Parse(layout, startDateStr)
	assert.Nil(t, err, "error not nil")
	assert.NotNil(t, timeResult, "Time is nil")
}
