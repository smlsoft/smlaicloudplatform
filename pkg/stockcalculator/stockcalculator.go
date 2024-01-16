package stockcalculator

import (
	"math"
	"smlcloudplatform/pkg/round"
)

type IStockCalculator interface {
	ApplyStock(qty float64, cost float64) (sumOfCost float64, averageCost float64)
	ReduceStock(qty float64) (sumOfCost float64, averageCost float64)
	ReduceStockWithCost(qty float64, cost float64) (sumOfCost float64, averageCost float64)
	BalanceAmount() float64
	BalanceQty() float64
	AverageCost() float64
}

type StockCalculator struct {
	ShopID        string
	Barcode       string
	AmountDigit   int8
	balanceQty    float64
	balanceAmount float64
	averageCost   float64
}

func NewStockCalculator(shopID string, barcode string, amountDigit int8, balanceQtyFirst float64, balanceAmountFirst float64) IStockCalculator {

	if amountDigit <= 0 {
		amountDigit = 2
	}

	averageCost := round.Round(balanceAmountFirst/balanceQtyFirst, amountDigit)
	if math.IsNaN(averageCost) {
		averageCost = 0
	}

	return &StockCalculator{
		ShopID:        shopID,
		Barcode:       barcode,
		AmountDigit:   amountDigit,
		balanceQty:    balanceQtyFirst,
		balanceAmount: balanceAmountFirst,
		averageCost:   averageCost,
	}
}

func (sc *StockCalculator) ApplyStock(qty float64, cost float64) (sumOfCost float64, averageCost float64) {
	cost = round.Round(cost, sc.AmountDigit)
	sc.balanceQty += qty
	sc.balanceAmount += cost

	average := round.Round(cost/qty, sc.AmountDigit)

	sc.averageCost = round.Round(sc.balanceAmount/sc.balanceQty, sc.AmountDigit)

	return sc.balanceQty, average
}

func (sc *StockCalculator) ReduceStock(qty float64) (sumOfCost float64, averageCost float64) {
	sc.balanceQty -= qty

	averageCostOut := sc.averageCost
	sumOfCostOut := averageCostOut * qty
	costOut := round.Round(sumOfCostOut, sc.AmountDigit)
	sc.balanceAmount -= costOut
	sc.averageCost = round.Round(sc.balanceAmount/sc.balanceQty, sc.AmountDigit)

	return costOut, averageCostOut
}

func (sc *StockCalculator) ReduceStockWithCost(qty float64, cost float64) (sumOfCost float64, averageCost float64) {

	sc.balanceQty -= qty

	averageCostOut := cost
	costOut := round.Round(averageCostOut*qty, sc.AmountDigit)
	sc.balanceAmount -= costOut
	sc.averageCost = round.Round(sc.balanceAmount/sc.balanceQty, sc.AmountDigit)

	return costOut, averageCostOut
}

func (sc *StockCalculator) BalanceAmount() float64 {
	return sc.balanceAmount
}

func (sc *StockCalculator) BalanceQty() float64 {
	return sc.balanceQty
}

func (sc *StockCalculator) AverageCost() float64 {
	return sc.averageCost
}
