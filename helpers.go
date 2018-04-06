package main

import (
	"log"
	"fmt"
	"os"
)

func Check(err error) {
	if err != nil {
		log.Fatalln(fmt.Sprintf("Fatal error : %s",err.Error()))
		os.Exit(2)
	}
}