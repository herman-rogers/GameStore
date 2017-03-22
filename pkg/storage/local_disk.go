package storage

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LocalDisk struct {
	File      string
	Directory string
}

func (save *LocalDisk) StoreData(data string) {
	logfile, err := save.getLogFile()
	if err != nil {
		fmt.Printf("File Save Error: %v", err.Error())
		return
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Print(data)
}

func (save *LocalDisk) getLogFile() (*os.File, error) {
	directory := save.Directory
	filename := save.defaultFile()
	path := directory + filename
	logfile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return nil, err
	}
	return logfile, nil
}

func (save *LocalDisk) defaultFile() string {
	time := time.Now().Local()
	extension := time.Format("2006-01-02" + ".log")
	return save.File + extension
}
