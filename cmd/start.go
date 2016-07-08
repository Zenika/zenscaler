package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// StartCmd parse configuration file and launch the scaler
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start autoscaler",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := parseConfig()
		if err != nil {
			log.Panicf("Fatal error reading config file: %s \n", err)
		}
		config.Initialize()
	},
}
