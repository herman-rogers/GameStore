package clock

import "time"

// TimeInfo is the interface to implement getting the system time, date, etc
type TimeInfo interface {
	ClockTime(serverTime time.Time) string
}
