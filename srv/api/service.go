package api

import (
	"context"

	"github.com/ranjbar-dev/golog"
	"github.com/ranjbar-dev/gowin/config"
	"github.com/ranjbar-dev/gowin/internal/httpserver"
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

		err := a.hs.Serve()
		if err != nil {

			golog.Logger.Error("Api server stopped", "server stopped", err)
		} else {

			golog.Logger.Trace("Api server stopped", "server stopped")
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
