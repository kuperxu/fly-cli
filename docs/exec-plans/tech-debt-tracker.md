# Tech Debt Tracker

Known issues, shortcuts, and deferred improvements.

## Open Items

| ID | Area | Description | Priority |
|----|------|-------------|----------|
| TD-001 | storage | Portfolio file has no write lock — concurrent writes from multiple terminals can race | Low |
| TD-002 | deps | `tablewriter` remains in `go.mod` as unused indirect dep | Low |
| TD-003 | api | Eastmoney NYSE auto-retry (`105.` → `106.`) is a heuristic; some US ETFs may fail silently | Medium |
| TD-004 | display | No pagination for large portfolios — all rows printed at once | Low |
| TD-005 | api | `emData.Volume`/`Amount` use `interface{}` to handle eastmoney returning `"-"` for some indices — fragile JSON parsing | Low |
| TD-006 | storage | Alerts file has no write lock — same race condition as TD-001 | Low |
| TD-007 | alert | Zero value (`0`) means "not set" — cannot alert on price 0.00 (acceptable for real stocks) | Low |

## Resolved Items

_None yet._
