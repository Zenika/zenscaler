package main

import (
	"zscaler/cmd"

	log "github.com/Sirupsen/logrus"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("%s", err)
	}
}
