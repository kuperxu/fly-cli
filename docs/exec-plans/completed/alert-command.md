# Completed: `fly alert` — Price Alert Command

## What was built

A price alert system that lets users set Entry (entry price), TP1 (take-profit), and SL (stop-loss)
levels for stocks. Alerts are checked automatically when running `fly pf` and triggered alerts
are displayed as warnings after the portfolio table.

## Commands

- `fly alert set <symbol> --entry X --tp1 Y --sl Z` — set/update alerts (at least one flag required)
- `fly alert ls` — list all configured alerts
- `fly alert rm <symbol>` — remove an alert

## Key decisions

1. **Separate storage file** (`alerts.yaml`) — keeps alert config independent from portfolio holdings.
   Users may want alerts on stocks they don't hold yet (e.g. entry signals).

2. **Merge semantics on upsert** — when updating an existing alert, only non-zero flags overwrite.
   This lets users do `fly alert set AAPL --tp1 200` without clearing existing SL/Entry values.

3. **Passive checking only** — alerts are evaluated when `fly pf` runs, not via background daemon.
   This keeps the tool simple and stateless (no running processes).

4. **Trigger logic**:
   - Entry: current price <= entry (price has come down to your buy zone)
   - TP1: current price >= tp1 (price has reached take-profit)
   - SL: current price <= sl (price has hit stop-loss)

## Files changed

| File | Change |
|------|--------|
| `internal/model/alert.go` | New — `Alert`, `TriggeredAlert` structs |
| `internal/storage/alert_store.go` | New — `AlertStore` with CRUD for `alerts.yaml` |
| `cmd/alert.go` | New — `alert set/ls/rm` subcommands |
| `cmd/root.go` | Modified — register `alertCmd` |
| `cmd/portfolio.go` | Modified — check alerts after portfolio display |
| `internal/display/table.go` | Modified — `PrintAlerts()`, `PrintTriggeredAlerts()` |

## Trade-offs

- No background/push notification — user must run `fly pf` to see alerts.
- Zero value (`0`) means "not set", so users cannot set an alert at exactly 0.00 price.
  This is acceptable since no tradeable stock has a price of zero.
