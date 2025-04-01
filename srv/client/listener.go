package client

import (
	"encoding/json"
	"time"

	"github.com/ranjbar-dev/gowin/tools/logger"
	"github.com/ranjbar-dev/gowin/types"
)

func (c *Client) listenForJobs() {

	for {

		// poll jobs
		resp, err := c.request().SetQueryParams(map[string]string{
			"client_id": c.id,
		}).Get(c.host + "/client/poll-updates")
		if err != nil {

			logger.Error("panic error").Message(err.Error()).Log()
			time.Sleep(time.Second * 5)
			continue
		}

		// check if error
		if resp.IsError() {

			logger.Error(resp.String()).Message(resp.String()).Log()
			time.Sleep(time.Second * 5)
			continue
		}

		// parse jobs
		var jobs []types.Job
		err = json.Unmarshal(resp.Body(), &jobs)
		if err != nil {

			logger.Error(err.Error()).Log()
			continue
		}

		// push jobs to channel
		for _, job := range jobs {

			c.jobs <- job
		}
	}
}
