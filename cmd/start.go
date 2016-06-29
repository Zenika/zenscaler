package cmd

import (
	"fmt"
	"zscaler/core"

	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start autoscaler",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := parseConfig()
		if err != nil {
			fmt.Errorf("Error: bad configuration file, %s\n", err)
		}
		core.Initialize(config)
	},
}
