package bitget

// DepthPushData represents Bitget order book stream data.
type DepthPushData struct {
	Action string      `json:"action"` // "snapshot" or "update"
	Arg    DepthArg    `json:"arg"`
	Data   []DepthData `json:"data"`
}

// DepthArg identifies the subscription.
type DepthArg struct {
	InstType string `json:"instType"` // "SPOT" or "UMCBL"
	Channel  string `json:"channel"`  // "books" or "books5"
	InstID   string `json:"instId"`   // "BTCUSDT"
}

// DepthData represents order book data.
type DepthData struct {
	Asks [][]string `json:"asks"` // [[price, size],...]
	Bids [][]string `json:"bids"`
	TS   string     `json:"ts"` // Timestamp string
}

// TradePushData represents Bitget market trade stream data.
type TradePushData struct {
	Action string      `json:"action"` // "snapshot" or "update"
	Arg    TradeArg    `json:"arg"`
	Data   []TradeData `json:"data"`
}

// TradeArg identifies the subscription.
type TradeArg struct {
	InstType string `json:"instType"`
	Channel  string `json:"channel"` // "trade"
	InstID   string `json:"instId"`
}

// TradeData represents individual trade data.
type TradeData struct {
	InstID  string `json:"instId"`
	TradeID string `json:"tradeId"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Side    string `json:"side"` // "buy" or "sell"
	TS      string `json:"ts"`   // Timestamp
}

// CandlePushData represents Bitget candle stream data.
// Bitget sends candle data in arrays: [[ts, open, high, low, close, volume], ...]
type CandlePushData struct {
	Action string    `json:"action"`
	Arg    CandleArg `json:"arg"`
	Data   [][]string `json:"data"` // [[ts, open, high, low, close, volume, ...]]
}

// CandleArg identifies the candle subscription.
type CandleArg struct {
	InstType string `json:"instType"`
	Channel  string `json:"channel"` // "candle1m"
	InstID   string `json:"instId"`
}
