package client

import (
	"fmt"
	"os/exec"

	"github.com/ranjbar-dev/gowin/tools/logger"
	"github.com/ranjbar-dev/gowin/types"
)

func (c *Client) handleJobs() {

	for {

		job := <-c.jobs

		switch job.Name {

		case types.JobNamePing:

			c.handlePing(job)

		case types.JobCopyClipboard:

			c.handleCopyClipboard(job)

		case types.JobLockScreen:

			c.handleLockScreen(job)

		default:
			logger.Error("unknown job: " + job.String()).Log()
		}

	}
}

// handle ping
func (c *Client) handlePing(job types.Job) {

	resp, err := c.request().SetQueryParams(map[string]string{
		"client_id": c.id,
	}).Get(c.host + "/client/ping")
	if err != nil {

		logger.Error("panic error").Message(err.Error()).Log()
		return
	}

	if resp.IsError() {

		logger.Error(resp.String()).Message(resp.String()).Log()
		return
	}
}

// handle copy clipboard
func (c *Client) handleCopyClipboard(job types.Job) {

	if len(job.Params) == 0 {

		logger.Error("no text to copy").Log()
		return
	}

	text := job.Params[0]

	// Copy text to clipboard
	cmd := exec.Command("cmd", "/c", fmt.Sprintf("echo %s", text), "|", "clip")
	err := cmd.Run()
	if err != nil {

		logger.Error("error copying text to clipboard").Message(err.Error()).Log()
		return
	}

	logger.Info("text copied to clipboard").Message(text).Log()
}

// handle lock screen
func (c *Client) handleLockScreen(job types.Job) {

	cmd := exec.Command("cmd", "/c", "rundll32.exe", "user32.dll,LockWorkStation")
	err := cmd.Run()
	if err != nil {

		logger.Error("error locking screen").Message(err.Error()).Log()
		return
	}

	logger.Info("screen locked").Log()
}
