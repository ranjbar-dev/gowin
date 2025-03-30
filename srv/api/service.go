package api

import (
	"fmt"

	"github.com/ranjbar-dev/gowin/config"
	"github.com/ranjbar-dev/gowin/internal/httpserver"
	"github.com/ranjbar-dev/gowin/tools/logger"
)

type Api struct {
	hs *httpserver.HttpServer
}

func (a *Api) Start() {

	a.registerRoutes()

	// start http server
	go func() {

		logger.Debug(fmt.Sprintf("Api server started http://%s:%s", config.ApiHost(), config.ApiPort())).Log()
		err := a.hs.Serve()
		if err != nil {

			logger.Error("Api server stopped").Message(err.Error()).Log()
		} else {

			logger.Debug("Api server stopped").Log()
		}

	}()

}

func NewApi() *Api {

	return &Api{
		hs: httpserver.NewHttpServer(config.ApiHost(), config.ApiPort(), config.ApiDebug()),
	}
}
