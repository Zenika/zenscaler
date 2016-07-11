package main

import (
	"os"
	"zscaler/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
