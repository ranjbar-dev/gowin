package data

import (
	"fmt"

	"github.com/ranjbar-dev/gowin/tools/logger"
)

// get logs
func (d *Data) GetLogs() []string {

	d.logsMutex.Lock()
	defer d.logsMutex.Unlock()

	return d.logs
}

// add log
func (d *Data) AddLog(id string, message string) {

	d.logsMutex.Lock()
	defer d.logsMutex.Unlock()

	// remove old result if more than 1000
	if len(d.logs) > 1000 {

		d.logs = d.logs[1:]
	}

	// add new result
	d.logs = append(d.logs, fmt.Sprintf("%s: %s", id, message))

	// log message
	logger.Info(fmt.Sprintf("Job: %s", id)).Message(message).Log()

	// TODO : send to telegram
}
