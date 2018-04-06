package main

import (
	"fmt"
	"strconv"
	"github.com/tkanos/gonfig"
	"github.com/c2h5oh/datasize"
)

const maxSize = 10 * datasize.KB
var configuration = Configuration{}

func CheckConfiguration(configuration Configuration) Configuration {
	if configuration != GetConfiguration() {
		message := "Configuration changed : "
		if configuration.ProcessName != GetConfiguration().ProcessName {
			message += fmt.Sprintf("process to watch changed for %s",GetConfiguration().ProcessName)
		}
		if configuration.RetryEvery != GetConfiguration().RetryEvery {
			message += fmt.Sprintf("refresh frequency changed for %s",strconv.Itoa(GetConfiguration().RetryEvery))
		}
		if configuration.Method != GetConfiguration().Method {
			message += fmt.Sprintf("starting method changed for %s",GetConfiguration().Method)
		}
		WriteLog(message,LogWarning)
	}
	return GetConfiguration()
}

func GetConfiguration() Configuration {
	configuration := Configuration{}
	err := gonfig.GetConf("configuration.json", &configuration)
	Check(err)
	return configuration
}

type Configuration struct {
	ProcessName string
	RetryEvery  int
	Method string
	MinimalLogLevel string
}