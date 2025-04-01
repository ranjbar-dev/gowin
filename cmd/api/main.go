package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ranjbar-dev/gowin/srv/api"
	"github.com/ranjbar-dev/gowin/tools/logger"
)

func main() {

	// create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	forever := make(chan struct{}, 1)
	go func() {

		// wait for signal to exit
		<-sigs
		logger.Info("Application stopped").Log()

		os.Exit(0)
	}()

	// start api
	apiService := api.NewApi()
	apiService.Start()

	// wait to exit from app
	<-forever
}
