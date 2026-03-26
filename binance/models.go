package binance

import "fmt"

// CandlePushData represents Binance kline/candlestick stream data.
type CandlePushData struct {
	EventType string     `json:"e"` // "kline"
	EventTime int64      `json:"E"`
	Symbol    string     `json:"s"`
	Kline     *KlineData `json:"k"`
}

// KlineData represents kline data within a push.
type KlineData struct {
	StartTime           int64  `json:"t"` // Kline start time
	EndTime             int64  `json:"T"` // Kline close time
	Symbol              string `json:"s"`
	Interval            string `json:"i"` // e.g. "1m"
	FirstTradeID        int64  `json:"f"`
	LastTradeID         int64  `json:"L"`
	Open                string `json:"o"`
	Close               string `json:"c"`
	High                string `json:"h"`
	Low                 string `json:"l"`
	Volume              string `json:"v"` // Base asset volume
	NumberOfTrades      int    `json:"n"`
	IsClosed            bool   `json:"x"`
	QuoteAssetVolume    string `json:"q"`
	TakerBuyBaseVolume  string `json:"V"`
	TakerBuyQuoteVolume string `json:"Q"`
	Ignore              string `json:"B"`
}

// ToSlice returns [timestamp, open, high, low, close, volume].
func (k *KlineData) ToSlice() []string {
	return []string{
		fmt.Sprintf("%d", k.StartTime),
		k.Open, k.High, k.Low, k.Close, k.Volume,
	}
}

// DepthPushData represents Binance order book depth stream data.
type DepthPushData struct {
	EventType     string     `json:"e"`            // "depthUpdate"
	EventTime     int64      `json:"E"`
	Symbol        string     `json:"s"`
	FirstUpdateID int64      `json:"U"`
	FinalUpdateID int64      `json:"u"`
	LastUpdateID  int64      `json:"lastUpdateId"` // For snapshot
	Bids          [][]string `json:"b"`
	Asks          [][]string `json:"a"`
}

// TradePushData represents Binance market trade stream data.
type TradePushData struct {
	EventType     string `json:"e"` // "trade" or "aggTrade"
	EventTime     int64  `json:"E"`
	Symbol        string `json:"s"`
	TradeID       int64  `json:"t"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	BuyerOrderID  int64  `json:"b"`
	SellerOrderID int64  `json:"a"`
	TradeTime     int64  `json:"T"`
	IsBuyerMaker  bool   `json:"m"`
	Ignore        bool   `json:"M"`
}
