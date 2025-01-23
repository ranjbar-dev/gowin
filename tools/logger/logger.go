package logger

import (
	"encoding/json"
	"fmt"

	"github.com/ranjbar-dev/gowin/tools/timetool"
)

type LogRecord struct {
	level   string
	title   string
	message string
	params  map[string]interface{}
}

func (l *LogRecord) Message(message string) *LogRecord {

	l.message = message
	return l
}

func (l *LogRecord) Params(params map[string]interface{}) *LogRecord {

	l.params = params
	return l
}

func (l *LogRecord) Log() {

	var log = fmt.Sprintf("[%s] %s - %s", l.level, timetool.TimeToDatetimeZ(timetool.Now()), l.title)

	if l.message != "" {

		log += "\n " + l.message
	}

	if l.params != nil {

		jsonData, err := json.Marshal(l.params)
		if err != nil {

			fmt.Println("Error marshalling log data:", err)
			return
		}

		log += "\n " + string(jsonData)
	}

	fmt.Println(log)
}

func Error(title string) *LogRecord {

	return &LogRecord{level: "ERROR", title: title}
}

func Warn(title string) *LogRecord {

	return &LogRecord{level: "WARNING", title: title}
}

func Info(title string) *LogRecord {

	return &LogRecord{level: "INFO", title: title}
}

func Debug(title string) *LogRecord {

	return &LogRecord{level: "DEBUG", title: title}
}
