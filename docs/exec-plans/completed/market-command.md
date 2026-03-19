# Completed: `fly market` — Major Index Quotes

**Date:** 2026-03-19

## What was built

New command `fly market` (alias `fly mk`) that displays real-time quotes for five major
market indicators:

| Indicator | Eastmoney secid |
|-----------|----------------|
| 上证指数 | `1.000001` |
| 科创50 | `1.000688` |
| 恒生指数 | `100.HSI` |
| 纳斯达克 | `100.NDX` |
| 中国10年期国债 | `171.CN10Y` |

## Key decisions

1. **Batch `ulist.np` endpoint** — Used `push2.eastmoney.com/api/qt/ulist.np/get` instead of
   the per-stock `qt/stock/get` endpoint. This supports batch queries for indices, bonds, and
   stocks using a single request with comma-separated secids. Field codes differ from the
   stock endpoint (e.g. `f2`=price vs `f43`, `f14`=name vs `f58`).

2. **Added `MarketIndex` ("IDX") market type** — Indices and bonds don't belong to SH/SZ/HK/US
   markets. A dedicated market type keeps the display clean (shows `CN10Y.IDX`).

3. **`interface{}` for Volume/Amount in `emData`** — Eastmoney returns `"-"` (string) instead
   of a number for `f47`/`f48` on some foreign indices (e.g. NDX). Changed these fields from
   `int64`/`float64` to `interface{}` with `toInt64()`/`toFloat64()` helpers. Tracked as
   TD-005 in tech debt.

4. **No fallback provider** — Tencent API doesn't support index/bond queries, so
   `GetIndexQuotes()` only uses the eastmoney provider.

## Files changed

- `internal/model/stock.go` — Added `MarketIndex` constant
- `internal/api/eastmoney.go` — Added `fetchIndices()` batch method using `ulist.np`, `toInt64()`, `toFloat64()`; changed `emData` Volume/Amount to `interface{}`
- `internal/api/client.go` — Added `GetIndexQuotes()` method
- `cmd/market.go` — New command file
- `cmd/root.go` — Registered `marketCmd`

## Trade-offs

- Index list is hardcoded in `cmd/market.go`. If users want custom indices, this would need
  a config mechanism. Deferred as YAGNI.
- `interface{}` fields in `emData` are less type-safe than before; tracked as tech debt.
