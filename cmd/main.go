package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/lxn/win"
	"github.com/ranjbar-dev/gowin/srv/api"
	"github.com/ranjbar-dev/gowin/srv/telegram"
	"github.com/ranjbar-dev/gowin/tools/logger"
)

var telegramService *telegram.Telegram
var apiService *api.Api
var sigs chan os.Signal

func main() {

	// hide console window
	console := win.GetConsoleWindow()
	if console != 0 {

		win.ShowWindow(console, win.SW_HIDE)
	}

	// create a channel to receive OS signals
	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	waitChannel := make(chan struct{}, 1)
	go func() {

		// wait for signal to exit
		<-sigs
		logger.Info("Application stopped").Log()

		// send application stopped message to telegram
		err := telegramService.SendApplicationStartedMessage()
		if err != nil {

			logger.Error("Failed to send application started message to telegram").Message(err.Error()).Log()
		}

		// exit the app
		waitChannel <- struct{}{}
	}()

	// start api
	apiService = api.NewApi()
	apiService.Start()

	// start telegram
	telegramService = telegram.NewTelegram()
	telegramService.Start()

	// sStart system tray
	systray.Run(onReady, onExit)

	// wait to exit from app
	<-waitChannel

	logger.Info("Application quit").Log()

	os.Exit(0)

	systray.Quit()
}

func onReady() {

	logger.Info("Application ready").Log()

	systray.SetIcon(getIcon())
	systray.SetTitle("Gowin")
	systray.SetTooltip("Gowin")

	// register menu items
	mQuit := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	// send application started message to telegram
	err := telegramService.SendApplicationStartedMessage()
	if err != nil {

		logger.Error("Failed to send application started message to telegram").Message(err.Error()).Log()
	}
}

func onExit() {

	logger.Info("Application exit").Log()

	// send terminate signal to stop services
	sigs <- syscall.SIGTERM
}

func getIcon() []byte {

	icon, err := os.ReadFile("assets/icon.ico")
	if err != nil {

		panic(err)
	}

	return icon
}
