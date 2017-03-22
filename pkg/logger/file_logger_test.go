package logger

import (
	"testing"
	"time"

	"github.com/go-errors/errors"
)

type MockStorage struct {
	callCount int
	mockData  string
}

func (save *MockStorage) StoreData(data string) {
	save.callCount++
	save.mockData = data
}

func TestFileSystemInfoFormatIsCorrect(t *testing.T) {
	data := "TEST_MOCK_INFO"
	time := time.Now().Local()
	timestamp := time.Format(timeFormat)
	mockinfo := "SYSTEM_INFO " + timestamp + ": " + data

	mocklog := FileSystemLogs{}
	loginfo := mocklog.formatLogData(data, infoTag, time)

	if mockinfo != loginfo {
		t.Errorf("Expected %v, but got %v", mockinfo, loginfo)
	}
}

func TestFileSystemErrorFormatIsCorrect(t *testing.T) {
	data := "TEST_ERROR_INFO"
	time := time.Now().Local()
	timestamp := time.Format(timeFormat)
	mockerror := "SYSTEM_ERROR " + timestamp + ": " + data

	mockLog := FileSystemLogs{}
	logerror := mockLog.formatLogData(data, errorTag, time)

	if mockerror != logerror {
		t.Errorf("Expected %v, but got %v", mockerror, logerror)
	}
}

func TestFileSystemLogsInfo(t *testing.T) {
	data := "TEST_MOCK_INFO"
	mocksave := MockStorage{}
	mocklog := FileSystemLogs{Save: &mocksave}
	mockinfo := mocklog.formatLogData(data, infoTag, time.Now().Local())
	mocklog.Info(data)

	if mocksave.callCount != 1 {
		t.Errorf("Expected Info to be called once, but was called %v", mocksave.callCount)
	}
	if mockinfo != mocksave.mockData {
		t.Errorf("Expected %v, but got %v", mockinfo, mocksave.mockData)
	}
}

func TestFilesystemLogsErrors(t *testing.T) {
	err := "MOCK_ERROR"
	mocksave := MockStorage{}
	mocklog := FileSystemLogs{Save: &mocksave}
	mockerror := mocklog.formatLogData(err, errorTag, time.Now().Local())
	mocklog.Error(err)

	if mocksave.callCount != 1 {
		t.Errorf("Expected Error to be called once, but was called %v", mocksave.callCount)
	}
	if mockerror != mocksave.mockData {
		t.Errorf("Expected %v, but got %v", mockerror, mocksave.mockData)
	}
}

func TestFileSystemLogsErrorsWithStackTraces(t *testing.T) {
	err := errors.Errorf("TEST_MOCK_ERROR")
	mocksave := MockStorage{}
	mocklog := FileSystemLogs{Save: &mocksave}
	mockerror := "MOCK_TEST: " + mocklog.formatLogData(err.Error(), errorTag, time.Now().Local())
	mockerror += mocklog.errStackTrack(err)
	mocklog.ErrorStackTrace("MOCK_TEST: ", err)

	if mocksave.callCount != 1 {
		t.Errorf("Expected ErrorStackTrace to be called once, but was called %v", mocksave.callCount)
	}
	if mockerror != mocksave.mockData {
		t.Errorf("Expected %v, but got %v", mockerror, mocksave.mockData)
	}
}
