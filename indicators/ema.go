package indicators

// CalculateEMA computes the Exponential Moving Average for the given period.
func CalculateEMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}
	k := 2.0 / float64(period+1)
	ema := average(prices[:period])
	for i := period; i < len(prices); i++ {
		ema = prices[i]*k + ema*(1-k)
	}
	return ema
}

// CalculateEMAList computes a full EMA series for each price point.
func CalculateEMAList(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}
	ema := make([]float64, len(prices))
	multiplier := 2.0 / float64(period+1)

	var sum float64
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i]-ema[i-1])*multiplier + ema[i-1]
	}
	return ema
}

func average(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}
