package apicontroller

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kbinani/screenshot"
	"github.com/mitchellh/go-ps"
)

func (controller *Controller) MonitorScreenShot(c *gin.Context) {

	n := screenshot.NumActiveDisplays()
	if n == 0 {
		controller.error(c, fmt.Errorf("no active displays found"))
		return
	}

	var allBounds image.Rectangle
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		allBounds = allBounds.Union(bounds)
	}

	allImg := image.NewRGBA(allBounds)
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			controller.error(c, err)
			return
		}
		draw.Draw(allImg, bounds, img, image.Point{}, draw.Src)
	}

	c.Header("Content-Type", "image/png")
	png.Encode(c.Writer, allImg)
}

func (controller *Controller) MonitorProcesses(c *gin.Context) {

	processes, err := ps.Processes()
	if err != nil {

		controller.error(c, err)
		return
	}

	var applications []string
	for _, process := range processes {

		applications = append(applications, process.Executable())
	}

	c.JSON(http.StatusOK, applications)
}
