package telegram

import (
	"log"
	"time"

	"github.com/ranjbar-dev/gowin/config"
	tele "gopkg.in/telebot.v4"
)

var bot *tele.Bot
var chatID int64

func init() {

	pref := tele.Settings{
		Token:  config.TelegramBotToken(),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {

		log.Fatal(err)
		return
	}

	bot = b
	chatID = config.TelegramChatID()
}

func SendMessage(message string) {

	_, err := bot.Send(&tele.User{ID: chatID}, message)
	if err != nil {

		log.Println(err)
	}
}
