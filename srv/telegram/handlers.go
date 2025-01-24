package telegram

import (
	"bytes"
	"fmt"
	"image/png"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/mitchellh/go-ps"
	"github.com/ranjbar-dev/gowin/config"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

var (
	// Universal markup builders.
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}

	// Reply buttons.
	btnTime          = menu.Text("Time")
	btnScreenShot    = menu.Text("Take screenshot")
	btnListProcesses = menu.Text("List processes")
)

func (t *Telegram) RegisterHandlers() {

	adminOnly := t.bot.Group()
	adminOnly.Use(middleware.Whitelist(config.TelegramChatID()))

	adminOnly.Handle("/start", func(c tele.Context) error {

		user := c.Sender()

		menu.Reply(menu.Row(btnScreenShot), menu.Row(btnListProcesses), menu.Row(btnTime))

		return c.Send("hello "+user.FirstName+" "+user.LastName, menu)
	})

	adminOnly.Handle(&btnTime, func(c tele.Context) error {

		return c.Send(fmt.Sprintf("Current time is %s", time.Now().Format("15:04:05")))
	})

	adminOnly.Handle(&btnScreenShot, func(c tele.Context) error {

		n := screenshot.NumActiveDisplays()
		if n == 0 {
			return c.Send("No active displays found")
		}

		c.Send("Processing ...")

		for i := 0; i < n; i++ {
			bounds := screenshot.GetDisplayBounds(i)
			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				return c.Send(fmt.Sprintf("Failed to capture screenshot: %v", err))
			}

			var buf bytes.Buffer
			if err := png.Encode(&buf, img); err != nil {
				return c.Send(fmt.Sprintf("Failed to encode screenshot: %v", err))
			}

			photo := &tele.Photo{File: tele.FromReader(&buf)}
			if err := c.Send(photo); err != nil {
				return c.Send(fmt.Sprintf("Failed to send photo: %v", err))
			}
		}

		c.Send("Done")

		return nil
	})

	adminOnly.Handle(&btnListProcesses, func(c tele.Context) error {

		processes, err := ps.Processes()
		if err != nil {

			return c.Send(fmt.Sprintf("Error in getting processes: %v", err))
		}

		var list string
		for i, process := range processes {

			list += fmt.Sprintf("%d: %s\n", i+1, process.Executable())
		}

		return c.Send(list) // TODO : fix message is too long 400max
	})

}
