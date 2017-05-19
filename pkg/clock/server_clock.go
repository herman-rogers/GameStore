package clock

import "time"

// ServerClock is the basic implementation of server time and date
type ServerClock struct{}

// ClockTime returns the server time in UTC formatted to YYYY-MM-DD HH:MM:SS
func (ServerClock) ClockTime(serverTime time.Time) string {
	const layout = "06-01-02 03:04:05"
	return serverTime.Format(layout)
}
