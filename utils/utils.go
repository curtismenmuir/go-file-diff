package utils

import "fmt"

var log = fmt.Println

// Logger will print a string to console when verbose flag is set
// Verbose flag can be overwritten (true) to log to console
func Logger(message string, verbose bool) {
	if !verbose {
		return
	}

	_, _ = log(message)
}
