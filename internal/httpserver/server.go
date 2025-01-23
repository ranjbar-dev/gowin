package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	host   string
	port   string
	debug  bool
	server *http.Server
	ge     *gin.Engine
}

func NewHttpServer(host string, port string, debug bool) *HttpServer {

	if !debug {

		gin.SetMode(gin.ReleaseMode)
	}

	ge := gin.New()

	ge.Use(gin.Recovery())

	ge.SetTrustedProxies(nil)

	return &HttpServer{
		host:  host,
		port:  port,
		debug: debug,
		server: &http.Server{
			Addr:    host + ":" + port,
			Handler: ge,
		},
		ge: ge,
	}
}
