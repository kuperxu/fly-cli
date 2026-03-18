package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"fly/internal/model"
)

// eastmoneyProvider fetches quotes from Eastmoney push2 API
type eastmoneyProvider struct {
	client *http.Client
}

func newEastmoneyProvider() *eastmoneyProvider {
	return &eastmoneyProvider{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *eastmoneyProvider) Name() string { return "eastmoney" }

// eastmoneySecID converts a canonical symbol to eastmoney secid format
// e.g. "600519.SH" -> "1.600519", "00700.HK" -> "116.00700", "AAPL" -> "105.AAPL"
func eastmoneySecID(symbol string) (string, error) {
	code, market, err := model.ParseSymbol(symbol)
	if err != nil {
		return "", err
	}
	switch market {
	case model.MarketSH:
		return "1." + code, nil
	case model.MarketSZ:
		return "0." + code, nil
	case model.MarketHK:
		// Pad HK codes to 5 digits
		for len(code) < 5 {
			code = "0" + code
		}
		return "116." + code, nil
	case model.MarketUS:
		// Try NASDAQ first (105), NYSE is 106 — eastmoney auto-routes by ticker
		return "105." + code, nil
	}
	return "", fmt.Errorf("unsupported market for %s", symbol)
}

type emResponse struct {
	RC   int    `json:"rc"`
	Data emData `json:"data"`
}

type emData struct {
	Price     float64 `json:"f43"`
	High      float64 `json:"f46"`
	Low       float64 `json:"f44"`
	Volume    int64   `json:"f47"`
	Amount    float64 `json:"f48"`
	Code      string  `json:"f57"`
	Name      string  `json:"f58"`
	PrevClose float64 `json:"f60"`
	Change    float64 `json:"f169"`
	ChangePct float64 `json:"f170"`
	Open      float64 `json:"f46o"` // not always present
}

func (p *eastmoneyProvider) GetQuotes(symbols []string) ([]*model.Quote, error) {
	quotes := make([]*model.Quote, 0, len(symbols))
	// Eastmoney doesn't support batching via a single request cleanly for mixed markets,
	// so we fetch each symbol individually but do it sequentially (fast enough for CLI use)
	for _, sym := range symbols {
		q, err := p.fetchOne(sym)
		if err != nil {
			return quotes, fmt.Errorf("failed to fetch %s: %w", sym, err)
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}

func (p *eastmoneyProvider) fetchOne(symbol string) (*model.Quote, error) {
	secid, err := eastmoneySecID(symbol)
	if err != nil {
		return nil, err
	}
	code, market, _ := model.ParseSymbol(symbol)

	fields := "f43,f44,f46,f47,f48,f57,f58,f60,f169,f170"
	url := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/stock/get?secid=%s&fields=%s&fltt=2&invt=2",
		secid, fields,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://finance.eastmoney.com/")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result emResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid response: %w", err)
	}
	if result.RC != 0 {
		return nil, fmt.Errorf("API error rc=%d", result.RC)
	}
	d := result.Data
	if d.Price == 0 && d.Name == "" {
		// US stocks on 105 may need 106 (NYSE) — try fallback
		if market == model.MarketUS {
			return p.fetchUSFallback(symbol, code)
		}
		return nil, fmt.Errorf("no data returned for %s", symbol)
	}

	currency := currencyForMarket(market)
	// Eastmoney returns ChangePct * 100 already as percent
	return &model.Quote{
		Code:      code,
		Name:      d.Name,
		Market:    market,
		Price:     d.Price,
		PrevClose: d.PrevClose,
		High:      d.High,
		Low:       d.Low,
		Volume:    d.Volume,
		Amount:    d.Amount,
		Change:    d.Change,
		ChangePct: d.ChangePct,
		Currency:  currency,
	}, nil
}

// fetchUSFallback tries NYSE (106) if NASDAQ (105) returned empty
func (p *eastmoneyProvider) fetchUSFallback(symbol, code string) (*model.Quote, error) {
	fields := "f43,f44,f46,f47,f48,f57,f58,f60,f169,f170"
	secid := "106." + strings.ToUpper(code)
	url := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/stock/get?secid=%s&fields=%s&fltt=2&invt=2",
		secid, fields,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://finance.eastmoney.com/")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result emResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid response: %w", err)
	}
	d := result.Data
	if d.Price == 0 {
		return nil, fmt.Errorf("no data for US stock: %s", symbol)
	}
	return &model.Quote{
		Code:      code,
		Name:      d.Name,
		Market:    model.MarketUS,
		Price:     d.Price,
		PrevClose: d.PrevClose,
		High:      d.High,
		Low:       d.Low,
		Volume:    d.Volume,
		Amount:    d.Amount,
		Change:    d.Change,
		ChangePct: d.ChangePct,
		Currency:  "USD",
	}, nil
}

func currencyForMarket(m model.Market) string {
	switch m {
	case model.MarketHK:
		return "HKD"
	case model.MarketUS:
		return "USD"
	default:
		return "CNY"
	}
}
