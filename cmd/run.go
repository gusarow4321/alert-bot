package cmd

import (
	"github.com/gusarow4321/binance-alert-bot/binance"
	"github.com/gusarow4321/binance-alert-bot/bot"
	"github.com/spf13/cobra"
	"log"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run bot",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString("bot-token")
		if err != nil {
			panic(err)
		}
		timeout, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			panic(err)
		}
		adminId, err := cmd.Flags().GetInt("admin-id")
		if err != nil {
			panic(err)
		}
		channelId, err := cmd.Flags().GetInt64("channel-id")
		if err != nil {
			panic(err)
		}
		fromTimestamp, err := cmd.Flags().GetString("from-timestamp")
		if err != nil {
			panic(err)
		}
		testnet, err := cmd.Flags().GetBool("testnet")
		if err != nil {
			panic(err)
		}

		b, err := bot.InitBot(token, timeout, adminId)
		if err != nil {
			log.Println(err) // TODO
			return
		}

		bin := binance.NewBinance("", fromTimestamp, nil)

		bot.AddHandlers(b, bin, channelId, testnet)

		go b.Start()
		go bin.Start(b, channelId)

		select {}
	},
}
