package cmd

import (
	"fmt"
	"time"

	"fly/internal/api"
	"fly/internal/display"
	"fly/internal/model"

	"github.com/spf13/cobra"
)

var defaultIndices = []string{
	"1.000001",  // 上证指数
	"1.000688",  // 科创50
	"100.HSI",   // 恒生指数
	"100.NDX",   // 纳斯达克
	"171.CN10Y", // 中国10年期国债
}

// printMarketClosedNotice prints a notice when today is a weekend (A-share market closed).
// It calculates the last weekday as the most recent trading day reference.
func printMarketClosedNotice(_ []*model.Quote) {
	cst := time.FixedZone("CST", 8*3600)
	now := time.Now().In(cst)

	wd := now.Weekday()
	if wd != time.Saturday && wd != time.Sunday {
		return
	}

	// Walk back to the most recent weekday (Friday for Sat/Sun)
	lastTrading := now
	for {
		lastTrading = lastTrading.AddDate(0, 0, -1)
		if wd := lastTrading.Weekday(); wd != time.Saturday && wd != time.Sunday {
			break
		}
	}
	fmt.Printf("休盘提示：今天休盘，以下为 %s 的大盘数据\n\n", lastTrading.Format("01月02日"))
}

var marketCmd = &cobra.Command{
	Use:     "market",
	Aliases: []string{"mk"},
	Short:   "查看今日大盘行情",
	Long: `查看今日主要指数行情：上证指数、科创50、恒生指数、纳斯达克、中国10年期国债。

示例:
  fly market
  fly mk`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient()
		quotes, err := client.GetIndexQuotes(defaultIndices)
		if err != nil {
			return fmt.Errorf("获取大盘行情失败: %w", err)
		}
		if len(quotes) == 0 {
			return fmt.Errorf("未获取到任何行情数据")
		}

		views := make([]*model.PositionView, 0, len(quotes))
		for _, q := range quotes {
			views = append(views, &model.PositionView{Quote: q})
		}

		printMarketClosedNotice(quotes)
		display.PrintMarketQuotes(views)
		return nil
	},
}
