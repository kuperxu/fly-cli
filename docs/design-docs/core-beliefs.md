# Core Beliefs

Guiding principles for the `fly-cli` project. These drive every design and implementation decision.

## 1. No external dependencies for data access

All market data is sourced from free public Chinese finance APIs. No API keys, no accounts, no
rate-limit tokens. The tool must work out of the box.

**Implication:** Eastmoney and Tencent are the only allowed data providers unless a new source
is equally frictionless.

## 2. Terminal-first, not TUI

`fly` is a command you run and read in one shot — not a live-updating dashboard. Each command
makes one fresh HTTP request and prints a table. No goroutines, no refresh loops, no
alternate screen.

**Implication:** Resist adding `--watch` flags or live-poll modes unless they are explicitly
requested and scoped.

## 3. CJK correctness is non-negotiable

Stock names are in Chinese. Tables must align correctly. `tablewriter` mis-measures CJK widths,
so the renderer is hand-rolled with `mattn/go-runewidth`.

**Implication:** Never swap back to `tablewriter` or any library that doesn't handle CJK widths.

## 4. Convention follows Chinese market norms

Color coding: **red = price up, green = price down**. This matches Chinese trading apps and is
opposite of Western convention. Do not change this without user request.

## 5. Portfolio file is user-editable YAML

`~/.fly-cli/portfolio.yaml` is a simple hand-editable config. It is not a database. It does
not need migrations or versioning.

**Implication:** Schema changes must be backward-compatible or provide a migration note.

## 6. Fail fast, fail loud

Error messages are printed to stderr via Cobra's `RunE`. There is no silent swallowing.
If a symbol can't be fetched, it's reported.
