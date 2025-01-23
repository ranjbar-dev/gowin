package apicontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vgo/robotgo"
)

func (controller *Controller) MouseMove(c *gin.Context) {

	data, ok := controller.queries(c, map[string]string{
		"x": "int",
		"y": "int",
	})
	if !ok {

		return
	}

	// TODO : not working correctly
	x := data["x"].(int)
	y := data["y"].(int)
	currentX, currentY := robotgo.Location()
	robotgo.Move(currentX+x, currentY+y, 3)

	controller.ok(c, nil)
}
