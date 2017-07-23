// Package cmd provide the cli informations
package cmd

import (
	"fmt"

	"github.com/Zenika/zenscaler/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is the defaut command
var RootCmd = &cobra.Command{
	Use:   "zenscaler",
	Short: "ZScaler is a simple yet flexible scaler",
	Long: `A Simple and Flexible scaler for various orchetrators.
Complete documentation is available at https://github.com/zenika/zenscaler/wiki`,
	Run: func(cmd *cobra.Command, args []string) {
		help := cmd.HelpFunc()
		help(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(dumpConfigCmd)
	RootCmd.AddCommand(versionCmd)
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "Activate debug output")
	_ = viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	startCmd.Flags().StringP("api-port", "l", ":3000", "API listening address and port")
	_ = viper.BindPFlag("api-port", startCmd.Flags().Lookup("api-port"))
	startCmd.Flags().Bool("allow-cmd-probe", false, "Allow the use of cmd probe over the API")
	_ = viper.BindPFlag("allow-cmd-probe", startCmd.Flags().Lookup("allow-cmd-probe"))
}

// VersionCmd display version number and informations
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version number",
	Long:  "Display version number and build informations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("zenscaler %s, build with %s\n", core.Version, core.GoVersion)
	},
}
