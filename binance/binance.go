package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"sort"
	"time"
)

type Binance struct {
	symbols []string
	lastTS  int64
	client  *futures.Client
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

func NewBinance(Symbols []string, ts string, Client *futures.Client) *Binance {
	return &Binance{
		symbols: Symbols,
		lastTS:  timestampToInt64(ts),
		client:  Client,
	}
}

func (b *Binance) Ready() bool {
	return b.client != nil && len(b.symbols) != 0
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

func (b *Binance) GetSymbols() []string {
	return b.symbols
}

func (b *Binance) AddSymbol(s string) bool {
	if b.client == nil {
		return false
	}

	_, err := b.client.NewPremiumIndexService().Symbol(s).Do(context.Background())
	if err != nil {
		return false
	}

	for _, v := range b.symbols {
		if v == s {
			return true
		}
	}

	b.symbols = append(b.symbols, s)
	return true
}

func (b *Binance) RemoveSymbol(s string) {
	for i, v := range b.symbols {
		if v == s {
			b.symbols = append(b.symbols[:i], b.symbols[i+1:]...)
			break
		}
	}
}

// https://testnet.binancefuture.com/en/futures/BTCUSDT
// https://binance-docs.github.io/apidocs/futures/en/#get-income-history-user_data
// https://binance-docs.github.io/apidocs/futures/en/#account-trade-list-user_data

var (
	lastOpenTS  map[string]int64
	lastCloseTS map[string]int64
)

func (b *Binance) Start(bot *tb.Bot, channelId int64) {
	for {
		if !b.Ready() || channelId == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
		var trades []*futures.AccountTrade
		for _, symbol := range b.symbols {
			t, err := b.client.NewListAccountTradeService().Symbol(symbol).StartTime(b.lastTS).Do(context.Background())
			if err != nil {
				log.Println(err)
				continue
			}
			trades = append(trades, t...)
		}

		sort.Slice(trades, func(i, j int) bool {
			return trades[i].Time < trades[j].Time
		})

		maxMessagesCount := 5
		cnt := 0

		for _, t := range trades {
			if cnt == maxMessagesCount {
				break
			}

			cnt++
			b.lastTS = t.Time + 1

			if t.RealizedPnl == "0" {
				position := "LONG"
				if t.Side == "SELL" {
					position = "SHORT"
				}
				if t.Time-lastOpenTS[t.Symbol] > 11000 { // 11 sec
					lastOpenTS[t.Symbol] = t.Time
					_, err := bot.Send(&tb.Chat{ID: channelId}, fmt.Sprintf("%s\nОткрываем позицию в %s по цене %s", t.Symbol, position, t.Price))
					if err != nil {
						log.Println(err)
					}
				}
			} else {
				if t.Time-lastCloseTS[t.Symbol] > 11000 {
					lastCloseTS[t.Symbol] = t.Time
					_, err := bot.Send(&tb.Chat{ID: channelId}, fmt.Sprintf("%s\nПозиция закрыта по цене %s", t.Symbol, t.Price))
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}
