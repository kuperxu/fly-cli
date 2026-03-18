package cmd

import (
	"fmt"

	"fly/internal/display"
	"fly/internal/model"
	"fly/internal/storage"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <symbol>",
	Aliases: []string{"rm", "del"},
	Short:   "删除持仓",
	Long: `从持仓列表中删除一只股票。

示例:
  fly remove 600519.SH
  fly rm AAPL
  fly del 00700.HK`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		norm, err := model.NormalizeCode(args[0])
		if err != nil {
			return fmt.Errorf("无效的股票代码: %w", err)
		}

		store, err := storage.NewStore()
		if err != nil {
			return err
		}

		if err := store.Remove(norm); err != nil {
			return err
		}

		display.PrintSuccess(fmt.Sprintf("已删除持仓: %s", norm))
		return nil
	},
}
