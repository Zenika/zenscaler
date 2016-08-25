package main

import (
	"os"

	"github.com/Zenika/zenscaler/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
