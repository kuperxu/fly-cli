# Architecture

`fly` is a layered CLI application. Each layer has a single responsibility and depends only
on layers below it.

```
┌─────────────────────────────────────────────┐
│  cmd/          CLI layer (Cobra commands)             │
│  quote.go  portfolio.go  add.go  remove.go  market.go│
└───────────────┬─────────────────────────────┘
                │ uses
    ┌───────────┴────────────────────────────┐
    │  internal/api/    Data fetching        │
    │  client.go  eastmoney.go  tencent.go   │
    └───────────┬────────────────────────────┘
                │ returns
    ┌───────────┴──────────┐  ┌─────────────────────────┐
    │  internal/model/     │  │  internal/storage/      │
    │  stock.go symbol.go  │  │  store.go               │
    └──────────────────────┘  └─────────────────────────┘
                                        used by cmd/
    ┌──────────────────────────────────────────┐
    │  internal/display/   Terminal rendering  │
    │  table.go                                │
    └──────────────────────────────────────────┘
                                        used by cmd/
```

## Packages

### `cmd/`
Cobra commands. One file per command. Commands parse flags, call `api` and `storage`,
then pass results to `display`. No business logic lives here.

### `internal/model/`
Domain types only. No I/O, no formatting.
- `stock.go` — `Quote`, `Holding`, `PositionView`, `Market`
- `symbol.go` — `ParseSymbol()`, `NormalizeCode()`

### `internal/api/`
HTTP data fetching.
- `provider.go` — `Provider` interface
- `client.go` — `Client` with primary+fallback logic
- `eastmoney.go` — primary provider (JSON)
- `tencent.go` — fallback provider (GBK text)

See [docs/design-docs/api-provider-pattern.md](docs/design-docs/api-provider-pattern.md).

### `internal/storage/`
YAML portfolio persistence at `~/.fly-cli/portfolio.yaml`.
- `Load()`, `Save()`, `Upsert()`, `Remove()`, `FindHolding()`

### `internal/display/`
Terminal table rendering with CJK-aware column sizing.
- `PrintQuotes()`, `PrintPortfolio()`, `PrintSuccess()`, `PrintError()`

See [docs/design-docs/display-renderer.md](docs/design-docs/display-renderer.md).

## Data Flow: `fly q 600519 AAPL`

```
user input
    → cmd/quote.go: parse & normalize symbols
    → api.Client.GetQuotes(): fetch from Eastmoney, fallback Tencent
    → []model.Quote
    → []model.PositionView (wraps Quote + optional Holding)
    → display.PrintQuotes()
    → terminal output
```

## Data Flow: `fly market`

```
cmd/market.go: hardcoded index definitions (secid + display code)
    → api.Client.GetIndexQuotes(): fetch via eastmoney raw secid
    → []model.Quote (Market = MarketIndex)
    → []model.PositionView (no Holding)
    → display.PrintQuotes()
    → terminal output
```

## Data Flow: `fly portfolio`

```
storage.Load() → []Holding
    → api.Client.GetQuotes() for all held codes
    → []PositionView (Quote + Holding merged)
    → display.PrintPortfolio()
    → terminal output with P&L summary
```
