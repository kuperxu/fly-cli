# Symbol Formats

All supported input formats and their resolved markets.

## Input Parsing Rules

| Input example | Resolved market | Rule |
|--------------|----------------|------|
| `600519` | SH (Shanghai) | Numeric, starts with 6, 5, or 9 |
| `000858` | SZ (Shenzhen) | Numeric, starts with 0, 1, 2, 3, or 4 |
| `600519.SH` | SH | Explicit `.SH` suffix |
| `000858.SZ` | SZ | Explicit `.SZ` suffix |
| `00700.HK` | HK | Explicit `.HK` suffix |
| `AAPL` | US | All-alphabetic ticker |

Implemented in `internal/model/symbol.go` → `ParseSymbol()`.

## Canonical Form

After parsing, symbols are normalized to:
- A-shares: `<CODE>.<MARKET>` — e.g. `600519.SH`, `000858.SZ`
- HK stocks: `<CODE>.HK` — e.g. `00700.HK`
- US stocks: bare ticker — e.g. `AAPL` (no suffix)

Used as the storage key in `~/.fly-cli/portfolio.yaml`. Matching is case-insensitive.

## HK Code Padding

HK codes are zero-padded to 5 digits when constructing API calls:
- User input `700.HK` → normalized `00700.HK`
- Eastmoney secid: `116.00700`
- Tencent symbol: `hk00700`
