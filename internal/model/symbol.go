package model

import (
	"fmt"
	"strings"
)

// ParseSymbol parses user input into a code and market.
// Accepts formats: "600519", "600519.SH", "000858.SZ", "00700.HK", "AAPL"
func ParseSymbol(input string) (code string, market Market, err error) {
	input = strings.TrimSpace(strings.ToUpper(input))
	if input == "" {
		return "", "", fmt.Errorf("empty symbol")
	}

	// Already has explicit suffix
	if strings.HasSuffix(input, ".SH") {
		return strings.TrimSuffix(input, ".SH"), MarketSH, nil
	}
	if strings.HasSuffix(input, ".SZ") {
		return strings.TrimSuffix(input, ".SZ"), MarketSZ, nil
	}
	if strings.HasSuffix(input, ".HK") {
		return strings.TrimSuffix(input, ".HK"), MarketHK, nil
	}

	// Pure numeric: infer SH or SZ by prefix
	if isNumeric(input) {
		if strings.HasPrefix(input, "6") || strings.HasPrefix(input, "5") || strings.HasPrefix(input, "9") {
			return input, MarketSH, nil
		}
		return input, MarketSZ, nil
	}

	// HK stocks are typically 5-digit numbers starting with 0
	if isNumeric(input) && len(input) == 5 {
		return input, MarketHK, nil
	}

	// Alphabetic: US stock ticker
	if isAlpha(input) {
		return input, MarketUS, nil
	}

	return "", "", fmt.Errorf("cannot determine market for symbol: %s", input)
}

// NormalizeCode returns canonical symbol string like "600519.SH" or "AAPL"
func NormalizeCode(input string) (string, error) {
	code, market, err := ParseSymbol(input)
	if err != nil {
		return "", err
	}
	if market == MarketUS {
		return code, nil
	}
	return code + "." + string(market), nil
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func isAlpha(s string) bool {
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '-' || c == '.') {
			return false
		}
	}
	return len(s) > 0
}
