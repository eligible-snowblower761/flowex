package candles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/KhavrTrading/flowex/models"
)

// FetchBybitCandles fetches historical klines from Bybit V5 REST API.
// interval: "1" (1m), "5" (5m), "15", "60", "240", "D", "W".
// limit: max 200 per request.
func FetchBybitCandles(symbol, interval string, limit int) ([]models.CandleHLCV, error) {
	url := fmt.Sprintf(
		"https://api.bybit.com/v5/market/kline?category=linear&symbol=%s&interval=%s&limit=%d",
		symbol, interval, limit,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("bybit klines request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bybit klines status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			List [][]string `json:"list"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("bybit klines decode: %w", err)
	}

	if result.RetCode != 0 {
		return nil, fmt.Errorf("bybit API error: %s", result.RetMsg)
	}

	// Bybit returns newest first — reverse for chronological order
	raw := result.Result.List
	candles := make([]models.CandleHLCV, 0, len(raw))
	for i := len(raw) - 1; i >= 0; i-- {
		row := raw[i]
		if len(row) < 6 {
			continue
		}
		// Bybit format: [timestamp, open, high, low, close, volume, turnover]
		c, err := models.NewCandleHLCVFromSlice(row[:6])
		if err != nil {
			continue
		}
		candles = append(candles, c)
	}

	return candles, nil
}
