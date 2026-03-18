package cmd

import (
	"fmt"
	"os"

	"fly/internal/api"
	"fly/internal/display"
	"fly/internal/model"
	"fly/internal/storage"

	"github.com/spf13/cobra"
)

var quoteCmd = &cobra.Command{
	Use:     "quote [symbols...]",
	Aliases: []string{"q"},
	Short:   "查询股票实时行情",
	Long: `查询一只或多只股票的实时行情。

示例:
  fly quote 600519
  fly quote 600519.SH 000858.SZ
  fly quote 00700.HK
  fly quote AAPL TSLA
  fly q 600519 00700.HK AAPL`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Normalize all symbols
		symbols := make([]string, 0, len(args))
		for _, arg := range args {
			norm, err := model.NormalizeCode(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "警告: 无法识别 %s: %v\n", arg, err)
				continue
			}
			symbols = append(symbols, norm)
		}
		if len(symbols) == 0 {
			return fmt.Errorf("没有有效的股票代码")
		}

		// Load portfolio to show P&L for held stocks
		store, err := storage.NewStore()
		if err != nil {
			return err
		}
		portfolio, err := store.Load()
		if err != nil {
			return err
		}

		// Fetch quotes
		client := api.NewClient()
		quotes, err := client.GetQuotes(symbols)
		if err != nil {
			return fmt.Errorf("获取行情失败: %w", err)
		}
		if len(quotes) == 0 {
			return fmt.Errorf("未获取到任何行情数据")
		}

		// Build position views
		views := make([]*model.PositionView, 0, len(quotes))
		for _, q := range quotes {
			sym := q.Code + "." + string(q.Market)
			if q.Market == model.MarketUS {
				sym = q.Code
			}
			holding := portfolio.FindHolding(sym)
			views = append(views, &model.PositionView{
				Quote:   q,
				Holding: holding,
			})
		}

		display.PrintQuotes(views)
		return nil
	},
}
