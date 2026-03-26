package indicators

// CalculateRSI computes the Relative Strength Index for the given period.
func CalculateRSI(closes []float64, period int) float64 {
	if len(closes) < period {
		return 0
	}

	gains := 0.0
	losses := 0.0

	for i := 1; i < period; i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains += change
		} else {
			losses += -change
		}
	}

	avgGain := gains / float64(period-1)
	avgLoss := losses / float64(period-1)

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}
