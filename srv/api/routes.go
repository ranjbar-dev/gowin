package api

import apicontroller "github.com/ranjbar-dev/gowin/srv/api/controllers"

func (a *Api) registerRoutes() {

	controller := apicontroller.NewController()

	// register middlewares //
	a.registerMiddlewares()

	// server //

	a.hs.RegisterGetRoute("/server/ping", controller.ServerPing)

	// client //

	a.hs.RegisterGetRoute("/client/ping", controller.ClientPing)

	a.hs.RegisterGetRoute("/client/poll-updates", controller.ClientPollUpdates)

	a.hs.RegisterGetRoute("/client/result-job", controller.ClientResultJob)

	a.hs.RegisterGetRoute("/client/add-job", controller.AddJob)

	// result //

	a.hs.RegisterGetRoute("/result/latest", controller.GetLatestResults)
}
