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
		tick = time.Tick(time.Duration(configuration.RetryEvery) * time.Second)
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
		log.Println("Process stopped. Attempt to reload.")
		launchProcess()
	} else {
		log.Println("Process is running fine.")
	}
}

func launchProcess(){
	var err error = nil
	method := configuration.Method
	switch configuration.Method {
	case "bin":
		err = exec.Command("bash", "-c", configuration.ProcessName).Start()
	case "systemctl":
		out, err := exec.Command(method,"status",configuration.ProcessName).Output()
		check(err)
		if strings.Contains(string(out),"not-found") {
			log.Println("Method used isn't compatible with process. Retry assuming it's a bin..")
			err = exec.Command("bash", "-c", configuration.ProcessName).Start()
		} else {
			err = exec.Command(method,"start",configuration.ProcessName).Run()
		}
	case "service":
		out, err := exec.Command(method,configuration.ProcessName,"status").Output()
		check(err)
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
		message := "Configuration changed : "
		if configuration.ProcessName != getConfiguration().ProcessName {
			message += fmt.Sprintf("process to watch changed for %s",getConfiguration().ProcessName)
		}
		if configuration.RetryEvery != getConfiguration().RetryEvery {
			message += fmt.Sprintf("refresh frequency changed for %s",strconv.Itoa(getConfiguration().RetryEvery))
		}
		if configuration.Method != getConfiguration().Method {
			message += fmt.Sprintf("starting method changed for %s",getConfiguration().Method)
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

type Configuration struct {
	ProcessName string
	RetryEvery  int
	Method string
}