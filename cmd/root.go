package cmd

import (
	"fmt"
	"zscaler/probe"
	"zscaler/service"

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

func parseConfig() (*service.Config, error) {
	// parse config file
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var config = &service.Config{
		Services: make([]service.Service, 0),
		Probes:   probe.Initialize(),
	}

	// mock scaler
	var mockScaler = new(service.MockScaler)

	for _, key := range viper.Sub("services").AllKeys() {
		fmt.Println("Add service: " + key)
		config.Services = append(config.Services, service.Service{
			Name:  key,
			Scale: mockScaler,
			Rule: func() bool {
				return p["DefaultScalingProbe"].Value() > 0.5
			},
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
