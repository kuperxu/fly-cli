package cmd

import (
	"fmt"

	"fly/internal/display"
	"fly/internal/model"
	"fly/internal/storage"

	"github.com/spf13/cobra"
)

var (
	alertEntry float64
	alertTP1   float64
	alertSL    float64
)

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "管理价格提醒 (Entry/TP1/SL)",
	Long: `设置、查看、删除股票的关键价格位提醒。

子命令:
  set  设置/更新价格提醒
  ls   列出所有提醒
  rm   删除提醒`,
}

var alertSetCmd = &cobra.Command{
	Use:   "set <symbol> [--entry X] [--tp1 Y] [--sl Z]",
	Short: "设置/更新价格提醒",
	Long: `为一只股票设置 Entry（入场价）、TP1（止盈价）、SL（止损价）。
至少需要提供一个价格 flag。

示例:
  fly alert set 600519.SH --entry 1400 --tp1 1600 --sl 1350
  fly alert set AAPL --tp1 200 --sl 150
  fly alert set 00700.HK --entry 300`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if alertEntry == 0 && alertTP1 == 0 && alertSL == 0 {
			return fmt.Errorf("至少需要提供一个价格 flag: --entry, --tp1, --sl")
		}

		norm, err := model.NormalizeCode(args[0])
		if err != nil {
			return fmt.Errorf("无效的股票代码: %w", err)
		}

		store, err := storage.NewAlertStore()
		if err != nil {
			return err
		}

		alert := &model.Alert{
			Code:  norm,
			Entry: alertEntry,
			TP1:   alertTP1,
			SL:    alertSL,
		}
		if err := store.Upsert(alert); err != nil {
			return err
		}

		msg := fmt.Sprintf("已设置提醒: %s", norm)
		if alertEntry != 0 {
			msg += fmt.Sprintf("  Entry=%.2f", alertEntry)
		}
		if alertTP1 != 0 {
			msg += fmt.Sprintf("  TP1=%.2f", alertTP1)
		}
		if alertSL != 0 {
			msg += fmt.Sprintf("  SL=%.2f", alertSL)
		}
		display.PrintSuccess(msg)
		return nil
	},
}

var alertLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "列出所有价格提醒",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewAlertStore()
		if err != nil {
			return err
		}
		book, err := store.Load()
		if err != nil {
			return err
		}
		if len(book.Alerts) == 0 {
			fmt.Println("暂无价格提醒。使用 'fly alert set <代码> --entry X --tp1 Y --sl Z' 添加。")
			return nil
		}
		display.PrintAlerts(book.Alerts)
		return nil
	},
}

var alertRmCmd = &cobra.Command{
	Use:     "rm <symbol>",
	Aliases: []string{"remove", "del"},
	Short:   "删除价格提醒",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		norm, err := model.NormalizeCode(args[0])
		if err != nil {
			return fmt.Errorf("无效的股票代码: %w", err)
		}

		store, err := storage.NewAlertStore()
		if err != nil {
			return err
		}

		if err := store.Remove(norm); err != nil {
			return err
		}

		display.PrintSuccess(fmt.Sprintf("已删除提醒: %s", norm))
		return nil
	},
}

func init() {
	alertSetCmd.Flags().Float64Var(&alertEntry, "entry", 0, "入场价")
	alertSetCmd.Flags().Float64Var(&alertTP1, "tp1", 0, "止盈价")
	alertSetCmd.Flags().Float64Var(&alertSL, "sl", 0, "止损价")

	alertCmd.AddCommand(alertSetCmd)
	alertCmd.AddCommand(alertLsCmd)
	alertCmd.AddCommand(alertRmCmd)
}
