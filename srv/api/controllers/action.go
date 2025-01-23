package apicontroller

import (
	"os/exec"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	user32              = syscall.MustLoadDLL("user32.dll")
	procLockWorkStation = user32.MustFindProc("LockWorkStation")
)

func lockWorkStation() error {

	r, _, err := procLockWorkStation.Call()
	if r == 0 {
		return err
	}
	return nil
}

func shutdownSystem() error {

	cmd := exec.Command("shutdown", "/s", "/t", "0")
	return cmd.Run()
}

func (controller *Controller) ActionLock(c *gin.Context) {

	if err := lockWorkStation(); err != nil {

		controller.error(c, err)
		return
	}

	controller.ok(c, "done")
}

func (controller *Controller) ActionShutdown(c *gin.Context) {

	if err := shutdownSystem(); err != nil {

		controller.error(c, err)
		return
	}

	controller.ok(c, "done")
}
