# AGENT.md

Navigation map for AI-assisted development on `fly-cli`.
This file is a table of contents — follow the links for deeper context.

## What is this project?

`fly` — a terminal CLI for real-time stock quotes and portfolio P&L tracking.
- Markets: A-shares (SH/SZ), Hong Kong, US stocks
- No API keys required
- Binary name: `fly` | Config: `~/.fly-cli/portfolio.yaml`, `~/.fly-cli/alerts.yaml`

## Key documents

| Document | What's in it |
|----------|-------------|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Package map, layer diagram, data flow |
| [docs/design-docs/core-beliefs.md](./docs/design-docs/core-beliefs.md) | Non-negotiable design principles |
| [docs/design-docs/api-provider-pattern.md](./docs/design-docs/api-provider-pattern.md) | Primary+fallback provider design |
| [docs/design-docs/display-renderer.md](./docs/design-docs/display-renderer.md) | Why the custom table renderer exists |
| [docs/references/symbol-formats.md](./docs/references/symbol-formats.md) | All symbol input formats and normalization rules |
| [docs/references/eastmoney-api.md](./docs/references/eastmoney-api.md) | Eastmoney secid format and field reference |
| [docs/references/tencent-api.md](./docs/references/tencent-api.md) | Tencent API symbol format and GBK decoding |
| [docs/exec-plans/tech-debt-tracker.md](./docs/exec-plans/tech-debt-tracker.md) | Known tech debt items |

## Commands

| Command | Aliases | File |
|---------|---------|------|
| `fly quote <symbols...>` | `fly q` | `cmd/quote.go` |
| `fly portfolio` | `fly pf`, `fly ls` | `cmd/portfolio.go` |
| `fly add <symbol>` | — | `cmd/add.go` |
| `fly remove <symbol>` | `fly rm`, `fly del` | `cmd/remove.go` |
| `fly market` | `fly mk` | `cmd/market.go` |
| `fly alert set <symbol>` | — | `cmd/alert.go` |
| `fly alert ls` | `fly alert list` | `cmd/alert.go` |
| `fly alert rm <symbol>` | `fly alert remove`, `fly alert del` | `cmd/alert.go` |

## Critical constraints (read before changing anything)

1. **Do not replace the table renderer** with `tablewriter` — it breaks CJK alignment.
   See [docs/design-docs/display-renderer.md](./docs/design-docs/display-renderer.md).
2. **Color convention is red=up, green=down** (Chinese market norm, not Western).
3. **HK codes are zero-padded to 5 digits** in all API calls.
4. **US stocks default NASDAQ (`105.`), auto-retry NYSE (`106.`)** in `eastmoney.go`.

## Agent protocol: when a requirement is complete

**Every time you finish implementing a feature, you must update the knowledge base.**
Do not ask the user — just do it as part of completing the task.

| What changed | Where to update |
|-------------|----------------|
| New feature shipped | `docs/exec-plans/completed/<slug>.md` (new file) + `docs/exec-plans/index.md` |
| New tech debt introduced | `docs/exec-plans/tech-debt-tracker.md` |
| Tech debt resolved | Mark resolved in `tech-debt-tracker.md` |
| New command added | `AGENT.md` Commands table |
| New invariant / constraint | `AGENT.md` Critical constraints |
| New package or data flow | `ARCHITECTURE.md` |
| Non-obvious design decision | New file in `docs/design-docs/` + update `index.md` |
| New external API or symbol format | `docs/references/` |

**AGENT.md must stay under ~100 lines.** Detail belongs in linked sub-documents.

## Common tasks

**Add a command:**
1. Create `cmd/<name>.go` with `var <name>Cmd = &cobra.Command{...}`
2. Register in `cmd/root.go` `init()`: `rootCmd.AddCommand(<name>Cmd)`

**Add a data provider:**
1. Implement `api.Provider` interface in `internal/api/<name>.go`
2. Wire into `api.NewClient()` in `client.go`

**Build & test:**
```bash
go build .
go install .
fly q 600519        # smoke test A-share
fly q AAPL          # smoke test US
fly portfolio       # smoke test portfolio
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/fatih/color` | ANSI terminal colors |
| `github.com/mattn/go-runewidth` | CJK-aware string width |
| `gopkg.in/yaml.v3` | Portfolio YAML serialization |
| `golang.org/x/text` | GBK → UTF-8 transcoding (Tencent API) |
