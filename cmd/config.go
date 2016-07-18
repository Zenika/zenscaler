package cmd

import (
	"errors"
	"fmt"
	"strings"
	"zscaler/core"
	"zscaler/core/probe"
	"zscaler/core/rule"
	"zscaler/core/scaler"
	"zscaler/swarm"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DumpConfigCmd definition
var DumpConfigCmd = &cobra.Command{
	Use:   "dumpconfig",
	Short: "Dump parsed config file to stdout",
	Long:  `Check, parse and dump the configuration to the standart output`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := parseConfig()
		if err != nil {
			log.Fatalf("Error in config file: %s", err)
		}

	},
}

func parseConfig() (*core.Config, error) {
	// parse config file
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		return nil, fmt.Errorf("Cannot read config file: %s \n", err)
	}

	// global configuration structure
	var config = &core.Config{
		Scalers: make(map[string]scaler.Scaler, 5),
		Rules:   make([]rule.Rule, 0),
	}
	// Parse scalers
	var scalers = viper.Sub("scalers")
	for _, name := range scalers.AllKeys() {
		log.Info("Add scaler [" + name + "]")
		var s = scalers.Sub(name)
		switch s.GetString("type") {
		case "docker-compose":
			if s.GetString("config") == "" {
				return nil, errors.New("No config specified for docker-compose scaler [" + name + "]")
			}
			if s.GetString("target") == "" {
				return nil, errors.New("No target specified for docker-compose scaler [" + name + "]")
			}
			config.Scalers[name] = scaler.NewComposeScaler(s.GetString("target"), s.GetString("config"))
		case "":
			return nil, fmt.Errorf("Unknown scaler type: %s\n", s.GetString("type"))
		}
	}

	// loop over the services
	rules := viper.Sub("rules")
	for _, r := range rules.AllKeys() {
		target := rules.Sub(r).GetString("target")
		log.Info("Add service [" + target + "]")

		// Check scaler
		// TODO external function
		scaler := rules.Sub(r).GetString("scaler")
		if config.Scalers[scaler] == nil {
			return nil, errors.New("No scaler specified for rule [" + r + "]")
		}

		// Parse rules
		up, err := rule.Decode(rules.Sub(r).GetString("up"))
		if err != nil {
			return nil, errors.New(target + fmt.Sprintf(": %v up", err))
		}
		down, err := rule.Decode(rules.Sub(r).GetString("down"))
		if err != nil {
			return nil, errors.New(target + fmt.Sprintf(": %v down", err))
		}

		// Pick probe
		var p probe.Probe
		refProbe := rules.Sub(r).GetString("probe")
		splittedProbe := strings.Split(refProbe, ".")
		if len(splittedProbe) < 2 {
			return nil, errors.New("Badly formated probe: " + refProbe)
		}
		switch splittedProbe[0] {
		case "swarm":
			// handle swarm probe
			p = &swarm.AverageCPU{Tag: target}
		case "hap":
			// HAproxy probes
			if len(splittedProbe) != 3 {
				return nil, errors.New("")
			}
			p = probe.HAproxy{
				Socket: "/home/maximilien/zenika/haproxy/haproxy.stats",
				Type:   splittedProbe[1],
				Item:   splittedProbe[2],
			}
		case "cmd":
			p = &probe.Command{
				Cmd: rules.Sub(r).GetString("cmd"),
			}
		default:
			return nil, errors.New("Unknown probe " + splittedProbe[0])
		}

		config.Rules = append(config.Rules, rule.FloatValue{
			ServiceName: target,
			Scale:       config.Scalers[scaler], // TODO externalize
			Probe:       p,
			RefreshRate: rules.Sub(r).GetDuration("refresh"),
			Up:          up,
			Down:        down,
		})
	}
	log.Info("Configuration complete !")
	return config, nil
}
