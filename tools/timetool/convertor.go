package timetool

import (
	"time"
)

func TimestampMsToTime(timestamp int64) time.Time {

	seconds := timestamp / 1000

	nanoseconds := (timestamp % 1000) * 1e6

	return time.Unix(seconds, nanoseconds)
}

func TimestampMsToDatetimeZ(timestamp int64) string {

	t := TimestampMsToTime(timestamp).UTC()

	return TimeToDatetimeZ(t)
}

func TimeToDatetimeZ(t time.Time) string {

	return t.Format(time.RFC3339)
}
