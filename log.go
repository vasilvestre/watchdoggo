package main

import (
	"os"
	"fmt"
	"log"
	"strings"
	"time"
)

const LogInfo = "INFO"
const LogWarning = "WARNING"
const LogFatal = "FATAL"

var logFile = CreateLog()
var logLevel = make(map[string]int)

func CreateLog() *os.File {
	if _, err := os.Stat("./log/"); os.IsNotExist(err) {
		os.Mkdir("./log/", 0777)
	}
	fileName := GetLogName()
	filePathAndName := fmt.Sprintf("./log/%s", fileName)
	file, _ := os.OpenFile(filePathAndName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	log.SetOutput(file)
	return file
}

func LogRotate(logFile *os.File) *os.File {
	fileInfo, err := logFile.Stat()
	Check(err)
	counterHistory := 1
	if fileInfo.Size() > int64(maxSize) {
		fileName := strings.Replace(fileInfo.Name(), ".log", "", -1)
		os.Rename(fmt.Sprintf("./log/%s", fileInfo.Name()), fmt.Sprintf("./log/%s.%d.log", fileName, counterHistory))
		logFile.Close()
		counterHistory++
		return CreateLog()
	}
	return logFile
}

func GetLogName() string {
	location,_ := time.LoadLocation("Europe/Rome")
	actualDateFormat := "2006_01_02_15_04"
	actualDate := time.Now().In(location).Format(actualDateFormat)
	return fmt.Sprintf("history_%s.log", actualDate)
}

func WriteLog(message string, level string) {
	if CanBeLogged(level) {
		location,_ := time.LoadLocation("Europe/Rome")
		log.Println(fmt.Sprintf("%s : %s",level, message))
		line := LogLine{}
		line.Level = level
		line.Message = message
		actualDateFormat := "2006_01_02_15_04_05"
		actualDate := time.Now().In(location).Format(actualDateFormat)
		db := getDatabase()
		err := db.Write("log_line",fmt.Sprintf("line_%s", actualDate),line)
		Check(err)
	}
}

func CanBeLogged(level string) bool {
	minimumLevel := configuration.MinimalLogLevel
	logLevel[LogInfo] = 1
	logLevel[LogWarning] = 2
	logLevel[LogFatal] = 3
	if logLevel[level] >= logLevel[minimumLevel] {
		return true
	}
	return false
}

type LogLine struct {
	Message string
	Level  string
}