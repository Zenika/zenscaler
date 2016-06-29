package cmd

import (
	"fmt"
	"os"
	"zscaler/core"

	"github.com/spf13/cobra"
)

// StartCmd parse configuration file and launch the scaler
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start autoscaler",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := parseConfig()
		if err != nil {
			os.Exit(fmt.Errorf("Error: bad configuration file, %s\n", err))

		}
		core.Initialize(config)
	},
}
