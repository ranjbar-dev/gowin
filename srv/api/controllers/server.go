package apicontroller

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/gowin/tools/timetool"
)

// @Summary       Ping
// @Description   ping
// @Tags          Server
// @Accept        json
// @Produce       json
// @Router        /server/ping [get, post]
// @Success       200     {object}   nil   "Success"
func (controller *Controller) ServerPing(c *gin.Context) {

	controller.ok(c, nil)
}

// @Summary       Server time
// @Description   Get server time
// @Tags          Server
// @Accept        json
// @Produce       json
// @Router        /server/time [get]
// @Success       200     {string}   string   "Success"
func (controller *Controller) ServerTime(c *gin.Context) {

	controller.ok(c, timetool.Now().Format("2006-01-02 15:04:05.000"))
}

// @Summary       Server timezone offset
// @Description   Get server timezone offset
// @Tags          Server
// @Accept        json
// @Produce       json
// @Router        /server/timezone-offset [get]
// @Success       200     {string}   string   "Success"
func (controller *Controller) ServerTimezoneOffset(c *gin.Context) {

	_, offset := timetool.Now().Zone()

	controller.ok(c, offset)
}

// @Summary       Server timezone
// @Description   Get server timezone
// @Tags          Server
// @Accept        json
// @Produce       json
// @Router        /server/timezone [get]
// @Success       200     {string}   string   "Success"
func (controller *Controller) ServerTimezone(c *gin.Context) {

	controller.ok(c, timetool.Timezone())
}

// @Summary       Server memory usage
// @Description   Get server memory usage
// @Tags          Server
// @Accept        json
// @Produce       json
// @Router        /server/memory-usage [get]
// @Success       200     {string}   string   "Success"
func (controller *Controller) ServerMemoryUsage(c *gin.Context) {

	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	controller.ok(c, fmt.Sprintf("Memory Usage: %v MiB, Alloc: %v MiB, Sys: %v MiB, NumGC: %v MiB", uint64(m.Alloc/1024/1024), uint64(m.TotalAlloc/1024/1024), uint64(m.Sys/1024/1024), m.NumGC))
}
