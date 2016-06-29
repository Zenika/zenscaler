package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// StartCmd parse configuration file and launch the scaler
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start autoscaler",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := parseConfig()
		if err != nil {
			panic(fmt.Sprintf("Fatal error config file: %s \n", err))
		}
		config.Initialize()
	},
}
