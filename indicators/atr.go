package indicators

import (
	"fmt"
	"math"

	"github.com/KhavrTrading/flowex/models"
)

// CalculateATR computes the Average True Range from CandleHLC data.
func CalculateATR(candles []models.CandleHLC, period int) float64 {
	if len(candles) <= period {
		return 0
	}

	var trSum float64
	startIdx := len(candles) - period

	for i := startIdx; i < len(candles); i++ {
		highLow := candles[i].High - candles[i].Low
		highPrevClose := math.Abs(candles[i].High - candles[i-1].GetClose())
		lowPrevClose := math.Abs(candles[i].Low - candles[i-1].GetClose())
		trueRange := math.Max(highLow, math.Max(highPrevClose, lowPrevClose))
		trSum += trueRange
	}

	return trSum / float64(period)
}

// EvaluateATR computes ATR and determines if it's rising.
func EvaluateATR(candles []models.CandleHLC, period int, atrThresholdPercent float64) (atr, atrThreshold float64, atrRising bool, err error) {
	if len(candles) <= period+1 {
		return 0, 0, false, fmt.Errorf("not enough candles")
	}

	latestPrice := candles[len(candles)-1].GetClose()
	atr = CalculateATR(candles[len(candles)-period-1:], period)
	prevATR := CalculateATR(candles[len(candles)-period-2:len(candles)-1], period)

	atrRising = atr > prevATR
	atrThreshold = latestPrice * atrThresholdPercent

	return atr, atrThreshold, atrRising, nil
}
