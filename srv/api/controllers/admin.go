package apicontroller

import "github.com/gin-gonic/gin"

func (controller *Controller) AdminDashboard(c *gin.Context) {

	c.File("api/index.html")
}
