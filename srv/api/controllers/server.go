package apicontroller

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/gowin/tools/timetool"
)

func (controller *Controller) ServerPing(c *gin.Context) {

	controller.ok(c, nil)
}

func (controller *Controller) ServerTime(c *gin.Context) {

	controller.ok(c, timetool.Now().Format("2006-01-02 15:04:05.000"))
}

func (controller *Controller) ServerTimezoneOffset(c *gin.Context) {

	_, offset := timetool.Now().Zone()

	controller.ok(c, offset)
}

func (controller *Controller) ServerTimezone(c *gin.Context) {

	controller.ok(c, timetool.Timezone())
}

func (controller *Controller) ServerMemoryUsage(c *gin.Context) {

	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	controller.ok(c, fmt.Sprintf("Memory Usage: %v MiB, Alloc: %v MiB, Sys: %v MiB, NumGC: %v MiB", uint64(m.Alloc/1024/1024), uint64(m.TotalAlloc/1024/1024), uint64(m.Sys/1024/1024), m.NumGC))
}
