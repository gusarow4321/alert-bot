package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

type Binance struct {
	symbol string
	lastTS int64
	client *futures.Client
}

func timestampToInt64(ts string) int64 {
	if ts == "" {
		return time.Now().UnixMilli()
	}

	t, err := time.Parse("2006-01-02 15:04:05", ts)
	if err != nil {
		log.Println(err)
		return time.Now().UnixMilli()
	}

	return t.UnixMilli()
}

func NewBinance(Symbol string, ts string, Client *futures.Client) *Binance {
	return &Binance{
		symbol: Symbol,
		lastTS: timestampToInt64(ts),
		client: Client,
	}
}

func (b *Binance) Ready() bool {
	return b.client != nil && b.symbol != ""
}

func (b *Binance) IsConnectedString() string {
	if b.client == nil {
		return "❗️Подключите ключи binance api"
	}
	return "Binance api подключен"
}

func (b *Binance) SetKeys(api, secret string, isTest bool) bool {
	if isTest {
		futures.UseTestnet = true
	}
	newClient := futures.NewClient(api, secret)
	_, err := newClient.NewListAccountTradeService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		return false
	}
	b.client = newClient
	return true
}

func (b *Binance) GetSymbol() string {
	return b.symbol
}

func (b *Binance) SetSymbol(s string) bool {
	if b.client == nil {
		return false
	}

	_, err := b.client.NewPremiumIndexService().Symbol(s).Do(context.Background())
	if err != nil {
		return false
	}

	b.symbol = s
	return true
}

// https://testnet.binancefuture.com/en/futures/BTCUSDT
// https://binance-docs.github.io/apidocs/futures/en/#get-income-history-user_data
// https://binance-docs.github.io/apidocs/futures/en/#account-trade-list-user_data

func (b *Binance) Start(bot *tb.Bot, channelId int64) {
	for {
		if !b.Ready() || channelId == 0 {
			time.Sleep(10 * time.Second)
			continue
		}

		trades, err := b.client.NewListAccountTradeService().Symbol(b.symbol).StartTime(b.lastTS).Do(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}

		lastOpen := ""
		maxMessagesCount := 10

		for i, t := range trades {
			if i == maxMessagesCount {
				break
			}

			b.lastTS = t.Time + 1

			if t.RealizedPnl == "0" {
				position := "LONG"
				if t.Side == "SELL" {
					position = "SHORT"
				}
				open := fmt.Sprintf("%s\nОткрываем позицию в %s по цене %s", b.symbol, position, t.Price)
				if open != lastOpen {
					lastOpen = open
					_, err = bot.Send(&tb.Chat{ID: channelId}, open)
					if err != nil {
						log.Println(err)
					}
				}
			} else {
				_, err = bot.Send(&tb.Chat{ID: channelId}, fmt.Sprintf("%s\nПозиция закрыта по цене %s", b.symbol, t.Price))
				if err != nil {
					log.Println(err)
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}
