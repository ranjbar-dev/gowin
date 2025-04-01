package apicontroller

import (
	"github.com/gin-gonic/gin"
)

func (controller *Controller) ServerPing(c *gin.Context) {

	controller.ok(c, nil)
}
