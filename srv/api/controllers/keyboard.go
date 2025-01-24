package apicontroller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-vgo/robotgo"
)

func (controller *Controller) KeyboardPress(c *gin.Context) {

	data, ok := controller.queries(c, map[string]string{
		"keys": "string",
	})
	if !ok {

		return
	}

	var keys []string
	err := json.Unmarshal([]byte(data["keys"].(string)), &keys)
	if err != nil {

		controller.error(c, err)
		return
	}

	for _, key := range keys {

		robotgo.KeyDown(key)
	}

	for _, key := range keys {

		robotgo.KeyUp(key)
	}

	controller.ok(c, "keys pressed")
}

func (controller *Controller) KeyboardType(c *gin.Context) {

	data, ok := controller.queries(c, map[string]string{
		"text": "string",
	})
	if !ok {

		return
	}

	robotgo.TypeStr(data["text"].(string))

	controller.ok(c, "text typed")
}
