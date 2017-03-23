package testutils

import (
	"fmt"
)

// MockLogger is a mock implementation of the logger interface
type MockLogger struct {
	InfoCallCount        int
	ErrorCallCount       int
	ErrorStackTraceCount int
}

// Info is a mock version of the real logger
func (log *MockLogger) Info(data string) {
	fmt.Printf(data)
	log.InfoCallCount++
}

// Error is a mock version of the real logger
func (log *MockLogger) Error(data string) {
	fmt.Printf(data)
	log.ErrorCallCount++
}

// ErrorStackTrace is a mock version of the real logger
func (log *MockLogger) ErrorStackTrace(message string, data error) {
	fmt.Printf(message + data.Error())
	log.ErrorStackTraceCount++
}
