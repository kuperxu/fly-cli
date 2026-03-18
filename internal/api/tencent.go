package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fly/internal/model"
)

// tencentProvider fetches quotes from Tencent Finance API
type tencentProvider struct {
	client *http.Client
}

func newTencentProvider() *tencentProvider {
	return &tencentProvider{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *tencentProvider) Name() string { return "tencent" }

// tencentSymbol converts canonical symbol to Tencent format
// "600519.SH" -> "sh600519", "00700.HK" -> "hk00700", "AAPL" -> "usAAPL"
func tencentSymbol(symbol string) (string, error) {
	code, market, err := model.ParseSymbol(symbol)
	if err != nil {
		return "", err
	}
	switch market {
	case model.MarketSH:
		return "sh" + code, nil
	case model.MarketSZ:
		return "sz" + code, nil
	case model.MarketHK:
		for len(code) < 5 {
			code = "0" + code
		}
		return "hk" + code, nil
	case model.MarketUS:
		return "us" + strings.ToUpper(code), nil
	}
	return "", fmt.Errorf("unsupported market: %s", symbol)
}

func (p *tencentProvider) GetQuotes(symbols []string) ([]*model.Quote, error) {
	// Build batch query
	tcSymbols := make([]string, 0, len(symbols))
	symbolMap := make(map[string]string) // tencent sym -> original sym
	for _, sym := range symbols {
		ts, err := tencentSymbol(sym)
		if err != nil {
			return nil, err
		}
		tcSymbols = append(tcSymbols, ts)
		symbolMap[ts] = sym
	}

	url := "https://qt.gtimg.cn/q=" + strings.Join(tcSymbols, ",")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://finance.qq.com/")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseTencentResponse(string(body))
}

// parseTencentResponse parses the tilde-delimited Tencent response
// Format: v_sh600519="1~贵州茅台~600519~1474.49~1485.00~...";
func parseTencentResponse(body string) ([]*model.Quote, error) {
	var quotes []*model.Quote
	lines := strings.Split(strings.TrimSpace(body), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		q, err := parseTencentLine(line)
		if err != nil {
			continue // skip malformed lines
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}

func parseTencentLine(line string) (*model.Quote, error) {
	// v_sh600519="...fields...";
	eqIdx := strings.Index(line, "=\"")
	if eqIdx < 0 {
		return nil, fmt.Errorf("unexpected format")
	}
	varName := line[:eqIdx] // e.g. "v_sh600519"
	content := line[eqIdx+2:]
	content = strings.TrimSuffix(content, "\";")
	content = strings.TrimSuffix(content, "\"")

	fields := strings.Split(content, "~")
	if len(fields) < 35 {
		return nil, fmt.Errorf("too few fields")
	}

	// Determine market from variable name prefix
	varName = strings.TrimPrefix(varName, "v_")
	market, code := detectMarketFromVar(varName)

	price := parseFloat(fields[3])
	prevClose := parseFloat(fields[4])
	open := parseFloat(fields[5])
	// For HK/US, fields 4 and 5 are swapped
	if market == model.MarketHK || market == model.MarketUS {
		prevClose = parseFloat(fields[5])
		open = parseFloat(fields[4])
	}

	high := parseFloat(fields[33])
	low := parseFloat(fields[34])
	change := parseFloat(fields[31])
	changePct := parseFloat(fields[32])
	volume := parseInt(fields[6])

	return &model.Quote{
		Code:      code,
		Name:      fields[1],
		Market:    market,
		Price:     price,
		PrevClose: prevClose,
		Open:      open,
		High:      high,
		Low:       low,
		Volume:    volume,
		Change:    change,
		ChangePct: changePct,
		Currency:  currencyForMarket(market),
	}, nil
}

func detectMarketFromVar(varName string) (model.Market, string) {
	if strings.HasPrefix(varName, "sh") {
		return model.MarketSH, strings.TrimPrefix(varName, "sh")
	}
	if strings.HasPrefix(varName, "sz") {
		return model.MarketSZ, strings.TrimPrefix(varName, "sz")
	}
	if strings.HasPrefix(varName, "hk") {
		return model.MarketHK, strings.TrimPrefix(varName, "hk")
	}
	if strings.HasPrefix(varName, "us") {
		return model.MarketUS, strings.TrimPrefix(varName, "us")
	}
	return model.MarketSH, varName
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}
