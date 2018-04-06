package main

import (
	"os"
	"fmt"
	"log"
	"strings"
	"time"
)

func createLog() *os.File {
	if _, err := os.Stat("./log/"); os.IsNotExist(err) {
		os.Mkdir("./log/", 0777)
	}
	fileName := getLogName()
	filePathAndName := fmt.Sprintf("./log/%s", fileName)
	file, _ := os.OpenFile(filePathAndName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	log.SetOutput(file)
	return file
}

func logRotate(logFile *os.File) *os.File {
	fileInfo, err := logFile.Stat()
	check(err)
	counterHistory := 1
	if fileInfo.Size() > int64(maxSize) {
		fileName := strings.Replace(fileInfo.Name(), ".log", "", -1)
		os.Rename(fmt.Sprintf("./log/%s", fileInfo.Name()), fmt.Sprintf("./log/%s.%d.log", fileName, counterHistory))
		logFile.Close()
		counterHistory++
		return createLog()
	}
	return logFile
}

func getLogName() string {
	actualDateFormat := "2006_01_02"
	actualDate := time.Now().UTC().Format(actualDateFormat)
	return fmt.Sprintf("history_%s.log", actualDate)
}