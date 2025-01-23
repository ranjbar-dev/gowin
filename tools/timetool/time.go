package timetool

import (
	"time"

	"github.com/ranjbar-dev/gowin/config"
)

var appLocation *time.Location

func init() {

	loc, err := time.LoadLocation(config.Timezone())
	if err != nil {
		panic(err)
	}

	appLocation = loc
}

func Now() time.Time {

	return time.Now().In(appLocation)
}

func Timezone() string {

	return appLocation.String()
}

func ParseInLocation(layout, value string) (time.Time, error) {

	return time.ParseInLocation(layout, value, appLocation)
}

func Date() string {

	return Now().Format("2006-01-02")
}
