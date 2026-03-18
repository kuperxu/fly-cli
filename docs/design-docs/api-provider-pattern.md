# API Provider Pattern

## Problem

Stock data APIs are unreliable. Eastmoney may be slow or return incomplete data for some symbols
(especially US stocks on NYSE vs NASDAQ). A single provider would mean silent failures.

## Design

```
api.Provider (interface)
    Name() string
    GetQuotes(symbols []string) ([]*model.Quote, error)

api.Client
    primary  Provider   // Eastmoney
    fallback Provider   // Tencent
```

`Client.GetQuotes()` logic:
1. Call primary provider.
2. If error OR returned count < requested count → call fallback for missing symbols.
3. Merge results.

## Provider: Eastmoney (primary)

- Endpoint: `push2.eastmoney.com/api/qt/stock/get`
- Format: JSON
- Symbol format (`secid`):
  - Shanghai: `1.600519`
  - Shenzhen: `0.000858`
  - HK: `116.00700` (zero-padded to 5 digits)
  - NASDAQ: `105.AAPL`
  - NYSE: `106.BRK-A` (auto-retried if `105.` returns no data)

## Provider: Tencent (fallback)

- Endpoint: `qt.gtimg.cn/q=<symbols>`
- Format: tilde-delimited text, GBK-encoded
- Decoding: `golang.org/x/text/encoding/simplifiedchinese` → UTF-8
- Symbol format: `sh600519`, `sz000858`, `hk00700`, `usAAPL`

## Adding a New Provider

1. Create `internal/api/<name>.go`
2. Implement `Provider` interface: `Name()` and `GetQuotes()`
3. Wire into `api.NewClient()` in `client.go` as primary or fallback
