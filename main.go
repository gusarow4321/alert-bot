package main

import (
	"github.com/gusarow4321/binance-alert-bot/cmd"
)

func main() {
	rootCmd := cmd.RootCmd

	cmd.RunCmd.Flags().StringP("bot-token", "b", "", "bot token")
	cmd.RunCmd.Flags().IntP("timeout", "t", 10, "polling timeout (default is 10)")
	cmd.RunCmd.Flags().IntP("admin-id", "a", 0, "admin user id")
	cmd.RunCmd.Flags().Int64P("channel-id", "c", 0, "channel chat id")
	cmd.RunCmd.Flags().StringP("from-timestamp", "f", "", "UTC timestamp to start with (format: YYYY-MM-DD hh:mm:ss, default is now)")
	cmd.RunCmd.Flags().Bool("testnet", false, "is testnet futures api (default is false)")

	rootCmd.AddCommand(
		cmd.RunCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
