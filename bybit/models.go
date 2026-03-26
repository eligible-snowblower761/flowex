package bybit

import "fmt"

// CandlePushData represents Bybit V5 kline stream data.
type CandlePushData struct {
	Topic string      `json:"topic"` // "kline.1.BTCUSDT"
	Type  string      `json:"type"`  // "snapshot" or "delta"
	Data  []KlineData `json:"data"`
	TS    int64       `json:"ts"`
}

// KlineData represents kline data within a push.
type KlineData struct {
	Start     int64  `json:"start"`     // Start time (ms)
	End       int64  `json:"end"`       // End time (ms)
	Interval  string `json:"interval"`  // "1" for 1 minute
	Open      string `json:"open"`
	Close     string `json:"close"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Volume    string `json:"volume"`
	Turnover  string `json:"turnover"`
	Confirm   bool   `json:"confirm"`   // Whether candle is closed
	Timestamp int64  `json:"timestamp"` // Candle timestamp (ms)
}

// ToSlice returns [timestamp, open, high, low, close, volume].
func (k *KlineData) ToSlice() []string {
	return []string{
		fmt.Sprintf("%d", k.Start),
		k.Open, k.High, k.Low, k.Close, k.Volume,
	}
}

// DepthPushData represents Bybit V5 order book stream data.
type DepthPushData struct {
	Topic string    `json:"topic"` // "orderbook.50.BTCUSDT"
	Type  string    `json:"type"`  // "snapshot" or "delta"
	TS    int64     `json:"ts"`
	Data  DepthData `json:"data"`
}

// DepthData represents order book data.
type DepthData struct {
	Symbol   string     `json:"s"`
	Bids     [][]string `json:"b"` // [[price, size],...]
	Asks     [][]string `json:"a"`
	UpdateID int64      `json:"u"`
	SeqNum   int64      `json:"seq"`
}

// TradePushData represents Bybit V5 public trade stream data.
type TradePushData struct {
	Topic string      `json:"topic"` // "publicTrade.BTCUSDT"
	Type  string      `json:"type"`
	TS    int64       `json:"ts"`
	Data  []TradeData `json:"data"`
}

// TradeData represents individual trade data.
type TradeData struct {
	Timestamp  int64  `json:"T"`  // Trade timestamp (ms)
	Symbol     string `json:"s"`
	Side       string `json:"S"`  // "Buy" or "Sell"
	Size       string `json:"v"`  // Trade size
	Price      string `json:"p"`
	Direction  string `json:"L"`  // PlusTick, ZeroPlusTick, MinusTick, ZeroMinusTick
	TradeID    string `json:"i"`
	BlockTrade bool   `json:"BT"`
}
