package candles

import "github.com/KhavrTrading/flowex/models"

// Aggregate1mTo5m aggregates 1-minute candles into 5-minute candles.
// Candles are grouped by their 5-minute boundary (timestamp aligned).
func Aggregate1mTo5m(candles []models.CandleHLCV) []models.CandleHLCV {
	return aggregate(candles, 5*60*1000) // 5 minutes in ms
}

// Aggregate1mTo15m aggregates 1-minute candles into 15-minute candles.
func Aggregate1mTo15m(candles []models.CandleHLCV) []models.CandleHLCV {
	return aggregate(candles, 15*60*1000) // 15 minutes in ms
}

// Aggregate aggregates candles into bars of the given duration (in milliseconds).
func Aggregate(candles []models.CandleHLCV, durationMs int64) []models.CandleHLCV {
	return aggregate(candles, durationMs)
}

func aggregate(candles []models.CandleHLCV, durationMs int64) []models.CandleHLCV {
	if len(candles) == 0 || durationMs <= 0 {
		return nil
	}

	var result []models.CandleHLCV
	var current *models.CandleHLCV
	var currentBucket int64

	for _, c := range candles {
		bucket := (c.Ts / durationMs) * durationMs

		if current == nil || bucket != currentBucket {
			// Start a new aggregated candle
			if current != nil {
				result = append(result, *current)
			}
			cp := c // copy
			current = &cp
			current.Ts = bucket
			currentBucket = bucket
		} else {
			// Merge into current candle
			if c.High > current.High {
				current.High = c.High
			}
			if c.Low < current.Low {
				current.Low = c.Low
			}
			current.Close = c.Close
			current.Volume += c.Volume
		}
	}

	if current != nil {
		result = append(result, *current)
	}

	return result
}
