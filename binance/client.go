package binance

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/KhavrTrading/flowex/ws"

	log "github.com/sirupsen/logrus"
)

// NewClient creates a Binance futures WebSocket client for one symbol.
func NewClient(symbol string) (*ws.BaseClient, error) {
	cfg := ws.DefaultClientConfig("Binance", "wss://fstream.binance.com/ws")
	// Binance doesn't need application-level ping (server pings at protocol level)

	client := ws.NewBaseClient(symbol, cfg)
	client.SetDispatch(makeDispatcher(client, symbol))

	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// SubscribeStream subscribes to a named Binance stream (e.g., "btcusdt@kline_1m").
func SubscribeStream(client *ws.BaseClient, streamName string, id int) error {
	req := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{streamName},
		"id":     id,
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal subscribe: %w", err)
	}
	return client.WriteMessage(payload)
}

// CandleStreamName returns the Binance stream name for candles.
// interval: "1m", "3m", "5m", "15m", "30m", "1h", "4h", "1d", etc.
func CandleStreamName(symbol, interval string) string {
	return fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
}

// DepthStreamName returns the Binance stream name for partial depth snapshots.
// levels: 5, 10, or 20. speed: "100ms", "250ms", or "500ms" (empty = default).
func DepthStreamName(symbol string, levels int, speed string) string {
	lower := strings.ToLower(symbol)
	if speed != "" {
		return fmt.Sprintf("%s@depth%d@%s", lower, levels, speed)
	}
	return fmt.Sprintf("%s@depth%d", lower, levels)
}

// DiffDepthStreamName returns the Binance stream name for incremental depth updates.
// speed: "100ms", "250ms", "500ms" (empty = default 250ms).
func DiffDepthStreamName(symbol, speed string) string {
	lower := strings.ToLower(symbol)
	if speed != "" {
		return fmt.Sprintf("%s@depth@%s", lower, speed)
	}
	return fmt.Sprintf("%s@depth", lower)
}

// AggTradeStreamName returns the stream name for aggregate trades (default).
func AggTradeStreamName(symbol string) string {
	return fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol))
}

// TradeStreamName returns the stream name for individual trades.
func TradeStreamName(symbol string) string {
	return fmt.Sprintf("%s@trade", strings.ToLower(symbol))
}

// Callbacks set by the manager
var (
	candleCallbacks = make(map[string]func(CandlePushData))
	depthCallbacks  = make(map[string]func(DepthPushData))
	tradeCallbacks  = make(map[string]func(TradePushData))
)

// SetCandleCallback registers a candle callback for a symbol.
func SetCandleCallback(symbol string, cb func(CandlePushData)) {
	candleCallbacks[symbol] = cb
}

// SetDepthCallback registers a depth callback for a symbol.
func SetDepthCallback(symbol string, cb func(DepthPushData)) {
	depthCallbacks[symbol] = cb
}

// SetTradeCallback registers a trade callback for a symbol.
func SetTradeCallback(symbol string, cb func(TradePushData)) {
	tradeCallbacks[symbol] = cb
}

func makeDispatcher(client *ws.BaseClient, symbol string) ws.DispatchFunc {
	return func(msg []byte) {
		var rawMsg map[string]interface{}
		if err := json.Unmarshal(msg, &rawMsg); err != nil {
			return
		}

		eventType, _ := rawMsg["e"].(string)

		// Detect depth snapshots (no "e" field but has "bids"/"asks")
		if eventType == "" {
			if _, hasBids := rawMsg["bids"]; hasBids {
				if _, hasAsks := rawMsg["asks"]; hasAsks {
					eventType = "depthSnapshot"
				}
			}
		}

		switch eventType {
		case "kline":
			if cb := candleCallbacks[symbol]; cb != nil {
				var push CandlePushData
				if err := json.Unmarshal(msg, &push); err == nil {
					cb(push)
				}
			}
		case "depthUpdate", "depthSnapshot":
			if cb := depthCallbacks[symbol]; cb != nil {
				var push DepthPushData
				if err := json.Unmarshal(msg, &push); err == nil {
					if push.EventTime == 0 {
						push.EventTime = time.Now().UnixMilli()
					}
					cb(push)
				}
			}
		case "aggTrade", "trade":
			if cb := tradeCallbacks[symbol]; cb != nil {
				var push TradePushData
				if err := json.Unmarshal(msg, &push); err == nil {
					cb(push)
				}
			}
		default:
			if eventType != "" {
				log.Debugf("[Binance:%s] unknown event: %s", symbol, eventType)
			}
		}
	}
}
