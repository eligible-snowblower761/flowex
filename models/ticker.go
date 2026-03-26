package models

// TickerData represents a real-time ticker quote from any exchange.
type TickerData struct {
	Symbol   string
	LastPr   float64 // Last price
	Bid      float64 // Numeric bid
	Ask      float64 // Numeric ask
	BidStr   string  // String bid
	AskStr   string  // String ask
	Price    float64 // Numeric price
	PriceStr string  // String price
}
