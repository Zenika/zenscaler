package main

import (
	"os"

	"github.com/Zenika/zscaler/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
