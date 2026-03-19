package model

// Market represents the stock exchange
type Market string

const (
	MarketSH    Market = "SH"  // Shanghai A-share
	MarketSZ    Market = "SZ"  // Shenzhen A-share
	MarketHK    Market = "HK"  // Hong Kong
	MarketUS    Market = "US"  // US stocks
	MarketIndex Market = "IDX" // Global index (no exchange suffix)
)

// Quote holds real-time stock quote data
type Quote struct {
	Code      string  // e.g. "600519"
	Name      string  // e.g. "贵州茅台"
	Market    Market  // e.g. MarketSH
	Price     float64 // current price
	PrevClose float64 // previous close
	Open      float64 // open price
	High      float64 // day high
	Low       float64 // day low
	Volume    int64   // volume in shares
	Amount    float64 // turnover amount
	Change    float64 // price change
	ChangePct float64 // price change percent
	Currency  string  // CNY / HKD / USD
}

// Symbol returns the canonical symbol string e.g. "600519.SH"
func (q *Quote) Symbol() string {
	return q.Code + "." + string(q.Market)
}

// Holding represents a stock position in the portfolio
type Holding struct {
	Code   string  `yaml:"code"`   // canonical symbol e.g. "600519.SH", "AAPL"
	Cost   float64 `yaml:"cost"`   // average cost per share
	Shares float64 `yaml:"shares"` // number of shares held
}

// PositionView combines a quote with holding info for display
type PositionView struct {
	Quote   *Quote
	Holding *Holding // nil if not in portfolio
}

// PnL calculates profit/loss for a position
func (p *PositionView) PnL() (amount float64, pct float64) {
	if p.Holding == nil || p.Holding.Shares == 0 {
		return 0, 0
	}
	amount = (p.Quote.Price - p.Holding.Cost) * p.Holding.Shares
	if p.Holding.Cost > 0 {
		pct = (p.Quote.Price - p.Holding.Cost) / p.Holding.Cost * 100
	}
	return
}

// MarketValue returns current market value of the position
func (p *PositionView) MarketValue() float64 {
	if p.Holding == nil {
		return 0
	}
	return p.Quote.Price * p.Holding.Shares
}

// CostValue returns total cost of the position
func (p *PositionView) CostValue() float64 {
	if p.Holding == nil {
		return 0
	}
	return p.Holding.Cost * p.Holding.Shares
}
