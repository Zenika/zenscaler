// Package cmd provide the cli informations
package cmd

import "github.com/spf13/cobra"

// Debug switch
var Debug bool

// RootCmd is the defaut command
var RootCmd = &cobra.Command{
	Use:   "zscaler",
	Short: "ZScaler is a simple yet flexible scaler",
	Long: `A Simple and Flexible scaler for various orchetrators.
Complete documentation is available at https://github.com/zenika/zscaler/wiki`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(StartCmd)
	RootCmd.AddCommand(DumpConfigCmd)
	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Activate debug output")
}
