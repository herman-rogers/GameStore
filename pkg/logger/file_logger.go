package logger

import (
	"time"

	"github.com/WestCoastOpenSource/GameStore/pkg/storage"
	"github.com/go-errors/errors"
)

// FileSystemLogs will format and output logs into a specified file
type FileSystemLogs struct {
	Save storage.DataStorage
}

const infoTag = "SYSTEM_INFO"
const errorTag = "SYSTEM_ERROR"
const timeFormat = "Mon Jan _2 15:04:05 UTC 2006"

// Info logs general data and info about the operation of the system
func (logs FileSystemLogs) Info(data string) {
	systemInfo := logs.formatLogData(data, infoTag, time.Now().Local())
	logs.Save.StoreData(systemInfo)
}

// Error logs error data from the system
func (logs FileSystemLogs) Error(err string) {
	systemErr := logs.formatLogData(err, errorTag, time.Now().Local())
	logs.Save.StoreData(systemErr)
}

// ErrorStackTrace logs error data from the system and outputs a stack trace
func (logs FileSystemLogs) ErrorStackTrace(message string, err error) {
	systemErr := message + logs.formatLogData(err.Error(), errorTag, time.Now().Local())
	systemErr += logs.errStackTrack(err)

	logs.Save.StoreData(systemErr)
}

func (logs FileSystemLogs) errStackTrack(err error) string {
	return "\n SYSTEM_ERROR_STACKTRACE: " + err.(*errors.Error).ErrorStack()
}

func (logs FileSystemLogs) formatLogData(data string, tag string, time time.Time) string {
	timestamp := time.Format(timeFormat)
	formattedData := tag + " " + timestamp + ": " + data
	return formattedData
}
