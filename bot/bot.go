package bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func InitBot(token string, pollingTimeout, adminId int) (*tb.Bot, error) {
	poller := &tb.LongPoller{Timeout: time.Duration(pollingTimeout) * time.Second}
	adminAccess := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		if upd.Message != nil && upd.Message.Sender != nil && upd.Message.Sender.ID != adminId {
			return false
		}
		return true
	})

	b, err := tb.NewBot(tb.Settings{
		Token:     token,
		Poller:    adminAccess,
		ParseMode: tb.ModeHTML,
	})

	if err != nil {
		return nil, err
	}
	return b, nil
}
