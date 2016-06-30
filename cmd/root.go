// Package cmd provide the cli informations
package cmd

import (
	"fmt"
	"zscaler/core"
	"zscaler/core/probe"
	"zscaler/core/rule"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
}

func parseConfig() (*core.Config, error) {
	// parse config file
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Sprintf("Fatal error config file: %s \n", err))
	}

	// global configuration structure
	var config = &core.Config{
		Rules: make([]rule.Rule, 0),
	}
	// loop over the services
	// create one default rule by service
	for _, name := range viper.Sub("services").AllKeys() {
		fmt.Println("Add service [" + name + "] using DefaultRule")
		config.Rules = append(config.Rules, rule.Default{
			Target: rule.MockService(name),
			Probe:  new(probe.DefaultScalingProbe),
		})
	}
	return config, nil
}

// DumpConfigCmd definition
var DumpConfigCmd = &cobra.Command{
	Use:   "dumpconfig",
	Short: "Dump parsed config file to stdout",
	Long:  `Check, parse and dump the configuration to the standart output`,
	Run: func(cmd *cobra.Command, args []string) {
		parseConfig()
	},
}
