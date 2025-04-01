package telegram

import (
	"bytes"
	"fmt"
	"image/png"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/kbinani/screenshot"
	"github.com/mitchellh/go-ps"
	"github.com/ranjbar-dev/gowin/config"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

var (
	user32              = syscall.MustLoadDLL("user32.dll")
	procLockWorkStation = user32.MustFindProc("LockWorkStation")
	procSendInput       = user32.MustFindProc("SendInput")

	// Universal markup builders.
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}

	// Reply buttons.
	btnHelp          = menu.Text("/help")
	btnTime          = menu.Text("/time")
	btnListProcesses = menu.Text("/processes")
	btnTerminate     = menu.Text("/terminate")
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
		text += "/terminate - terminate a process\n"
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

		var list string
		for _, process := range processes {

			list += fmt.Sprintf("%d: %s\n", process.Pid(), process.Executable())
		}

		// if list size is more than 4000 char, chunks it and send it in multiple messages
		if len(list) > 4000 {

			start := 0
			end := 4000
			for {

				if end > len(list) {

					end = len(list)
				}

				c.Send(list[start:end])
				start += 4000
				end += 4000

				if start >= len(list) {

					break
				}
			}
			return nil
		} else {

			return c.Send(list)
		}
	})

	adminOnly.Handle(&btnTerminate, func(c tele.Context) error {

		args := c.Args()
		if len(args) == 0 {

			return c.Send("you must provide a pid after command")
		}

		pid := args[0]
		pidInt, err := strconv.Atoi(pid)
		if err != nil {

			return c.Send("invalid pid")
		}

		// Terminate process using Windows API
		handle, err := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, uint32(pidInt))
		if err != nil {

			return c.Send("failed to open process")
		}
		defer syscall.CloseHandle(handle)

		err = syscall.TerminateProcess(handle, 1)
		if err != nil {

			return c.Send("failed to terminate process")
		}

		return c.Send("process terminated")
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

	adminOnly.Handle("/type", func(c tele.Context) error {
		args := c.Args()
		if len(args) == 0 {
			return c.Send("you must provide text after command")
		}

		text := strings.Join(args, " ")
		err := simulateKeyboardInput(text)
		if err != nil {
			return c.Send(fmt.Sprintf("Error typing text: %v", err))
		}

		return c.Send("Text typed successfully")
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

	menu = &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(
		menu.Row(btnHelp, btnTime),
		menu.Row(btnListProcesses, btnScreenShot),
		menu.Row(btnLock, btnShutdown),
	)

	_, err = t.bot.Send(tele.ChatID(config.TelegramChatID()), info, menu)
	if err != nil {

		return fmt.Errorf("Error sending application started message: %v", err)
	}

	return nil
}

func (t *Telegram) SendApplicationStoppedMessage() error {

	_, err := t.bot.Send(tele.ChatID(config.TelegramChatID()), "Application stopped")

	return err
}

func simulateKeyboardInput(text string) error {
	for _, char := range text {
		// Convert character to virtual key code
		keyCode := charToKeyCode(char)

		// Create input event for key down
		input := struct {
			Type uint32
			Ki   struct {
				Vk          uint16
				Scan        uint16
				Flags       uint32
				Time        uint32
				DwExtraInfo uintptr
			}
		}{
			Type: 1, // INPUT_KEYBOARD
			Ki: struct {
				Vk          uint16
				Scan        uint16
				Flags       uint32
				Time        uint32
				DwExtraInfo uintptr
			}{
				Vk:    keyCode,
				Flags: 0, // KEYEVENTF_KEYDOWN
			},
		}

		// Send key down
		ret, _, err := procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
		if ret == 0 {
			return fmt.Errorf("failed to send key down: %v", err)
		}

		// Create input event for key up
		input.Ki.Flags = 0x0002 // KEYEVENTF_KEYUP
		ret, _, err = procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
		if ret == 0 {
			return fmt.Errorf("failed to send key up: %v", err)
		}

		time.Sleep(10 * time.Millisecond) // Small delay between keystrokes
	}

	return nil
}

func charToKeyCode(char rune) uint16 {
	// Simple mapping for basic characters
	if char >= 'a' && char <= 'z' {
		return uint16(char - 'a' + 0x41) // VK_A through VK_Z
	}
	if char >= 'A' && char <= 'Z' {
		return uint16(char - 'A' + 0x41)
	}
	if char >= '0' && char <= '9' {
		return uint16(char - '0' + 0x30) // VK_0 through VK_9
	}

	// Space
	if char == ' ' {
		return 0x20 // VK_SPACE
	}

	// Enter
	if char == '\n' {
		return 0x0D // VK_RETURN
	}

	// Default to space if character not mapped
	return 0x20
}

// btnScreenShot    = menu.Text("Take screenshot")
