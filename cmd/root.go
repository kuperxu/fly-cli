package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fly",
	Short: "股票行情查询工具",
	Long:  "fly - 在终端查看 A股、港股、美股实时行情及持仓盈亏",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(quoteCmd)
	rootCmd.AddCommand(portfolioCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(marketCmd)
	rootCmd.AddCommand(alertCmd)
}
