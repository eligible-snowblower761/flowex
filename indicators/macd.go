package indicators

// CalculateMACD computes the MACD line, signal line, and histogram.
// Uses standard 12/26/9 parameters.
func CalculateMACD(prices []float64) (macdLine, signalLine, histogram []float64) {
	ema12 := CalculateEMAList(prices, 12)
	ema26 := CalculateEMAList(prices, 26)

	length := len(prices)
	macd := make([]float64, length)

	for i := 0; i < length; i++ {
		if i < len(ema12) && i < len(ema26) {
			macd[i] = ema12[i] - ema26[i]
		}
	}

	signal := CalculateEMAList(macd, 9)

	minLen := len(macd)
	if len(signal) < minLen {
		minLen = len(signal)
	}
	hist := make([]float64, minLen)
	for i := 0; i < minLen; i++ {
		hist[i] = macd[i] - signal[i]
	}

	return macd, signal, hist
}
