package cmd

import (
	"fmt"

	"fly/internal/display"
	"fly/internal/model"
	"fly/internal/storage"

	"github.com/spf13/cobra"
)

var (
	addCost   float64
	addShares float64
)

var addCmd = &cobra.Command{
	Use:   "add <symbol> --cost <price> --shares <qty>",
	Short: "添加或更新持仓",
	Long: `添加或更新一只股票的持仓信息。

示例:
  fly add 600519.SH --cost 1800 --shares 100
  fly add 00700.HK --cost 320 --shares 200
  fly add AAPL --cost 150 --shares 50`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if addCost <= 0 {
			return fmt.Errorf("--cost 必须大于 0")
		}
		if addShares <= 0 {
			return fmt.Errorf("--shares 必须大于 0")
		}

		norm, err := model.NormalizeCode(args[0])
		if err != nil {
			return fmt.Errorf("无效的股票代码: %w", err)
		}

		store, err := storage.NewStore()
		if err != nil {
			return err
		}

		holding := &model.Holding{
			Code:   norm,
			Cost:   addCost,
			Shares: addShares,
		}
		if err := store.Upsert(holding); err != nil {
			return err
		}

		display.PrintSuccess(fmt.Sprintf("已保存持仓: %s  成本 %.2f  数量 %.0f", norm, addCost, addShares))
		return nil
	},
}

func init() {
	addCmd.Flags().Float64Var(&addCost, "cost", 0, "持仓成本价 (必填)")
	addCmd.Flags().Float64Var(&addShares, "shares", 0, "持股数量 (必填)")
	_ = addCmd.MarkFlagRequired("cost")
	_ = addCmd.MarkFlagRequired("shares")
}
