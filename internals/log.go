package internals

import (
	"fmt"
	"os"
)

// Fatal prints the error and exits with status 1
func Fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Warn prints the error
func Warn(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
