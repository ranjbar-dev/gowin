package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/lxn/win"
	"github.com/ranjbar-dev/gowin/srv/client"
	"github.com/ranjbar-dev/gowin/tools/logger"
	"github.com/ranjbar-dev/gowin/tools/telegram"
)

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

	forever := make(chan struct{}, 1)
	go func() {

		// wait for signal to exit
		<-sigs
		logger.Info("Application stopped").Log()

		// send application stopped message to telegram
		telegram.SendMessage("Client stopped")

		os.Exit(0)

		systray.Quit()
	}()

	c := client.NewClient("work", "http://localhost:3761")
	c.Start()

	// Start system tray
	systray.Run(onReady, onExit)

	// wait to exit from app
	<-forever
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
	telegram.SendMessage("Client started")
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
