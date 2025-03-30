package srv

// func StartMain() {

// 	// Create a channel to receive OS signals
// 	sigs := make(chan os.Signal, 1)
// 	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

// 	waitChannel := make(chan struct{}, 1)
// 	go func() {

// 		// we can exit from app now
// 		defer func() {

// 			waitChannel <- struct{}{}
// 		}()

// 		// wait for signal to exit
// 		<-sigs
// 		logger.Debug("Application cancelled").Log()
// 	}()

// 	// start api
// 	a := api.NewApi()
// 	a.Start()

// 	// start telegram
// 	t := telegram.NewTelegram()
// 	t.Start()

// 	// wait to exit from app
// 	<-waitChannel

// 	logger.Debug("Application terminated").Log()
// }
