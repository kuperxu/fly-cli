# Completed: Market Closed Notice

**Date:** 2026-03-21

## What was built

Added a "休盘提示" notice to `fly market` output. When the user runs the command on a weekend (Saturday/Sunday, CST), a line is printed above the table stating the market is closed and showing the last trading day's date (the most recent weekday).

Example output:
```
休盘提示：今天休盘，以下为 03月20日 的大盘数据

┌────────┬────────────────┬ ...
```

## Key decisions

- **Calendar math only, no extra API call.** The Eastmoney ulist batch endpoint's `f86` field turned out not to be a Unix timestamp (returns values like `25.47`). The single-stock endpoint does return a real Unix timestamp in `f86`, but making a separate HTTP request just for this was not justified. Weekend detection is reliable and covers the most common case.

- **Function signature accepts quotes but ignores them.** The signature `func printMarketClosedNotice(_ []*model.Quote)` was chosen to leave the door open for future use of quote data (e.g., detecting holidays via Price == PrevClose heuristic) without changing the call site.

## Trade-offs

- **Public holidays not detected.** Only Saturday/Sunday triggers the notice. Chinese public holidays (Golden Week, Spring Festival, etc.) will silently show stale data. This is a known gap — see tech-debt TD-008.

## Files changed

- `cmd/market.go` — added `printMarketClosedNotice()` helper
- No model or API changes required
