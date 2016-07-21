// Package cmd provide the cli informations
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.3-alpha"

// Debug switch
var Debug bool

// RootCmd is the defaut command
var RootCmd = &cobra.Command{
	Use:   "zscaler",
	Short: "ZScaler is a simple yet flexible scaler",
	Long: `A Simple and Flexible scaler for various orchetrators.
Complete documentation is available at https://github.com/zenika/zscaler/wiki`,
	Run: func(cmd *cobra.Command, args []string) {
		help := cmd.HelpFunc()
		help(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(StartCmd)
	RootCmd.AddCommand(DumpConfigCmd)
	RootCmd.AddCommand(VersionCmd)
	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Activate debug output")
}

// VersionCmd display version number and informations
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version number",
	Long:  "Display version number and build informations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("zScaler %s\n", version)
		return
	},
}
