package apicontroller

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/gowin/types"
)

func (controller *Controller) AddJob(c *gin.Context) {

	clientID := c.Query("client_id")
	if clientID == "" {

		controller.error(c, errors.New("client_id required"))
		return
	}

	name := c.Query("name")
	if name == "" {

		controller.error(c, errors.New("job required"))
		return
	}

	params := c.Query("params")
	if params == "" {

		controller.error(c, errors.New("params required"))
		return
	}

	var paramsArray []string
	err := json.Unmarshal([]byte(params), &paramsArray)
	if err != nil {

		controller.error(c, err)
		return
	}

	// add job to server
	controller.data.AddJob(clientID, types.NewJob(clientID, types.JobNameFromString(name), paramsArray))

	controller.ok(c, nil)
}
