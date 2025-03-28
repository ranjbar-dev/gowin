package telegram

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-vgo/robotgo"
	"github.com/mitchellh/go-ps"
	"github.com/ranjbar-dev/gowin/config"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

var (
	user32              = syscall.MustLoadDLL("user32.dll")
	procLockWorkStation = user32.MustFindProc("LockWorkStation")

	// Universal markup builders.
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}

	// Reply buttons.
	btnHelp          = menu.Text("/help")
	btnTime          = menu.Text("/time")
	btnListProcesses = menu.Text("/processes")
	btnMessage       = menu.Text("/message")
	btnTypeWrite     = menu.Text("/write")
	btnLock          = menu.Text("/lock")
)

func (t *Telegram) RegisterHandlers() {

	adminOnly := t.bot.Group()
	adminOnly.Use(middleware.Whitelist(config.TelegramChatID()))

	adminOnly.Handle("/start", func(c tele.Context) error {

		user := c.Sender()

		text := fmt.Sprintf("hello %s %s, you can call /help to see the commands", user.FirstName, user.LastName)

		return c.Send(text, menu)
	})

	adminOnly.Handle(&btnHelp, func(c tele.Context) error {

		text := "you can use the following commands:\n"
		text += "/help - show this help\n"
		text += "/time - show the system time\n"
		text += "/processes - list all processes\n"
		text += "/message - open a new message box and put the text in it\n"
		text += "/write - type text using the keyboard\n"
		text += "/lock - lock the screen\n"

		return c.Send(text)
	})

	adminOnly.Handle(&btnTime, func(c tele.Context) error {

		return c.Send(fmt.Sprintf("Current time is %s", time.Now().Format("15:04:05")))
	})

	adminOnly.Handle(&btnListProcesses, func(c tele.Context) error {

		processes, err := ps.Processes()
		if err != nil {

			return c.Send(fmt.Sprintf("Error in getting processes: %v", err))
		}

		localMap := make(map[string]struct{})
		for _, process := range processes {

			localMap[process.Executable()] = struct{}{}
		}

		var list string
		var i int
		for process := range localMap {

			i++
			list += fmt.Sprintf("%d: %s\n", i+1, process)
		}

		return c.Send(list) // TODO : fix message is too long 400max
	})

	adminOnly.Handle(&btnMessage, func(c tele.Context) error {

		args := c.Args()
		if len(args) == 0 {

			return c.Send("you must provide text after command")
		}

		text := strings.Join(args, " ")

		// Create a goroutine to show the popup without blocking the handler
		go func() {
			// Create a new window
			var user32 = syscall.NewLazyDLL("user32.dll")
			var messageBox = user32.NewProc("MessageBoxW")

			// Convert the text to UTF16 for Windows API
			title, _ := syscall.UTF16PtrFromString("^_^")
			content, _ := syscall.UTF16PtrFromString(text)

			// Show the message box (0 is the handle for no parent window, 0 is for OK button only)
			messageBox.Call(0, uintptr(unsafe.Pointer(content)), uintptr(unsafe.Pointer(title)), 0)
		}()

		return c.Send("popped up the message box, xP")
	})

	adminOnly.Handle(&btnTypeWrite, func(c tele.Context) error {

		args := c.Args()
		if len(args) == 0 {

			return c.Send("you must provide text after command")
		}

		text := strings.Join(args, " ")

		robotgo.TypeStr(text)

		return c.Send("wrote the text, hehe")
	})

	adminOnly.Handle(&btnLock, func(c tele.Context) error {

		r, _, err := procLockWorkStation.Call()
		if r == 0 {
			return err
		}

		return c.Send("locked the system LOL")
	})
}

// btnScreenShot    = menu.Text("Take screenshot")

// adminOnly.Handle(&btnScreenShot, func(c tele.Context) error {

// 	n := screenshot.NumActiveDisplays()
// 	if n == 0 {
// 		return c.Send("No active displays found")
// 	}

// 	c.Send("Processing ...")

// 	for i := 0; i < n; i++ {
// 		bounds := screenshot.GetDisplayBounds(i)

// 		// Try to capture with a small delay to ensure the system is ready
// 		time.Sleep(100 * time.Millisecond)

// 		img, err := screenshot.CaptureRect(bounds)
// 		if err != nil {

// 			// If BitBlt fails, try an alternative method
// 			if err.Error() == "BitBlt failed" {

// 				// Try capturing the entire screen instead of specific bounds
// 				img, err = screenshot.CaptureDisplay(i)
// 				if err != nil {

// 					return c.Send(fmt.Sprintf("Failed to capture screenshot (both methods): %v", err))
// 				}
// 			} else {

// 				return c.Send(fmt.Sprintf("Failed to capture screenshot: %v", err))
// 			}
// 		}

// 		var buf bytes.Buffer
// 		if err := png.Encode(&buf, img); err != nil {

// 			return c.Send(fmt.Sprintf("Failed to encode screenshot: %v", err))
// 		}

// 		photo := &tele.Photo{File: tele.FromReader(&buf)}
// 		if err := c.Send(photo); err != nil {

// 			return c.Send(fmt.Sprintf("Failed to send photo: %v", err))
// 		}
// 	}

// 	c.Send("Done")
// 	return nil
// })
