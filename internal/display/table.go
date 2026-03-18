package display

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"

	"fly/internal/model"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

var (
	colorUp      = color.New(color.FgRed)
	colorDown    = color.New(color.FgGreen)
	colorNeutral = color.New(color.Reset)
	colorBold    = color.New(color.Bold)
	colorCyan    = color.New(color.FgCyan)
	colorDim     = color.New(color.FgHiBlack)

	// Strip ANSI escape sequences when measuring string display width
	ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

// visibleWidth returns the terminal display width of a string, ignoring ANSI codes.
func visibleWidth(s string) int {
	return runewidth.StringWidth(ansiRe.ReplaceAllString(s, ""))
}

// padLeft right-aligns s within a cell of width w (display-width aware).
func padLeft(s string, w int) string {
	vw := visibleWidth(s)
	if vw >= w {
		return s
	}
	return strings.Repeat(" ", w-vw) + s
}

// padRight left-aligns s within a cell of width w (display-width aware).
func padRight(s string, w int) string {
	vw := visibleWidth(s)
	if vw >= w {
		return s
	}
	return s + strings.Repeat(" ", w-vw)
}

// padCenter centers s within a cell of width w (display-width aware).
func padCenter(s string, w int) string {
	vw := visibleWidth(s)
	if vw >= w {
		return s
	}
	total := w - vw
	left := total / 2
	right := total - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

type column struct {
	header    string
	rightAlign bool // false = center header + right-align data; true = right-align both
}

var columns = []column{
	{"代码", false},
	{"名称", false},
	{"当前价", true},
	{"涨跌额", true},
	{"涨跌%", true},
	{"最高", true},
	{"最低", true},
	{"成本", true},
	{"持股", true},
	{"市值", true},
	{"盈亏", true},
	{"盈亏%", true},
}

// PrintQuotes prints a nicely aligned table of stock quotes.
func PrintQuotes(views []*model.PositionView) {
	rows := make([][]string, len(views))
	for i, v := range views {
		rows[i] = buildQuoteRow(v)
	}
	printTable(rows)
}

// PrintPortfolio prints portfolio with a summary line.
func PrintPortfolio(views []*model.PositionView) {
	if len(views) == 0 {
		fmt.Println("持仓为空。使用 'fly add <代码> --cost <价格> --shares <数量>' 添加持仓。")
		return
	}
	PrintQuotes(views)

	var totalCost, totalValue, totalPnL float64
	for _, v := range views {
		if v.Holding != nil {
			totalCost += v.CostValue()
			totalValue += v.MarketValue()
			pnl, _ := v.PnL()
			totalPnL += pnl
		}
	}

	fmt.Println()
	colorBold.Print("汇总  ")
	fmt.Printf("总成本: %.2f  ", totalCost)
	fmt.Printf("总市值: %.2f  ", totalValue)
	colorChange(totalPnL).Printf("总盈亏: %s (%.2f%%)\n",
		formatPnL(totalPnL),
		safePct(totalPnL, totalCost),
	)
}

// PrintError prints an error message to stderr.
func PrintError(msg string) {
	color.New(color.FgRed).Fprintf(os.Stderr, "错误: %s\n", msg)
}

// PrintSuccess prints a success message.
func PrintSuccess(msg string) {
	color.New(color.FgGreen).Println(msg)
}

// printTable renders a table with auto-sized columns using runewidth for correct CJK alignment.
func printTable(rows [][]string) {
	n := len(columns)

	// Compute column widths: max of header width and each row's cell width
	widths := make([]int, n)
	for i, col := range columns {
		widths[i] = visibleWidth(col.header)
	}
	for _, row := range rows {
		for i := 0; i < n && i < len(row); i++ {
			if w := visibleWidth(row[i]); w > widths[i] {
				widths[i] = w
			}
		}
	}

	sep := colorDim.Sprint("│")
	topLine := buildBorderLine("┌", "┬", "┐", widths)
	midLine := buildBorderLine("├", "┼", "┤", widths)
	botLine := buildBorderLine("└", "┴", "┘", widths)

	// Top border
	colorDim.Println(topLine)

	// Header row
	fmt.Print(sep)
	for i, col := range columns {
		cell := padCenter(col.header, widths[i])
		fmt.Printf(" %s ", colorBold.Sprint(cell))
		fmt.Print(sep)
	}
	fmt.Println()

	// Header/body separator
	colorDim.Println(midLine)

	// Data rows
	for _, row := range rows {
		fmt.Print(sep)
		for i := 0; i < n; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			// First two columns left-align, rest right-align
			if i < 2 {
				fmt.Printf(" %s ", padRight(cell, widths[i]))
			} else {
				fmt.Printf(" %s ", padLeft(cell, widths[i]))
			}
			fmt.Print(sep)
		}
		fmt.Println()
	}

	// Bottom border
	colorDim.Println(botLine)
}

func buildBorderLine(left, mid, right string, widths []int) string {
	var sb strings.Builder
	sb.WriteString(left)
	for i, w := range widths {
		sb.WriteString(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			sb.WriteString(mid)
		}
	}
	sb.WriteString(right)
	return sb.String()
}

func buildQuoteRow(v *model.PositionView) []string {
	q := v.Quote
	symbol := q.Code + "." + string(q.Market)
	if q.Market == model.MarketUS {
		symbol = q.Code
	}

	changeStr := colorizeChange(fmt.Sprintf("%+.2f", q.Change), q.Change)
	changePctStr := colorizeChange(fmt.Sprintf("%+.2f%%", q.ChangePct), q.ChangePct)
	priceStr := colorizeChange(fmt.Sprintf("%.2f", q.Price), q.Change)

	row := []string{
		colorCyan.Sprint(symbol),
		q.Name,
		priceStr,
		changeStr,
		changePctStr,
		fmt.Sprintf("%.2f", q.High),
		fmt.Sprintf("%.2f", q.Low),
	}

	if v.Holding != nil {
		pnl, pnlPct := v.PnL()
		row = append(row,
			fmt.Sprintf("%.2f", v.Holding.Cost),
			fmt.Sprintf("%.0f", v.Holding.Shares),
			fmt.Sprintf("%.2f", v.MarketValue()),
			colorizeChange(fmt.Sprintf("%+.2f", pnl), pnl),
			colorizeChange(fmt.Sprintf("%+.2f%%", pnlPct), pnl),
		)
	} else {
		row = append(row, "-", "-", "-", "-", "-")
	}

	return row
}

func colorizeChange(s string, val float64) string {
	if val > 0 {
		return colorUp.Sprint(s)
	}
	if val < 0 {
		return colorDown.Sprint(s)
	}
	return s
}

func colorChange(v float64) *color.Color {
	if v > 0 {
		return colorUp
	}
	if v < 0 {
		return colorDown
	}
	return colorNeutral
}

func formatPnL(v float64) string {
	if v >= 0 {
		return fmt.Sprintf("+%.2f", v)
	}
	return fmt.Sprintf("%.2f", v)
}

func safePct(numerator, denominator float64) float64 {
	if denominator == 0 || math.IsNaN(denominator) {
		return 0
	}
	return numerator / denominator * 100
}
