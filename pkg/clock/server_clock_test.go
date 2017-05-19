package clock

import (
	"testing"
	"time"
)

func TestServerClockReturnsCorrectTime(t *testing.T) {
	clock := ServerClock{}
	mockTime := time.Date(2006, 12, 25, 12, 30, 01, 0, time.UTC)
	expectedTime := "06-12-25 12:30:01"

	time := clock.ClockTime(mockTime)
	if time == "" {
		t.Errorf("Expected time to not be nil")
	}
	if time != expectedTime {
		t.Errorf("Expected time format to be %v, got %v", expectedTime, time)
	}
}
