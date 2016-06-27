package cmd

import (
	"fmt"
	"zscaler/core"
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
		parseConfig()
		core.Initialize()
	},
}

func init() {
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
		Probes:   make(map[string]probe.Probe, 0),
	}

	// add some mocks probes
	p := config.Probes
	p["DefaultScalingProbe"] = new(probe.DefaultScalingProbe)
	// mock scaler
	var mockScaler = new(service.MockScaler)

	for _, key := range viper.Sub("services").AllKeys() {
		fmt.Println("Add service: " + key)
		config.Services = append(config.Services, service.Service{
			Name:  key,
			Scale: mockScaler,
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
