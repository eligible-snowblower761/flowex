package bitget

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/KhavrTrading/flowex/ws"

	log "github.com/sirupsen/logrus"
)

// NewClient creates a Bitget WebSocket client for one symbol.
func NewClient(symbol string) (*ws.BaseClient, error) {
	cfg := ws.DefaultClientConfig("Bitget", "wss://ws.bitget.com/v2/ws/public")

	// Bitget needs a 20-second ping (plain "ping" string, not JSON)
	cfg.PingInterval = 20 * time.Second
	cfg.PingMessage = func() ([]byte, error) {
		return []byte("ping"), nil
	}

	client := ws.NewBaseClient(symbol, cfg)
	client.SetDispatch(makeDispatcher(client, symbol))

	if err := client.Connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// SubscribeStream subscribes to a Bitget stream.
func SubscribeStream(client *ws.BaseClient, instType, channel, instId string) error {
	req := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"instType": instType,
				"channel":  channel,
				"instId":   instId,
			},
		},
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return err
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
		raw := string(msg)

		// Bitget sends "pong" as plain text
		if raw == "pong" {
			return
		}

		var rawMsg map[string]interface{}
		if err := json.Unmarshal(msg, &rawMsg); err != nil {
			return
		}

		// Get channel from arg.channel
		argRaw, ok := rawMsg["arg"]
		if !ok {
			return
		}
		argMap, ok := argRaw.(map[string]interface{})
		if !ok {
			return
		}
		channel, _ := argMap["channel"].(string)

		switch {
		case strings.HasPrefix(channel, "candle"):
			if cb := candleCallbacks[symbol]; cb != nil {
				var push CandlePushData
				if err := json.Unmarshal(msg, &push); err == nil {
					cb(push)
				}
			}
		case strings.HasPrefix(channel, "books"):
			if cb := depthCallbacks[symbol]; cb != nil {
				var push DepthPushData
				if err := json.Unmarshal(msg, &push); err == nil {
					cb(push)
				}
			}
		case channel == "trade":
			if cb := tradeCallbacks[symbol]; cb != nil {
				var push TradePushData
				if err := json.Unmarshal(msg, &push); err == nil {
					cb(push)
				}
			}
		default:
			log.Debugf("[Bitget:%s] unknown channel: %s", symbol, channel)
		}
	}
}

// ToSimpleSymbol converts "BTCUSDT:USDT" to "BTCUSDT".
func ToSimpleSymbol(symbol string) string {
	if idx := strings.Index(symbol, ":"); idx > 0 {
		return symbol[:idx]
	}
	return symbol
}

// ParseTimestampMs parses a string timestamp to int64 milliseconds.
func ParseTimestampMs(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}
