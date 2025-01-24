package telegram

import (
	"time"

	"github.com/ranjbar-dev/gowin/config"
	tele "gopkg.in/telebot.v4"
)

type Telegram struct {
	bot *tele.Bot
}

func (t *Telegram) Start() {

	t.RegisterHandlers()

	go t.bot.Start()
}

func NewTelegram() *Telegram {

	pref := tele.Settings{
		Token:  config.TelegramBotToken(),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {

		panic(err)
	}

	return &Telegram{
		bot: bot,
	}
}
