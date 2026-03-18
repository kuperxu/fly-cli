# AGENT.md

Project context for AI-assisted development.

## Project

`fly` — a terminal CLI for real-time stock quotes and portfolio P&L tracking.
Written in Go. Binary name: `fly`. Config dir: `~/.fly-cli/`.

## Architecture

```
fly-cli/
├── main.go                      # Entrypoint
├── cmd/                         # Cobra commands (one file per command)
│   ├── root.go                  # Root command, registers all subcommands
│   ├── quote.go                 # `fly quote` / `fly q`
│   ├── portfolio.go             # `fly portfolio` / `fly pf` / `fly ls`
│   ├── add.go                   # `fly add`
│   └── remove.go                # `fly remove` / `fly rm` / `fly del`
└── internal/
    ├── model/
    │   ├── stock.go             # Core types: Quote, Holding, PositionView
    │   └── symbol.go            # Symbol parsing and normalization
    ├── api/
    │   ├── provider.go          # Provider interface
    │   ├── client.go            # Client with primary + fallback logic
    │   ├── eastmoney.go         # Eastmoney push2 API (primary, JSON)
    │   └── tencent.go           # Tencent Finance API (fallback, text)
    ├── storage/
    │   └── store.go             # YAML portfolio read/write (~/.fly-cli/portfolio.yaml)
    └── display/
        └── table.go             # Colored, CJK-aligned table rendering
```

## Key types

```go
// model/stock.go
type Market string  // "SH" | "SZ" | "HK" | "US"

type Quote struct {
    Code, Name, Currency string
    Market               Market
    Price, PrevClose, Open, High, Low float64
    Volume               int64
    Amount, Change, ChangePct float64
}

type Holding struct {
    Code   string  `yaml:"code"`   // e.g. "600519.SH", "AAPL"
    Cost   float64 `yaml:"cost"`
    Shares float64 `yaml:"shares"`
}

type PositionView struct {
    Quote   *Quote
    Holding *Holding  // nil if not in portfolio
}
```

## Data sources

| Provider | API endpoint | Markets | Format |
|----------|-------------|---------|--------|
| Eastmoney (primary) | `push2.eastmoney.com/api/qt/stock/get` | A/HK/US | JSON |
| Tencent (fallback) | `qt.gtimg.cn/q=` | A/HK/US | `~`-delimited text |

**Eastmoney secid format:**
- Shanghai: `1.600519`
- Shenzhen: `0.000858`
- HK: `116.00700`
- NASDAQ: `105.AAPL`
- NYSE: `106.BRK-A` (fallback tried automatically)

**Tencent symbol format:** `sh600519`, `sz000858`, `hk00700`, `usAAPL`

## Symbol parsing rules (`model/symbol.go`)

| Input | Resolved market |
|-------|----------------|
| `600519` | SH (starts with 6/5/9) |
| `000858` | SZ (other numeric) |
| `600519.SH` | SH (explicit) |
| `000858.SZ` | SZ (explicit) |
| `00700.HK` | HK (explicit) |
| `AAPL` | US (alphabetic) |

Canonical form: `CODE.MARKET` (e.g. `600519.SH`) except US tickers which are bare (e.g. `AAPL`).

## Display (`internal/display/table.go`)

- Hand-rolled table renderer — **do not switch back to tablewriter**, it mis-measures CJK widths.
- Uses `mattn/go-runewidth` for display-width-aware column sizing.
- ANSI codes stripped via regex before width measurement so colored cells align correctly.
- Color convention: **red = up, green = down** (Chinese market convention).
- First two columns (code, name) are left-aligned; all others right-aligned.

## Storage (`internal/storage/store.go`)

- Config file: `~/.fly-cli/portfolio.yaml`
- `Store.Load()` returns empty portfolio (not error) when file doesn't exist.
- `Store.Upsert()` matches by normalized (uppercased) code.

## Adding a new command

1. Create `cmd/<name>.go` with a `var <name>Cmd = &cobra.Command{...}`.
2. Register it in `cmd/root.go` under `func init()` with `rootCmd.AddCommand(<name>Cmd)`.

## Adding a new data provider

1. Implement the `api.Provider` interface (`Name() string`, `GetQuotes([]string) ([]*model.Quote, error)`).
2. Wire it into `api.NewClient()` in `client.go` as primary or fallback.

## Build & run

```bash
go build .          # build binary in current dir
go install .        # install to $GOPATH/bin (already in PATH)
fly q 600519        # smoke test
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/fatih/color` | ANSI terminal colors |
| `github.com/mattn/go-runewidth` | CJK-aware string width |
| `gopkg.in/yaml.v3` | Portfolio config serialization |
| `github.com/olekukonko/tablewriter` | In go.mod but **not used** — replaced by custom renderer |

## Known constraints

- No real-time streaming — each command makes a fresh HTTP request.
- Eastmoney US stocks default to NASDAQ (`105.`); NYSE stocks auto-retry with `106.`.
- HK codes are zero-padded to 5 digits in API calls.
- Portfolio file is not locked — concurrent writes from multiple terminals can race.
