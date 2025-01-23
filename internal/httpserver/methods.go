package httpserver

import (
	"context"

	"github.com/gin-gonic/gin"
)

func (hs *HttpServer) GetRouter() *gin.Engine {

	return hs.ge
}

func (hs *HttpServer) RegisterGetRoute(path string, callback func(c *gin.Context)) {

	hs.ge.GET(path, callback)
}

func (hs *HttpServer) RegisterPostRoute(path string, callback func(c *gin.Context)) {

	hs.ge.POST(path, callback)
}

func (hs *HttpServer) Shutdown(ctx context.Context) error {

	return hs.server.Shutdown(ctx)
}

func (hs *HttpServer) Serve() error {

	return hs.server.ListenAndServe()
}
