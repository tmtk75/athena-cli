package main

import (
	"fmt"
	"os"
)

func init() {
	logger.Printf = func(format string, v ...interface{}) {}
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var logger struct {
	Printf func(format string, v ...interface{})
}
