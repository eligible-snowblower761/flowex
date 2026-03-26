package indicators

import "github.com/KhavrTrading/flowex/models"

// SupportResistance identifies pivot-based support and resistance levels.
// Returns the distance (%) from current price to S/R, and a break/retest score.
func SupportResistance(candles []models.CandleHLC, lookback, retWindow int) (supportPct, resistancePct, score float64) {
	n := len(candles)
	if n < lookback*2 {
		return 0, 0, 0
	}

	currentPrice := candles[n-1].GetClose()

	var supportLevel, resistanceLevel float64
	lastSupportIdx := -1
	lastResistanceIdx := -1

	for i := lookback; i < n-lookback; i++ {
		isPivotLow := true
		isPivotHigh := true

		for j := i - lookback; j <= i+lookback; j++ {
			if candles[j].Low < candles[i].Low {
				isPivotLow = false
			}
			if candles[j].High > candles[i].High {
				isPivotHigh = false
			}
			if !isPivotLow && !isPivotHigh {
				break
			}
		}

		if isPivotLow {
			lastSupportIdx = i
		}
		if isPivotHigh {
			lastResistanceIdx = i
		}
	}

	if lastSupportIdx != -1 {
		supportLevel = candles[lastSupportIdx].Low
		supportPct = ((supportLevel - currentPrice) / currentPrice) * 100.0
		if supportPct < 0 {
			supportPct = -supportPct
		}
		if supportPct > 100 {
			supportPct = 100
		}
	}

	if lastResistanceIdx != -1 {
		resistanceLevel = candles[lastResistanceIdx].High
		resistancePct = ((resistanceLevel - currentPrice) / currentPrice) * 100.0
		if resistancePct < 0 {
			resistancePct = -resistancePct
		}
		if resistancePct > 100 {
			resistancePct = 100
		}
	}

	recent := retWindow
	if recent > n-1 {
		recent = n - 1
	}

	var breaksUp, breaksDown, retestsUp, retestsDown int
	for i := n - recent; i < n; i++ {
		c := candles[i]
		if resistanceLevel > 0 && c.GetClose() > resistanceLevel {
			breaksUp++
		}
		if supportLevel > 0 && c.GetClose() < supportLevel {
			breaksDown++
		}
		if resistanceLevel > 0 && c.Low <= resistanceLevel && c.GetClose() >= resistanceLevel {
			retestsUp++
		}
		if supportLevel > 0 && c.High >= supportLevel && c.GetClose() <= supportLevel {
			retestsDown++
		}
	}

	total := breaksUp + breaksDown + retestsUp + retestsDown
	if total > 0 {
		score = (float64(retestsUp+retestsDown)*2.0 + float64(breaksUp+breaksDown)) / float64(total*2) * 100.0
	}

	return supportPct, resistancePct, score
}
