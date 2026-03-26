package depth

import "math"

// DepthMetrics represents computed order book metrics.
// Standardized across all exchanges.
type DepthMetrics struct {
	// Basic info
	Timestamp int64  `json:"timestamp"` // Milliseconds
	Symbol    string `json:"symbol"`

	// Spread
	BestBid   float64 `json:"best_bid"`
	BestAsk   float64 `json:"best_ask"`
	Spread    float64 `json:"spread"`     // ask - bid
	SpreadBps float64 `json:"spread_bps"` // spread/mid * 10000
	MidPrice  float64 `json:"mid_price"`  // (bid + ask) / 2

	// Liquidity depth (USD value) at different levels
	BidLiquidity5  float64 `json:"bid_liquidity_5"`
	AskLiquidity5  float64 `json:"ask_liquidity_5"`
	BidLiquidity10 float64 `json:"bid_liquidity_10"`
	AskLiquidity10 float64 `json:"ask_liquidity_10"`
	BidLiquidity20 float64 `json:"bid_liquidity_20"`
	AskLiquidity20 float64 `json:"ask_liquidity_20"`
	BidLiquidity50 float64 `json:"bid_liquidity_50"`
	AskLiquidity50 float64 `json:"ask_liquidity_50"`

	// Volumes (coin size) at different levels
	BidVolume5  float64 `json:"bid_volume_5"`
	AskVolume5  float64 `json:"ask_volume_5"`
	BidVolume10 float64 `json:"bid_volume_10"`
	AskVolume10 float64 `json:"ask_volume_10"`
	BidVolume20 float64 `json:"bid_volume_20"`
	AskVolume20 float64 `json:"ask_volume_20"`
	BidVolume50 float64 `json:"bid_volume_50"`
	AskVolume50 float64 `json:"ask_volume_50"`

	// Order book imbalance (bid_liq / ask_liq: >1 bullish, <1 bearish)
	ImbalanceRatio5  float64 `json:"imbalance_ratio_5"`
	ImbalanceRatio10 float64 `json:"imbalance_ratio_10"`
	ImbalanceRatio20 float64 `json:"imbalance_ratio_20"`
	ImbalanceRatio50 float64 `json:"imbalance_ratio_50"`

	// Imbalance delta: (bid-ask)/(bid+ask)*100, range -100..+100
	ImbalanceDelta10 float64 `json:"imbalance_delta_10"`
	ImbalanceDelta20 float64 `json:"imbalance_delta_20"`

	// Largest single orders (walls)
	LargestBidSize  float64 `json:"largest_bid_size"`
	LargestBidPrice float64 `json:"largest_bid_price"`
	LargestBidValue float64 `json:"largest_bid_value"`
	LargestAskSize  float64 `json:"largest_ask_size"`
	LargestAskPrice float64 `json:"largest_ask_price"`
	LargestAskValue float64 `json:"largest_ask_value"`

	// Slippage estimation (% for standardized USD order sizes)
	SlippageBuy100   float64 `json:"slippage_buy_100"`
	SlippageSell100  float64 `json:"slippage_sell_100"`
	SlippageBuy1K    float64 `json:"slippage_buy_1k"`
	SlippageSell1K   float64 `json:"slippage_sell_1k"`
	SlippageBuy5K    float64 `json:"slippage_buy_5k"`
	SlippageSell5K   float64 `json:"slippage_sell_5k"`
	SlippageBuy10K   float64 `json:"slippage_buy_10k"`
	SlippageSell10K  float64 `json:"slippage_sell_10k"`
	SlippageBuy50K   float64 `json:"slippage_buy_50k"`
	SlippageSell50K  float64 `json:"slippage_sell_50k"`
	SlippageBuy100K  float64 `json:"slippage_buy_100k"`
	SlippageSell100K float64 `json:"slippage_sell_100k"`
	SlippageBuy500K  float64 `json:"slippage_buy_500k"`
	SlippageSell500K float64 `json:"slippage_sell_500k"`
	SlippageBuy1M    float64 `json:"slippage_buy_1m"`
	SlippageSell1M   float64 `json:"slippage_sell_1m"`

	// Velocity metrics (rate-of-change, requires history)
	LiquidityVelocity10 float64 `json:"liquidity_velocity_10"`
	LiquidityVelocity50 float64 `json:"liquidity_velocity_50"`
	ImbalanceVelocity   float64 `json:"imbalance_velocity"`
	SpreadVelocity      float64 `json:"spread_velocity"`
	WallVelocity        float64 `json:"wall_velocity"`

	// Momentum
	BuyPressureMomentum  float64 `json:"buy_pressure_momentum"`
	SellPressureMomentum float64 `json:"sell_pressure_momentum"`
	WallBuildingBid      bool    `json:"wall_building_bid"`
	WallBuildingAsk      bool    `json:"wall_building_ask"`

	// Statistical z-scores (how unusual vs recent history)
	LiquidityZScore10 float64 `json:"liquidity_zscore_10"`
	ImbalanceZScore   float64 `json:"imbalance_zscore"`
	SpreadZScore      float64 `json:"spread_zscore"`

	// Depth quality
	BidLevelsCount             int     `json:"bid_levels_count"`
	AskLevelsCount             int     `json:"ask_levels_count"`
	AvgBidSize10               float64 `json:"avg_bid_size_10"`
	AvgAskSize10               float64 `json:"avg_ask_size_10"`
	TopBidConcentration5       float64 `json:"top_bid_concentration_5"`
	TopAskConcentration5       float64 `json:"top_ask_concentration_5"`
	SpreadNormImbalanceDelta10 float64 `json:"spread_norm_imbalance_delta_10"`
	SpreadNormImbalanceDelta20 float64 `json:"spread_norm_imbalance_delta_20"`
	SlippageGradientBuy        float64 `json:"slippage_gradient_buy"`
	SlippageGradientSell       float64 `json:"slippage_gradient_sell"`
	SlippageSkew1K             float64 `json:"slippage_skew_1k"`
	SlippageSkew10K            float64 `json:"slippage_skew_10k"`
}

// ComputeDepthMetrics computes order book metrics from raw bids/asks.
// bids and asks are slices of [price, quantity] string pairs.
func ComputeDepthMetrics(symbol string, timestampMs int64, bids, asks [][]string) DepthMetrics {
	m := DepthMetrics{
		Timestamp:      timestampMs,
		Symbol:         symbol,
		BidLevelsCount: len(bids),
		AskLevelsCount: len(asks),
	}

	bidPrices, bidSizes := parseLevels(bids)
	askPrices, askSizes := parseLevels(asks)

	// Spread
	if len(bidPrices) > 0 && len(askPrices) > 0 {
		m.BestBid = bidPrices[0]
		m.BestAsk = askPrices[0]
		m.Spread = m.BestAsk - m.BestBid
		m.MidPrice = (m.BestBid + m.BestAsk) / 2
		if m.MidPrice > 0 {
			m.SpreadBps = (m.Spread / m.MidPrice) * 10000
		}
	}

	// Liquidity & volume at different depths
	computeLiquidity(bidPrices, bidSizes, 5, &m.BidLiquidity5, &m.BidVolume5)
	computeLiquidity(askPrices, askSizes, 5, &m.AskLiquidity5, &m.AskVolume5)
	computeLiquidity(bidPrices, bidSizes, 10, &m.BidLiquidity10, &m.BidVolume10)
	computeLiquidity(askPrices, askSizes, 10, &m.AskLiquidity10, &m.AskVolume10)
	computeLiquidity(bidPrices, bidSizes, 20, &m.BidLiquidity20, &m.BidVolume20)
	computeLiquidity(askPrices, askSizes, 20, &m.AskLiquidity20, &m.AskVolume20)
	computeLiquidity(bidPrices, bidSizes, 50, &m.BidLiquidity50, &m.BidVolume50)
	computeLiquidity(askPrices, askSizes, 50, &m.AskLiquidity50, &m.AskVolume50)

	// Imbalance ratios
	m.ImbalanceRatio5 = safeDiv(m.BidLiquidity5, m.AskLiquidity5)
	m.ImbalanceRatio10 = safeDiv(m.BidLiquidity10, m.AskLiquidity10)
	m.ImbalanceRatio20 = safeDiv(m.BidLiquidity20, m.AskLiquidity20)
	m.ImbalanceRatio50 = safeDiv(m.BidLiquidity50, m.AskLiquidity50)

	// Imbalance deltas
	m.ImbalanceDelta10 = imbalanceDelta(m.BidLiquidity10, m.AskLiquidity10)
	m.ImbalanceDelta20 = imbalanceDelta(m.BidLiquidity20, m.AskLiquidity20)

	// Largest orders (walls)
	m.LargestBidSize, m.LargestBidPrice, m.LargestBidValue = findLargest(bidPrices, bidSizes)
	m.LargestAskSize, m.LargestAskPrice, m.LargestAskValue = findLargest(askPrices, askSizes)

	// Slippage estimation
	computeSlippage(askPrices, askSizes, m.MidPrice, true, &m)
	computeSlippage(bidPrices, bidSizes, m.MidPrice, false, &m)

	// Avg sizes, concentration
	if len(bidSizes) >= 10 {
		m.AvgBidSize10 = avgSlice(bidSizes[:10])
	}
	if len(askSizes) >= 10 {
		m.AvgAskSize10 = avgSlice(askSizes[:10])
	}
	if m.BidLiquidity50 > 0 {
		m.TopBidConcentration5 = m.BidLiquidity5 / m.BidLiquidity50 * 100
	}
	if m.AskLiquidity50 > 0 {
		m.TopAskConcentration5 = m.AskLiquidity5 / m.AskLiquidity50 * 100
	}

	// Spread-normalized imbalance
	if m.SpreadBps > 0 {
		m.SpreadNormImbalanceDelta10 = m.ImbalanceDelta10 / m.SpreadBps
		m.SpreadNormImbalanceDelta20 = m.ImbalanceDelta20 / m.SpreadBps
	}

	// Slippage gradient and skew
	if m.SlippageBuy1K > 0 && m.SlippageBuy100K > 0 {
		m.SlippageGradientBuy = (m.SlippageBuy100K - m.SlippageBuy1K) / 2 // per decade
	}
	if m.SlippageSell1K > 0 && m.SlippageSell100K > 0 {
		m.SlippageGradientSell = (m.SlippageSell100K - m.SlippageSell1K) / 2
	}
	m.SlippageSkew1K = m.SlippageBuy1K - m.SlippageSell1K
	m.SlippageSkew10K = m.SlippageBuy10K - m.SlippageSell10K

	return m
}

// parseLevels converts [[price, qty], ...] string pairs to float slices.
func parseLevels(levels [][]string) (prices, sizes []float64) {
	prices = make([]float64, 0, len(levels))
	sizes = make([]float64, 0, len(levels))
	for _, lv := range levels {
		if len(lv) < 2 {
			continue
		}
		p := parseFloatFast(lv[0])
		s := parseFloatFast(lv[1])
		prices = append(prices, p)
		sizes = append(sizes, s)
	}
	return
}

func computeLiquidity(prices, sizes []float64, n int, liq, vol *float64) {
	limit := n
	if limit > len(prices) {
		limit = len(prices)
	}
	var totalLiq, totalVol float64
	for i := 0; i < limit; i++ {
		totalLiq += prices[i] * sizes[i]
		totalVol += sizes[i]
	}
	*liq = totalLiq
	*vol = totalVol
}

func findLargest(prices, sizes []float64) (size, price, value float64) {
	for i := range prices {
		v := prices[i] * sizes[i]
		if v > value {
			size = sizes[i]
			price = prices[i]
			value = v
		}
	}
	return
}

func computeSlippage(prices, sizes []float64, mid float64, isBuy bool, m *DepthMetrics) {
	if mid <= 0 || len(prices) == 0 {
		return
	}

	targets := []struct {
		usd  float64
		buy  *float64
		sell *float64
	}{
		{100, &m.SlippageBuy100, &m.SlippageSell100},
		{1_000, &m.SlippageBuy1K, &m.SlippageSell1K},
		{5_000, &m.SlippageBuy5K, &m.SlippageSell5K},
		{10_000, &m.SlippageBuy10K, &m.SlippageSell10K},
		{50_000, &m.SlippageBuy50K, &m.SlippageSell50K},
		{100_000, &m.SlippageBuy100K, &m.SlippageSell100K},
		{500_000, &m.SlippageBuy500K, &m.SlippageSell500K},
		{1_000_000, &m.SlippageBuy1M, &m.SlippageSell1M},
	}

	for _, t := range targets {
		slip := estimateSlippage(prices, sizes, mid, t.usd)
		if isBuy {
			*t.buy = slip
		} else {
			*t.sell = slip
		}
	}
}

func estimateSlippage(prices, sizes []float64, mid, targetUSD float64) float64 {
	var filled float64
	var cost float64
	for i := range prices {
		levelUSD := prices[i] * sizes[i]
		remaining := targetUSD - filled
		if remaining <= 0 {
			break
		}
		if levelUSD >= remaining {
			frac := remaining / levelUSD
			cost += prices[i] * sizes[i] * frac
			filled += remaining
		} else {
			cost += levelUSD
			filled += levelUSD
		}
	}
	if filled <= 0 {
		return 0
	}
	avgPrice := cost / (filled / mid) // weighted average fill price approx
	_ = avgPrice
	// Simplified: slippage = (vwap - mid) / mid * 100
	vwap := cost / filled * mid
	return math.Abs(vwap-mid) / mid * 100
}

func safeDiv(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func imbalanceDelta(bid, ask float64) float64 {
	total := bid + ask
	if total == 0 {
		return 0
	}
	return (bid - ask) / total * 100
}

func avgSlice(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	var sum float64
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s))
}

func parseFloatFast(s string) float64 {
	v := 0.0
	neg := false
	i := 0
	if i < len(s) && s[i] == '-' {
		neg = true
		i++
	}
	for ; i < len(s) && s[i] != '.'; i++ {
		v = v*10 + float64(s[i]-'0')
	}
	if i < len(s) && s[i] == '.' {
		i++
		frac := 0.1
		for ; i < len(s); i++ {
			v += float64(s[i]-'0') * frac
			frac *= 0.1
		}
	}
	if neg {
		return -v
	}
	return v
}
