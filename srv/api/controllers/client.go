package apicontroller

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

func (controller *Controller) ClientPing(c *gin.Context) {

	clientID := c.Query("client_id")
	if clientID == "" {

		controller.error(c, errors.New("client_id required"))
		return
	}

	controller.data.UpdateClientLastSeen(clientID)

	controller.ok(c, nil)
}

func (controller *Controller) ClientPollUpdates(c *gin.Context) {

	clientID := c.Query("client_id")
	if clientID == "" {

		controller.error(c, errors.New("client_id required"))
		return
	}

	// Create timeout channel
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)

	defer ticker.Stop()

	for {
		// Check for jobs first
		jobs := controller.data.PullJobs(clientID)
		if len(jobs) > 0 {

			controller.ok(c, jobs)
			return
		}

		// Wait for either new jobs or timeout
		select {

		// Return empty response on timeout
		case <-timeout:
			controller.ok(c, []string{})
			return

		// each 500ms check for jobs
		case <-ticker.C:

			// Check for jobs first
			jobs := controller.data.PullJobs(clientID)
			if len(jobs) > 0 {

				controller.ok(c, jobs)
				return
			}

			continue

		// Client disconnected
		case <-c.Request.Context().Done():
			return
		}
	}
}

func (controller *Controller) ClientResultJob(c *gin.Context) {

	clientID := c.Query("client_id")
	if clientID == "" {

		controller.error(c, errors.New("client_id required"))
		return
	}

	jobID := c.Query("job_id")
	if jobID == "" {

		controller.error(c, errors.New("job_id required"))
		return
	}

	message := c.Query("message")
	if message == "" {

		controller.error(c, errors.New("message required"))
		return
	}

	controller.data.AddLog(jobID, message)

	controller.ok(c, nil)
}
