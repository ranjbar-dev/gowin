package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/gowin/config"
)

func (a *Api) registerMiddlewares() {

	router := a.hs.GetRouter()

	router.Use(basicAuthMiddleware)
}

func basicAuthMiddleware(ctx *gin.Context) {

	if ctx.Request.URL.Path == "/server/ping" {

		ctx.Next()
		return
	}

	username, password, hasAuth := ctx.Request.BasicAuth()
	if !hasAuth || username != config.ApiBasicUsername() || password != config.ApiBasicPassword() {

		ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		ctx.AbortWithStatus(401)
		return
	}

	ctx.Next()
}
