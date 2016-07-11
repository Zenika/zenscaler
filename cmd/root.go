// Package cmd provide the cli informations
package cmd

import (
	"errors"
	"fmt"
	"zscaler/core"
	"zscaler/core/rule"
	"zscaler/swarm"

	log "github.com/Sirupsen/logrus"
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
		log.Panicf("Fatal error in config file: %s \n", err)
	}

	// global configuration structure
	var config = &core.Config{
		Rules: make([]rule.Rule, 0),
	}
	// loop over the services
	// create one default rule by service
	rules := viper.Sub("rules")
	for _, r := range rules.AllKeys() {
		target := rules.Sub(r).GetString("target")
		log.Info("Add service [" + target + "]")

		// Parse rules
		up, err := rule.Decode(rules.Sub(r).GetString("up"))
		if err != nil {
			return nil, errors.New(target + fmt.Sprintf(": %v up", err))
		}
		down, err := rule.Decode(rules.Sub(r).GetString("down"))
		if err != nil {
			return nil, errors.New(target + fmt.Sprintf(": %v down", err))
		}
		config.Rules = append(config.Rules, rule.FloatValue{
			Scale:       rule.NewComposeScaler(target),
			Probe:       &swarm.AverageCPU{Tag: target},
			RefreshRate: rules.Sub(r).GetDuration("refresh"),
			Up:          up,
			Down:        down,
		})
	}
	log.Info("Configuration complete !")
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
