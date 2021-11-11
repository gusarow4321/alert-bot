package bot

import (
	"fmt"
	"github.com/gusarow4321/binance-alert-bot/binance"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"strings"
)

type State int

const (
	nothing State = iota
	keysUpdate
	symbolUpdate
)

var botState = nothing

func AddHandlers(b *tb.Bot, bin *binance.Binance, channelId int64, isTest bool) {
	chat, err := b.ChatByID(strconv.Itoa(int(channelId)))
	title := ""
	if err == nil {
		title = chat.Title
	}

	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnKeys := menu.Text("Подключить ключи")
	btnSymbol := menu.Text("Изменить пару")
	menu.Reply(menu.Row(btnKeys), menu.Row(btnSymbol))

	b.Handle("/start", func(m *tb.Message) {
		botState = nothing
		_, err = b.Send(
			m.Sender,
			fmt.Sprintf("Подключенный канал: %s\n\n%s\n\nВыбранная пара: %s", title, bin.IsConnectedString(), bin.GetSymbol()),
			menu,
		)
		if err != nil {
			log.Println(err)
		}
	})

	b.Handle("Подключить ключи", func(m *tb.Message) {
		botState = keysUpdate
		_, err = b.Send(m.Sender, "Отправьте ключи одним сообщением:\n\n<code>ApiKey SecretKey</code>")
		if err != nil {
			log.Println(err)
		}
	})

	b.Handle("Изменить пару", func(m *tb.Message) {
		botState = symbolUpdate
		_, err = b.Send(m.Sender, "Введите название пары")
		if err != nil {
			log.Println(err)
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		switch botState {
		case keysUpdate:
			keys := strings.Fields(m.Text)
			if len(keys) == 2 && bin.SetKeys(keys[0], keys[1], isTest) {
				botState = nothing
				_, err = b.Send(m.Sender, "Ключи подключены")
			} else {
				_, err = b.Send(m.Sender, "Неверно введены ключи\n\nОтправьте ключи одним сообщением:\n\n<code>ApiKey SecretKey</code>\n\nДля отмены - /start")
			}
			if err != nil {
				log.Println(err)
			}
		case symbolUpdate:
			if bin.SetSymbol(m.Text) {
				botState = nothing
				_, err = b.Send(m.Sender, "Пара обновлена")
			} else {
				_, err = b.Send(m.Sender, "Неверное название пары. Попробуйте еще раз\n\nДля отмены - /start")
			}
			if err != nil {
				log.Println(err)
			}
		}

		if botState == nothing {
			_, err = b.Send(
				m.Sender,
				fmt.Sprintf("Подключенный канал: %s\n\n%s\n\nВыбранная пара: %s", title, bin.IsConnectedString(), bin.GetSymbol()),
				menu,
			)
			if err != nil {
				log.Println(err)
			}
		}
	})
}
