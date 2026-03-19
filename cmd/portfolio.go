package cmd

import (
	"fmt"

	"fly/internal/api"
	"fly/internal/display"
	"fly/internal/model"
	"fly/internal/storage"

	"github.com/spf13/cobra"
)

var portfolioCmd = &cobra.Command{
	Use:     "portfolio",
	Aliases: []string{"pf", "ls"},
	Short:   "查看持仓及盈亏",
	Long: `显示所有持仓的实时行情和盈亏情况。

示例:
  fly portfolio
  fly pf
  fly ls`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStore()
		if err != nil {
			return err
		}
		portfolio, err := store.Load()
		if err != nil {
			return err
		}
		if len(portfolio.Holdings) == 0 {
			fmt.Println("持仓为空。使用 'fly add <代码> --cost <成本> --shares <数量>' 添加持仓。")
			return nil
		}

		// Collect symbols from portfolio
		symbols := make([]string, 0, len(portfolio.Holdings))
		for _, h := range portfolio.Holdings {
			symbols = append(symbols, h.Code)
		}

		// Fetch quotes
		client := api.NewClient()
		quotes, err := client.GetQuotes(symbols)
		if err != nil {
			return fmt.Errorf("获取行情失败: %w", err)
		}

		// Build position views preserving holding order
		views := make([]*model.PositionView, 0, len(portfolio.Holdings))
		quoteMap := make(map[string]*model.Quote)
		for _, q := range quotes {
			sym := q.Code + "." + string(q.Market)
			if q.Market == model.MarketUS {
				sym = q.Code
			}
			quoteMap[sym] = q
		}
		for _, h := range portfolio.Holdings {
			q, ok := quoteMap[h.Code]
			if !ok {
				fmt.Printf("警告: 未能获取 %s 的行情\n", h.Code)
				continue
			}
			views = append(views, &model.PositionView{
				Quote:   q,
				Holding: h,
			})
		}

		display.PrintPortfolio(views)

		// Check alerts
		alertStore, err := storage.NewAlertStore()
		if err == nil {
			book, err := alertStore.Load()
			if err == nil && len(book.Alerts) > 0 {
				var triggered []model.TriggeredAlert
				for _, v := range views {
					sym := v.Quote.Code + "." + string(v.Quote.Market)
					if v.Quote.Market == model.MarketUS {
						sym = v.Quote.Code
					}
					a := book.FindAlert(sym)
					if a == nil {
						continue
					}
					if a.Entry != 0 && v.Quote.Price <= a.Entry {
						triggered = append(triggered, model.TriggeredAlert{
							Code: sym, Name: v.Quote.Name, Price: v.Quote.Price,
							Type: "Entry", Target: a.Entry,
						})
					}
					if a.TP1 != 0 && v.Quote.Price >= a.TP1 {
						triggered = append(triggered, model.TriggeredAlert{
							Code: sym, Name: v.Quote.Name, Price: v.Quote.Price,
							Type: "TP1", Target: a.TP1,
						})
					}
					if a.SL != 0 && v.Quote.Price <= a.SL {
						triggered = append(triggered, model.TriggeredAlert{
							Code: sym, Name: v.Quote.Name, Price: v.Quote.Price,
							Type: "SL", Target: a.SL,
						})
					}
				}
				display.PrintTriggeredAlerts(triggered)
			}
		}

		return nil
	},
}
