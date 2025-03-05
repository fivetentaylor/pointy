package utils

import "time"

func UnixNanoToTime(unixNano int64) *time.Time {
	timestamp := unixNano / int64(time.Millisecond) // Convert nanoseconds to milliseconds
	seconds := timestamp / 1e3                      // Convert milliseconds to seconds
	nanoseconds := (timestamp % 1e3) * 1e6

	t := time.Unix(seconds, nanoseconds)
	return &t
}
