package indicators

import (
	"math"

	"github.com/KhavrTrading/flowex/models"
)

// BollingerMeanDeviation computes the Bollinger Mean Deviation score and oscillator
// standard deviation. baseLength is the SMA period, deviationLength is the lookback
// for the oscillator z-score.
func BollingerMeanDeviation(candles []models.CandleHLC, baseLength, deviationLength int) (score, oscSD float64) {
	n := len(candles)
	minRequired := baseLength + deviationLength
	if n < minRequired {
		return 0.0, 1.0
	}

	closes := make([]float64, n)
	for i := 0; i < n; i++ {
		closes[i] = candles[i].GetClose()
	}

	mSeries := make([]float64, n)
	for i := baseLength - 1; i < n; i++ {
		window := closes[i-baseLength+1 : i+1]
		sma := sma(window)
		stdDev := stdDev(window, sma)
		if stdDev <= 0 {
			mSeries[i] = 0.0
		} else {
			mSeries[i] = (closes[i] - sma) / stdDev
		}
	}

	validMValues := make([]float64, 0, deviationLength)
	for i := n - deviationLength; i < n; i++ {
		if i >= baseLength-1 {
			validMValues = append(validMValues, mSeries[i])
		}
	}

	if len(validMValues) < deviationLength {
		return 0.0, 1.0
	}

	mMean := sma(validMValues)
	oscSD = stdDev(validMValues, mMean)

	if oscSD <= 0 {
		return 0.0, 1.0
	}

	score = mSeries[n-1] / oscSD
	return score, oscSD
}

// BMD is a convenience wrapper using timeframe-dependent parameters.
// "1m" uses 40/50 periods; other timeframes use 20/25.
func BMD(candles []models.CandleHLC, timeframe string) (float64, float64) {
	if timeframe == "1m" {
		return BollingerMeanDeviation(candles, 40, 50)
	}
	return BollingerMeanDeviation(candles, 20, 25)
}

func sma(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func stdDev(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sumSq := 0.0
	for _, v := range values {
		d := v - mean
		sumSq += d * d
	}
	return math.Sqrt(sumSq / float64(len(values)))
}
