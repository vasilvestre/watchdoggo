package main

import (
	"os"
)

func Check(err error) {
	if err != nil {
		WriteLog(err.Error(), LogFatal)
		os.Exit(2)
	}
}