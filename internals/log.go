package internals

import (
	"fmt"
	"os"
)

func Fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Warn(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
