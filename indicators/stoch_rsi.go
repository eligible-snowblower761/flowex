package indicators

// CalculateStochRSI computes the Stochastic RSI oscillator.
// rsiP is the RSI period, stochP is the stochastic lookback period.
func CalculateStochRSI(closes []float64, rsiP, stochP int) []float64 {
	n := len(closes)
	outLen := n - rsiP - stochP + 1
	if outLen <= 0 {
		return nil
	}

	out := make([]float64, outLen)
	for i := 0; i < outLen; i++ {
		window := closes[i : i+rsiP]
		rsi := CalculateRSI(window, rsiP)

		minR, maxR := rsi, rsi
		for j := i + 1; j < i+stochP; j++ {
			r := CalculateRSI(closes[j:j+rsiP], rsiP)
			if r < minR {
				minR = r
			}
			if r > maxR {
				maxR = r
			}
		}

		if maxR == minR {
			out[i] = 0
		} else {
			out[i] = (rsi - minR) / (maxR - minR) * 100
		}
	}
	return out
}
