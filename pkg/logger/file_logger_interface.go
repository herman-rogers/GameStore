package logger

// Log is the base interface for implementing a logging system
type Log interface {
	Info(data string)
	Error(err string)
	ErrorStackTrace(message string, data error)
}
