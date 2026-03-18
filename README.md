# fly

A terminal CLI for real-time stock quotes and portfolio tracking. Supports A-shares (Shanghai/Shenzhen), Hong Kong stocks, and US stocks.

## Install

```bash
go install github.com/yourusername/fly-cli@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/fly-cli.git
cd fly-cli
go install .
```

## Usage

### Query quotes

```bash
# Single stock — auto-detects market from code format
fly quote 600519        # Shanghai A-share
fly quote 000858.SZ     # Shenzhen A-share
fly quote 00700.HK      # Hong Kong
fly quote AAPL TSLA     # US stocks

# Short alias
fly q 600519 00700.HK AAPL
```

### Manage portfolio

```bash
# Add or update a position
fly add 600519.SH --cost 1800 --shares 100
fly add 00700.HK  --cost 320  --shares 200
fly add AAPL      --cost 150  --shares 50

# View all positions with live P&L
fly portfolio
fly pf     # alias
fly ls     # alias

# Remove a position
fly remove 600519.SH
fly rm AAPL   # alias
fly del AAPL  # alias
```

### Example output

```
┌───────────┬──────────┬─────────┬────────┬────────┬─────────┬─────────┬─────────┬──────┬───────────┬───────────┬─────────┐
│   代码    │   名称   │ 当前价  │ 涨跌额 │ 涨跌%  │  最高   │  最低   │  成本   │ 持股 │   市值    │   盈亏    │  盈亏%  │
├───────────┼──────────┼─────────┼────────┼────────┼─────────┼─────────┼─────────┼──────┼───────────┼───────────┼─────────┤
│ 600519.SH │ 贵州茅台 │ 1474.98 │ -10.02 │ -0.67% │ 1496.50 │ 1468.00 │ 1500.00 │  100 │ 147498.00 │  -2502.00 │  -1.67% │
│ 00700.HK  │ 腾讯控股 │  546.50 │  -3.50 │ -0.64% │  550.50 │  542.50 │  500.00 │  200 │ 109300.00 │  +9300.00 │  +9.30% │
│ AAPL      │ 苹果     │  254.23 │  +1.41 │ +0.56% │  255.13 │  252.18 │  200.00 │   50 │  12711.50 │  +2711.50 │ +27.11% │
└───────────┴──────────┴─────────┴────────┴────────┴─────────┴─────────┴─────────┴──────┴───────────┴───────────┴─────────┘

汇总  总成本: 260000.00  总市值: 269509.50  总盈亏: +9509.50 (3.66%)
```

Price changes and P&L are color-coded following the Chinese convention: red = up, green = down.

## Symbol format

| Input | Market | Notes |
|-------|--------|-------|
| `600519` / `600519.SH` | Shanghai A-share | Codes starting with 6/5/9 default to SH |
| `000858` / `000858.SZ` | Shenzhen A-share | Other numeric codes default to SZ |
| `00700.HK` | Hong Kong | Explicit `.HK` suffix required |
| `AAPL`, `TSLA` | US stocks | Alphabetic tickers |

## Portfolio config

Positions are stored at `~/.fly-cli/portfolio.yaml`:

```yaml
holdings:
  - code: 600519.SH
    cost: 1800.00
    shares: 100
  - code: 00700.HK
    cost: 320.00
    shares: 200
  - code: AAPL
    cost: 150.00
    shares: 50
```

## Data sources

| Source | Markets | Format |
|--------|---------|--------|
| [Eastmoney](https://www.eastmoney.com) (primary) | A-shares, HK, US | JSON |
| [Tencent Finance](https://finance.qq.com) (fallback) | A-shares, HK, US | Text |

Both are free public APIs with no authentication required. The client automatically falls back to Tencent Finance if Eastmoney is unavailable.

## Requirements

- Go 1.21+
- No API keys needed
