# Initial Implementation

## Summary

Built the initial `fly` CLI from scratch.

## Scope

- `fly quote` / `fly q` — real-time stock quotes for A/HK/US markets
- `fly add` — add or update a portfolio position
- `fly portfolio` / `fly pf` / `fly ls` — view portfolio with live P&L
- `fly remove` / `fly rm` / `fly del` — remove a position

## Key Decisions

- **Eastmoney as primary, Tencent as fallback** — Eastmoney is more structured (JSON) and
  better for US stocks; Tencent covers gaps.
- **Custom table renderer** — Replaced `tablewriter` due to CJK width mis-measurement.
- **GBK decoding for Tencent** — Tencent returns GBK-encoded text; added
  `golang.org/x/text` for proper UTF-8 conversion.
- **Symbol normalization layer** — `model/symbol.go` handles the many input formats users
  might type (bare numeric, `.SH`/`.SZ`/`.HK` suffix, alphabetic ticker).

## Status

Completed. All commands functional.
