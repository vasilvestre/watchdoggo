package main

import (
	"time"
	"os/exec"
	"strings"
	"log"
	"fmt"
)

func LaunchWatchdog() {
	configuration = GetConfiguration()
	tick := time.Tick(time.Duration(configuration.RetryEvery) * time.Second)
	for {
		logFile = LogRotate(logFile)
		configuration = CheckConfiguration(configuration)
		tick = time.Tick(time.Duration(configuration.RetryEvery) * time.Second)
		select {
		case <-tick:
			KeepAliveProcess()
		}
	}
}

func KeepAliveProcess() {
	out, err := exec.Command("bash", "-c", "ps cax | grep -v grep | grep "+configuration.ProcessName+" | awk '{print $1}'").Output()
	Check(err)
	pid := strings.Replace(string(out), " ", "", -1)
	if len(pid) < 1 {
		log.Println("Process stopped. Attempt to reload.")
		LaunchProcess()
	} else {
		log.Println("Process is running fine.")
	}
}

func LaunchProcess(){
	var err error = nil
	method := configuration.Method
	switch configuration.Method {
	case "bin":
		fmt.Println("slt")
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
	Check(err)
}

