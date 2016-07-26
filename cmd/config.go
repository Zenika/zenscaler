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

const defaultMapEnties = 5

// DumpConfigCmd definition
var DumpConfigCmd = &cobra.Command{
	Use:   "dumpconfig",
	Short: "Dump parsed config file to stdout",
	Long:  `Check, parse and dump the configuration to the standart output`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := parseConfig()
		if err != nil {
			log.Fatalf("Error in config file: %s", err)
		}
		// TODO pretty output for config
		fmt.Printf("%v", config)
	},
}

func parseConfig() (*core.Configuration, error) {
	// parse config file
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		return nil, fmt.Errorf("Cannot read config file: %s \n", err)
	}

	// global configuration structure
	var config = &core.Configuration{
		Scalers: make(map[string]scaler.Scaler, defaultMapEnties),
		Rules:   make(map[string]rule.Rule, defaultMapEnties),
	}

	// check endpoint
	if viper.GetString("endpoint") == "" {
		return nil, fmt.Errorf("No endpoint specified")
	}

	err = parseScalers(config)
	if err != nil {
		return nil, err
	}

	err = parseRules(config)
	if err != nil {
		return nil, err
	}

	log.Info("Configuration complete !")
	return config, nil
}

func parseScalers(config *core.Configuration) error {
	var scalers = viper.Sub("scalers")
	for _, name := range scalers.AllKeys() {
		log.Info("Add scaler [" + name + "]")
		var s = scalers.Sub(name)
		switch s.GetString("type") {
		case "docker-compose":
			if s.GetString("config") == "" {
				return errors.New("No config specified for docker-compose scaler [" + name + "]")
			}
			if s.GetString("target") == "" {
				return errors.New("No target specified for docker-compose scaler [" + name + "]")
			}
			config.Scalers[name] = scaler.NewComposeScaler(s.GetString("target"), s.GetString("config"))
		case "docker-service":
			if s.GetString("service") == "" {
				return errors.New("No service specified for docker-service scaler [" + name + "]")
			}
			config.Scalers[name] = &scaler.ServiceScaler{
				ServiceID:    s.GetString("service"),
				EngineSocket: viper.GetString("endpoint"),
			}
		case "":
			return fmt.Errorf("Unknown scaler type: %s\n", s.GetString("type"))
		}
	}
	return nil
}

func parseRules(config *core.Configuration) error {
	rules := viper.Sub("rules")
	for _, r := range rules.AllKeys() {
		target := rules.Sub(r).GetString("target")
		log.Info("Add service [" + target + "]")

		// check scaler
		// TODO external function
		scaler := rules.Sub(r).GetString("scaler")
		if config.Scalers[scaler] == nil {
			return errors.New("No scaler specified for rule [" + r + "]")
		}

		// pick probe
		p, err := parseProbe(config, r)
		if err != nil {
			return err
		}

		var floatValueRule = &rule.FloatValue{
			ServiceName:    target,
			Scale:          config.Scalers[scaler], // TODO externalize
			Probe:          p,
			RefreshRate:    rules.Sub(r).GetDuration("refresh"),
			UpDefinition:   rules.Sub(r).GetString("up"),
			DownDefinition: rules.Sub(r).GetString("down"),
		}
		err = floatValueRule.Parse()
		if err != nil {
			return err
		}
		config.Rules[r] = floatValueRule
		fmt.Printf("%v\n", config.Rules[r])
	}
	return nil
}

func parseProbe(config *core.Configuration, r string) (p probe.Probe, err error) {
	rules := viper.Sub("rules")
	target := rules.Sub(r).GetString("target")

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
			return nil, errors.New("hap probe need to be like hap.foo.bar")
		}
		p = &probe.HAproxy{
			Socket: "/home/maximilien/zenika/haproxy/haproxy.stats",
			Type:   splittedProbe[1],
			Item:   splittedProbe[2],
		}
	case "cmd":
		p = &probe.Command{
			Cmd: rules.Sub(r).GetString("cmd"),
		}
	case "prom":
		if splittedProbe[1] == "http" {
			if rules.Sub(r).GetString("url") == "" {
				return nil, errors.New("No url specified for Prometheus probe")
			}
			if rules.Sub(r).GetString("key") == "" {
				return nil, errors.New("No url specified for Prometheus probe")
			}
			p = &probe.Prometheus{
				URL: rules.Sub(r).GetString("url"),
				Key: rules.Sub(r).GetString("key"),
			}
		}
	case "mock":
		p = &probe.DefaultScalingProbe{}
	default:
		return nil, errors.New("Unknown probe " + splittedProbe[0])
	}
	return p, nil
}
