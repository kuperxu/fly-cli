# Display Renderer

## Why a custom renderer?

`github.com/olekukonko/tablewriter` is in `go.mod` but **not used**. It was the original
renderer and was replaced because it mis-measures CJK (Chinese/Japanese/Korean) character
widths. CJK characters are "wide" — each occupies 2 terminal columns — but tablewriter
counts them as 1.

The result: columns containing Chinese stock names appear misaligned, breaking the table grid.

## Solution

`internal/display/table.go` is a hand-rolled Unicode box-drawing renderer with:

- `mattn/go-runewidth` for display-width-aware column sizing
- ANSI escape code stripping (via regex) before width measurement, so colored cells align
- Manual column padding computed from `runewidth.StringWidth()`

## Color Convention

**Chinese market convention (opposite of Western):**

| Direction | Color |
|-----------|-------|
| Price up | Red |
| Price down | Green |

Source: `github.com/fatih/color`

## Column Alignment

- Columns 0 and 1 (code, name): **left-aligned**
- All other numeric columns: **right-aligned**

## Public API

```go
display.PrintQuotes(quotes []*model.Quote)
display.PrintPortfolio(views []*model.PositionView)
display.PrintSuccess(msg string)
display.PrintError(msg string)
```

## Constraint

Do not replace this renderer with any library that does not explicitly support
`mattn/go-runewidth`-based width measurement.
