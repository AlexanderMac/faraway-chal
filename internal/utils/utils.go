package utils

import (
	"fmt"
	"os"
)

func CheckError(err error) {
	CheckErrorWithMessage(err, "")
}

func CheckErrorWithMessage(err error, message string) {
	if message == "" {
		message = "Fatal error"
	}
	if err != nil {
		fmt.Printf("%s: %s\n", message, err.Error())
		os.Exit(1)
	}
}
