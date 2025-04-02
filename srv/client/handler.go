package client

import (
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/ranjbar-dev/gowin/tools/logger"
	"github.com/ranjbar-dev/gowin/types"
)

var (
	user32       = syscall.MustLoadDLL("user32.dll")
	getCursorPos = user32.MustFindProc("GetCursorPos")
	setCursorPos = user32.MustFindProc("SetCursorPos")
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

		case types.JobMoveMouse:

			c.handleUpdateMousePosition(job)

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

// update mouse position
func (c *Client) handleUpdateMousePosition(job types.Job) {

	if len(job.Params) == 0 {

		logger.Error("no mouse position to update").Log()
		return
	}

	// Get current cursor position
	var point struct {
		X, Y int32
	}
	_, _, err := getCursorPos.Call(uintptr(unsafe.Pointer(&point)))
	if err != nil && err.Error() != "The operation completed successfully." {
		return
	}

	// Parse relative movement from job
	x, err := strconv.Atoi(job.Params[0])
	if err != nil {

		return
	}
	y, err := strconv.Atoi(job.Params[1])
	if err != nil {

		return
	}

	// Calculate new position
	newX := point.X + int32(x)
	newY := point.Y + int32(y)

	// Move mouse to new position
	_, _, err = setCursorPos.Call(uintptr(newX), uintptr(newY))
	if err != nil {

		logger.Error("error moving mouse to new position").Message(err.Error()).Log()
		return
	}

	logger.Info("mouse moved to new position").Message(fmt.Sprintf("x: %d, y: %d", newX, newY)).Log()
}
