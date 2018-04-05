package main

import (
	"log"
	"strings"
	"time"
	"os/exec"
	"fmt"
	"github.com/tkanos/gonfig"
	"os"
	"strconv"
	"github.com/c2h5oh/datasize"
)

const maxSize = datasize.KB
var configuration = Configuration{}

func check(err error) {
	if err != nil {
		log.Fatalln(fmt.Sprintf("Fatal error : %s",err.Error()))
		os.Exit(2)
	}
}

func main() {
	logFile := createLog()
	launchWatchdog(logFile)
	logFile.Close()
}

func launchWatchdog(logFile *os.File) {
	configuration = getConfiguration()
	tick := time.Tick(time.Duration(configuration.RetryEvery) * time.Second)
	for {
		logFile = logRotate(logFile)
		configuration = checkConfiguration(configuration)
		select {
		case <-tick:
			keepAliveProcess()
		}
	}
}

func keepAliveProcess() {
	out, err := exec.Command("bash", "-c", "ps cax | grep -v grep | grep "+configuration.ProcessName+" | awk '{print $1}'").Output()
	check(err)
	pid := strings.Replace(string(out), " ", "", -1)
	if len(pid) < 1 {
		log.Println("Service coupé, redémarrage en cours")
		launchProcess()
	} else {
		log.Println("Service en cours")
	}
}

func launchProcess(){
	var err error = nil
	method := configuration.Method
	switch configuration.Method {
	case "bin":
		err = exec.Command("bash", "-c", configuration.ProcessName).Start()
	case "systemctl":
		out, _ := exec.Command(method,"status",configuration.ProcessName).Output()
		if strings.Contains(string(out),"not-found") {
			log.Println("Method used isn't compatible with process. Retry assuming it's a bin..")
			err = exec.Command("bash", "-c", configuration.ProcessName).Start()
		} else {
			err = exec.Command(method,"start",configuration.ProcessName).Run()
		}
	case "service":
		out, _ := exec.Command(method,configuration.ProcessName,"status").Output()
		if strings.Contains(string(out),"unrecognized") {
			log.Println("Method used isn't compatible with process. Retry assuming it's a bin..")
			err = exec.Command("bash", "-c", configuration.ProcessName).Start()
		} else {
			err = exec.Command(method,configuration.ProcessName,"start").Run()
		}
	}
	check(err)
}

func checkConfiguration(configuration Configuration) Configuration {
	if configuration != getConfiguration() {
		message := "Changement de la configuration : "
		if configuration.ProcessName != getConfiguration().ProcessName {
			message += "le processus à surveiller est devenu " + getConfiguration().ProcessName
		}
		if configuration.RetryEvery != getConfiguration().RetryEvery {
			message += "la fréquence de rafraichissement est devenu " + strconv.Itoa(getConfiguration().RetryEvery)
		}
		if configuration.Method != getConfiguration().Method {
			message += "la méthode de démarrage est devenu " + getConfiguration().Method
		}
		log.Println(message)
	}
	return getConfiguration()
}

func getConfiguration() Configuration {
	configuration := Configuration{}
	err := gonfig.GetConf("watchdog-go.json", &configuration)
	check(err)
	return configuration
}

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

type Configuration struct {
	ProcessName string
	RetryEvery  int
	Method string
}