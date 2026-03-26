package ws

import (
	"time"

	"github.com/KhavrTrading/flowex/depth"
	"github.com/KhavrTrading/flowex/models"
)

// StreamType identifies a WebSocket data stream.
type StreamType string

const (
	StreamCandle StreamType = "candle"
	StreamDepth  StreamType = "depth"
	StreamTrade  StreamType = "trade"
)

// CandleHandler is called when a new candle update arrives.
type CandleHandler func(candle models.CandleHLCV)

// DepthHandler is called when a new order book snapshot arrives.
type DepthHandler func(bids, asks [][]string, timestamp int64)

// TradeHandler is called when a new trade arrives.
type TradeHandler func(trade models.NormalizedTrade)

// Snapshot is an immutable, point-in-time view of a symbol's state.
// External readers get these via atomic load — no locks needed.
type Snapshot struct {
	Timestamp  time.Time
	Candles    []models.CandleHLCV
	DepthStore *depth.Store
	Trades     []models.NormalizedTrade
}

// Manager coordinates WebSocket clients and workers for one exchange.
type Manager interface {
	SubscribeCandle(symbol string, handler CandleHandler) error
	SubscribeDepth(symbol string, handler DepthHandler) error
	SubscribeTrade(symbol string, handler TradeHandler) error
	SubscribeAll(symbol string, ch CandleHandler, dh DepthHandler, th TradeHandler) error
	Unsubscribe(symbol string, streamType StreamType) error
	UnsubscribeAll(symbol string) error
	GetSnapshot(symbol string) *Snapshot
	GetStatus() map[string]interface{}
	Shutdown()
}
