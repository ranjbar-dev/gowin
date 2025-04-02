package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/ranjbar-dev/gowin/srv/client"
	"github.com/ranjbar-dev/gowin/tools/logger"
	"github.com/ranjbar-dev/gowin/tools/telegram"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	sigs    chan os.Signal
	guiChan chan struct{} = make(chan struct{})
	fyneApp fyne.App
)

func main() {
	// Initialize Fyne app once
	fyneApp = app.New()

	// create a channel to receive OS signals
	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	forever := make(chan struct{}, 1)
	go func() {
		<-sigs
		logger.Info("Application stopped").Log()
		telegram.SendMessage("Client stopped")
		if fyneApp != nil {
			fyneApp.Quit()
		}
		os.Exit(0)
		systray.Quit()
	}()

	c := client.NewClient("work", "http://localhost:3761")
	c.Start()

	// Start system tray in a goroutine
	go systray.Run(onReady, onExit)

	// Main event loop
	for {
		select {
		case <-forever:
			return
		case <-guiChan:
			showGUI()
		}
	}
}

func onReady() {
	logger.Info("Application ready").Log()

	systray.SetIcon(getIcon())
	systray.SetTitle("Gowin")
	systray.SetTooltip("Gowin")

	// register menu items
	openGUI := systray.AddMenuItem("Open GUI", "Open the GUI")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-openGUI.ClickedCh:
				// Signal main thread to show GUI
				guiChan <- struct{}{}
			}
		}
	}()

	telegram.SendMessage("Client started")
}

func onExit() {
	logger.Info("Application exit").Log()
	sigs <- syscall.SIGTERM
}

func getIcon() []byte {
	icon, err := os.ReadFile("assets/icon.ico")
	if err != nil {
		panic(err)
	}
	return icon
}

func showGUI() {
	fmt.Println("GUI opened")

	w := fyneApp.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	// Set close callback
	w.SetOnClosed(func() {
		fmt.Println("GUI closed")
	})
	w.Show()
	w.RequestFocus()
	w.ShowAndRun()

	w.Close()
}
