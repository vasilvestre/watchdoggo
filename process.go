package main

import (
	"time"
	"os/exec"
	"strings"
	"fmt"
	"strconv"
)

var uptime = 0

func LaunchWatchdog() {
	configuration = GetConfiguration()
	tick := time.Tick(time.Duration(configuration.RetryEvery) * time.Second)
	KeepAliveProcess()
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
		WriteLog("Process stopped. Attempt to reload.", LogWarning)
		uptime = 0
		LaunchProcess()
	} else {
		WriteLog(fmt.Sprintf("Process is running fine. %s seconds uptime.",strconv.Itoa(uptime)), LogInfo)
		uptime += configuration.RetryEvery
	}
}

func LaunchProcess(){
	var err error = nil
	method := configuration.Method
	switch configuration.Method {
	case "bin":
		err = exec.Command("bash", "-c", configuration.ProcessName).Start()
	case "systemctl":
		out, _ := exec.Command(method,"status",configuration.ProcessName).Output()
		if strings.Contains(string(out),"not-found") {
			WriteLog("Method used isn't compatible with process. Retry assuming it's a bin..", LogWarning)
			err = exec.Command("bash", "-c", configuration.ProcessName).Start()
		} else {
			err = exec.Command(method,"start",configuration.ProcessName).Run()
		}
	case "service":
		out, _ := exec.Command(method,configuration.ProcessName,"status").Output()
		if strings.Contains(string(out),"unrecognized") {
			WriteLog("Method used isn't compatible with process. Retry assuming it's a bin..", LogWarning)
			err = exec.Command("bash", "-c", configuration.ProcessName).Start()
		} else {
			err = exec.Command(method,configuration.ProcessName,"start").Run()
		}
	}
	Check(err)
	WriteLog("Process launched successfully.",LogWarning)
}

