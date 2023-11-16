package main

import (
	"os"

	"github.com/CCOLLOT/appnametochange/cmd"
)

func main() {
	if err := cmd.InitAndRunCommand(); err != nil {
		os.Exit(3)
	}
}
