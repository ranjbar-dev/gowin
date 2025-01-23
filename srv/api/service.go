package api

import (
	"context"

	"github.com/ranjbar-dev/gowin/config"
	"github.com/ranjbar-dev/gowin/internal/httpserver"
	"github.com/ranjbar-dev/gowin/tools/logger"
)

type Api struct {
	ctx    context.Context
	cancel context.CancelFunc
	hs     *httpserver.HttpServer
}

func (a *Api) Start() {

	a.registerRoutes()

	// start http server
	go func() {

		logger.Debug("Api server started").Log()
		err := a.hs.Serve()
		if err != nil {

			logger.Error("Api server stopped").Message(err.Error()).Log()
		} else {

			logger.Debug("Api server stopped").Log()
		}

	}()

	// shutdown server when context is done
	go func() {

		<-a.ctx.Done()

		a.hs.Shutdown(a.ctx)
	}()
}

func NewApi(ctx context.Context, cancel context.CancelFunc) *Api {

	return &Api{
		ctx:    ctx,
		cancel: cancel,
		hs:     httpserver.NewHttpServer(config.ApiHost(), config.ApiPort(), config.ApiDebug()),
	}
}
