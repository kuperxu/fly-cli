package model

// Alert represents price alert levels for a stock
type Alert struct {
	Code  string  `yaml:"code"`            // canonical symbol e.g. "600519.SH", "AAPL"
	Entry float64 `yaml:"entry,omitempty"` // entry price (trigger when price <= entry)
	TP1   float64 `yaml:"tp1,omitempty"`   // take-profit price (trigger when price >= tp1)
	SL    float64 `yaml:"sl,omitempty"`    // stop-loss price (trigger when price <= sl)
}

// TriggeredAlert represents an alert that has been triggered by current price
type TriggeredAlert struct {
	Code   string
	Name   string
	Price  float64
	Type   string  // "Entry" / "TP1" / "SL"
	Target float64
}
