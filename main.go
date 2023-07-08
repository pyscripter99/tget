/*
Copyright © 2023 Ryder Retzlaff <ryder@retzlaff.family>
*/
package main

import (
	"fmt"

	"github.com/pyscripter99/tget/cmd"
)

var version string = "UNSET"

func main() {
	fmt.Println("tget " + version)
	cmd.Execute()
}
