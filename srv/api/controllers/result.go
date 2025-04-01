package apicontroller

import (
	"github.com/gin-gonic/gin"
)

func (controller *Controller) GetLatestResults(c *gin.Context) {

	results := controller.data.GetLogs()

	controller.ok(c, results)
}
