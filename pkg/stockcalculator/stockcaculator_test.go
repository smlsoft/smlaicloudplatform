package stockcalculator_test

import (
	"smlaicloudplatform/pkg/stockcalculator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStockCalculator(t *testing.T) {

	calculator := stockcalculator.NewStockCalculator("TESTSHOP", "TESTBARCODE", 2, 0, 0)

	assert.Equal(t, 0.0, calculator.BalanceAmount())
	assert.Equal(t, 0.0, calculator.BalanceQty())
	assert.Equal(t, 0.0, calculator.AverageCost())

	calculator.ApplyStock(100, 1000)
	assert.Equal(t, 1000.0, calculator.BalanceAmount())
	assert.Equal(t, 100.0, calculator.BalanceQty())
	assert.Equal(t, 10.0, calculator.AverageCost())

	calculator.ApplyStock(2, 22)
	assert.Equal(t, 1022.0, calculator.BalanceAmount())
	assert.Equal(t, 102.0, calculator.BalanceQty())
	assert.Equal(t, 10.02, calculator.AverageCost())

}

func TestStockCalculator2(t *testing.T) {

	calculator := stockcalculator.NewStockCalculator("TESTSHOP", "TESTBARCODE", 2, 0, 0)

	assert.Equal(t, 0.0, calculator.BalanceAmount())
	assert.Equal(t, 0.0, calculator.BalanceQty())
	assert.Equal(t, 0.0, calculator.AverageCost())

	calculator.ApplyStock(100, 1000)
	assert.Equal(t, 1000.0, calculator.BalanceAmount())
	assert.Equal(t, 100.0, calculator.BalanceQty())
	assert.Equal(t, 10.0, calculator.AverageCost())

	calculator.ApplyStock(2, 22)
	assert.Equal(t, 1022.0, calculator.BalanceAmount())
	assert.Equal(t, 102.0, calculator.BalanceQty())
	assert.Equal(t, 10.02, calculator.AverageCost())

	calculator.ApplyStock(3, 30)
	assert.Equal(t, 1022.0, calculator.BalanceAmount())
	assert.Equal(t, 102.0, calculator.BalanceQty())
	assert.Equal(t, 10.02, calculator.AverageCost())

}
