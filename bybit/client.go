package bybit

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/KhavrTrading/flowex/ws"

	log "github.com/sirupsen/logrus"
)

// NewClient creates a Bybit V5 linear WebSocket client for one symbol.
func NewClient(symbol string) (*ws.BaseClient, error) {
	cfg := ws.DefaultClientConfig("Bybit", "wss://stream.bybit.com/v5/public/linear")

	// Bybit needs a 15-second ping
	cfg.PingInterval = 15 * time.Second
	cfg.PingMessage = func() ([]byte, error) {
		return json.Marshal(map[string]interface{}{"op": "ping"})
	}

	client := ws.NewBaseClient(symbol, cfg)
	client.SetDispatch(makeDispatcher(client, symbol))

	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// SubscribeStream subscribes to a named Bybit stream (e.g., "kline.1.BTCUSDT").
func SubscribeStream(client *ws.BaseClient, streamName string) error {
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []string{streamName},
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal subscribe: %w", err)
	}
	return client.WriteMessage(payload)
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

		// Pong response
		if op, ok := rawMsg["op"].(string); ok && op == "pong" {
			return
		}

		topic, _ := rawMsg["topic"].(string)

		switch {
		case strings.HasPrefix(topic, "kline."):
			if cb := candleCallbacks[symbol]; cb != nil {
				var push CandlePushData
				if err := json.Unmarshal(msg, &push); err == nil {
					cb(push)
				}
			}
		case strings.HasPrefix(topic, "orderbook."):
			if cb := depthCallbacks[symbol]; cb != nil {
				var push DepthPushData
				if err := json.Unmarshal(msg, &push); err == nil {
					cb(push)
				}
			}
		case strings.HasPrefix(topic, "publicTrade."):
			if cb := tradeCallbacks[symbol]; cb != nil {
				var push TradePushData
				if err := json.Unmarshal(msg, &push); err == nil {
					cb(push)
				}
			}
		default:
			if topic != "" {
				log.Debugf("[Bybit:%s] unknown topic: %s", symbol, topic)
			}
		}
	}
}

// ToSimpleSymbol converts "BTC/USDT:USDT" to "BTCUSDT".
func ToSimpleSymbol(symbol string) string {
	s := strings.Replace(symbol, "/", "", 1)
	if idx := strings.Index(s, ":"); idx > 0 {
		s = s[:idx]
	}
	return s
}
