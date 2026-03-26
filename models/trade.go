package models

// NormalizedTrade represents a trade from any exchange in a unified format.
type NormalizedTrade struct {
	Timestamp int64   // Unix timestamp in milliseconds
	Price     float64 // Trade price
	Size      float64 // Trade size (base currency)
	SizeUSD   float64 // Trade size in USD
	Side      string  // "buy" or "sell"
	TradeID   string  // Trade ID
	Symbol    string  // Trading pair (e.g., "BTCUSDT")
	Exchange  string  // Exchange name: "binance", "bybit", "bitget"
}
