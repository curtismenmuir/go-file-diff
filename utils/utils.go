package utils

import "fmt"

var log = fmt.Println

func Logger(message string, verbose bool) {
	if !verbose {
		return
	}

	log(message)
}
