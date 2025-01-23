package api

import apicontroller "github.com/ranjbar-dev/gowin/srv/api/controllers"

func (a *Api) registerRoutes() {

	controller := apicontroller.NewController()

	// register middlewares //
	a.registerMiddlewares()

	// server //

	a.hs.RegisterGetRoute("/server/ping", controller.ServerPing)

	a.hs.RegisterGetRoute("/server/timezone", controller.ServerTimezone)

	a.hs.RegisterGetRoute("/server/timezone-offset", controller.ServerTimezoneOffset)

	a.hs.RegisterGetRoute("/server/time", controller.ServerTime)

	a.hs.RegisterGetRoute("/server/memory-usage", controller.ServerMemoryUsage)

	// monitor //

	a.hs.RegisterGetRoute("/monitor/screen-shot", controller.MonitorScreenShot)

	a.hs.RegisterGetRoute("/monitor/processes", controller.MonitorProcesses)

	// action //

	a.hs.RegisterGetRoute("/action/lock", controller.ActionLock)

	a.hs.RegisterGetRoute("/action/shutdown", controller.ActionShutdown)

}
