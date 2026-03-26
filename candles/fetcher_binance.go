package candles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/KhavrTrading/flowex/models"
)

// FetchBinanceCandles fetches historical klines from Binance Futures REST API.
// interval: "1m", "5m", "15m", "1h", etc.
// limit: max 1500 per request.
func FetchBinanceCandles(symbol, interval string, limit int) ([]models.CandleHLCV, error) {
	url := fmt.Sprintf(
		"https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=%s&limit=%d",
		symbol, interval, limit,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("binance klines request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("binance klines status %d: %s", resp.StatusCode, string(body))
	}

	var raw [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("binance klines decode: %w", err)
	}

	candles := make([]models.CandleHLCV, 0, len(raw))
	for _, row := range raw {
		if len(row) < 6 {
			continue
		}
		ts := int64(row[0].(float64))
		open := row[1].(string)
		high := row[2].(string)
		low := row[3].(string)
		cl := row[4].(string)
		vol := row[5].(string)

		c, err := models.NewCandleHLCVFromStrings(ts, open, high, low, cl, vol)
		if err != nil {
			continue
		}
		candles = append(candles, c)
	}

	return candles, nil
}

// FetchBinanceCandleHLC fetches historical klines as CandleHLC.
func FetchBinanceCandleHLC(symbol, interval string, limit int) ([]models.CandleHLC, error) {
	url := fmt.Sprintf(
		"https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=%s&limit=%d",
		symbol, interval, limit,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("binance klines request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("binance klines status %d: %s", resp.StatusCode, string(body))
	}

	var raw [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("binance klines decode: %w", err)
	}

	candles := make([]models.CandleHLC, 0, len(raw))
	for _, row := range raw {
		if len(row) < 6 {
			continue
		}
		slice := make([]string, 6)
		slice[0] = strconv.FormatInt(int64(row[0].(float64)), 10)
		slice[1] = row[1].(string) // open (unused but keeps index)
		slice[2] = row[2].(string) // high
		slice[3] = row[3].(string) // low
		slice[4] = row[4].(string) // close
		slice[5] = row[5].(string) // volume (unused)

		c, err := models.NewCandleHLCFromSlice(slice)
		if err != nil {
			continue
		}
		candles = append(candles, c)
	}

	return candles, nil
}
