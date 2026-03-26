package candles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/KhavrTrading/flowex/models"
)

// FetchBitgetCandles fetches historical klines from Bitget V2 REST API.
// granularity: "1m", "5m", "15m", "1H", "4H", "1D", etc.
// limit: max 200 per request.
func FetchBitgetCandles(symbol, granularity string, limit int) ([]models.CandleHLCV, error) {
	url := fmt.Sprintf(
		"https://api.bitget.com/api/v2/mix/market/candles?productType=USDT-FUTURES&symbol=%s&granularity=%s&limit=%d",
		symbol, granularity, limit,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("bitget klines request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bitget klines status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Code string     `json:"code"`
		Msg  string     `json:"msg"`
		Data [][]string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("bitget klines decode: %w", err)
	}

	if result.Code != "00000" {
		return nil, fmt.Errorf("bitget API error: %s", result.Msg)
	}

	// Bitget returns newest first — reverse for chronological order
	raw := result.Data
	candles := make([]models.CandleHLCV, 0, len(raw))
	for i := len(raw) - 1; i >= 0; i-- {
		row := raw[i]
		if len(row) < 6 {
			continue
		}
		// Bitget format: [timestamp, open, high, low, close, volume, ...]
		c, err := models.NewCandleHLCVFromSlice(row[:6])
		if err != nil {
			continue
		}
		candles = append(candles, c)
	}

	return candles, nil
}
