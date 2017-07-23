package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/probe"
	"github.com/Zenika/zenscaler/core/rule"
	"github.com/Zenika/zenscaler/core/scaler"
	"github.com/Zenika/zenscaler/core/tls"
	"github.com/Zenika/zenscaler/core/types"
	"github.com/Zenika/zenscaler/swarm"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultMapEntries = 5

// DumpConfigCmd definition
var dumpConfigCmd = &cobra.Command{
	Use:   "dumpconfig",
	Short: "Dump parsed config file to stdout",
	Long:  `Check, parse and dump the configuration to the standart output`,
	Run: func(cmd *cobra.Command, args []string) {
		setConfigPath()
		config, err := parseConfig()
		if err != nil {
			log.Fatalf("Error in config file: %s", err)
		}
		// TODO pretty output for config
		fmt.Printf("%v", config)
	},
}

func setConfigPath() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // look for config in the working directory
}

func parseConfig() (*types.Configuration, error) {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return nil, fmt.Errorf("cannot read config file: %s", err)
	}

	// global configuration structure
	var config = &types.Configuration{
		Scalers: make(map[string]types.Scaler, defaultMapEntries),
		Rules:   make(map[string]types.Rule, defaultMapEntries),
	}
	// set it as global
	core.Config = config

	err = parseOrchestrator(config)
	if err != nil {
		return nil, err
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

func parseOrchestrator(config *types.Configuration) error {
	sub := viper.Sub("orchestrator")
	if sub == nil {
		return fmt.Errorf("'orchestrator' section required, missing")
	}
	config.Orchestrator = types.OrchestratorConfig{
		Engine:        sub.GetString("engine"),
		Endpoint:      sub.GetString("endpoint"),
		TLSCACertPath: sub.GetString("tls-cacert"),
		TLSCertPath:   sub.GetString("tls-cert"),
		TLSKeyPath:    sub.GetString("tls-key"),
	}
	// check endpoint
	if config.Orchestrator.Endpoint == "" && config.Orchestrator.Engine == "docker" {
		config.Orchestrator.Endpoint = "unix:///var/run/docker.sock" // default socket
	}
	if config.Orchestrator.Endpoint == "" {
		return fmt.Errorf("No endpoint specified")
	}
	// check TLS requirements
	return tls.CheckTLS(config)
}

func parseScalers(config *types.Configuration) error {
	var scalers = viper.Sub("scalers")
	for _, name := range scalers.AllKeys() {
		log.Info("Add scaler [" + name + "]")
		var s = scalers.Sub(name)
		switch s.GetString("type") {
		case "docker-compose-cmd":
			cs, err := parseScalerDockerComposeCmd(s)
			if err != nil {
				return fmt.Errorf("cannot create docker-compose-cmd scaler [%s]: %s", name, err)
			}
			config.Scalers[name] = cs
		case "docker-service":
			ss, err := parseScalerDockerService(s)
			if err != nil {
				return fmt.Errorf("cannot create docker-service scaler [%s]: %s", name, err)
			}
			config.Scalers[name] = ss
		case "":
			return fmt.Errorf("unknown scaler type: %s", s.GetString("type"))
		}
	}
	return nil
}

func parseScalerDockerComposeCmd(s *viper.Viper) (*scaler.ComposeCmdScaler, error) {
	cs, err := scaler.NewComposeCmdScaler(s.GetString("target"), s.GetString("project"), s.GetString("config"))
	if err != nil {
		return nil, err
	}
	// set optional parameter
	if s.IsSet("upper_count_limit") {
		cs.UpperCountLimit = uint64(s.GetInt("upper_count_limit"))
	}
	if s.IsSet("lower_count_limit") {
		cs.LowerCountLimit = uint64(s.GetInt("lower_count_limit"))
	}
	return cs, nil
}

func parseScalerDockerService(s *viper.Viper) (*scaler.ServiceScaler, error) {
	if s.GetString("service") == "" {
		return nil, errors.New("No service specified")
	}
	ss := &scaler.ServiceScaler{
		ServiceID:       s.GetString("service"),
		EngineSocket:    viper.GetString("endpoint"),
		LowerCountLimit: 1,
		UpperCountLimit: 0,
	}
	if s.IsSet("upper_count_limit") {
		ss.UpperCountLimit = uint64(s.GetInt("upper_count_limit"))
	}
	if s.IsSet("lower_count_limit") {
		ss.LowerCountLimit = uint64(s.GetInt("lower_count_limit"))
	}
	return ss, nil
}

func parseRules(config *types.Configuration) error {
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
			Scale:          config.Scalers[scaler],
			ScalerID:       scaler,
			Probe:          p,
			ProbeID:        rules.Sub(r).GetString("probe"),
			RefreshRate:    rules.Sub(r).GetDuration("refresh"),
			UpDefinition:   rules.Sub(r).GetString("up"),
			DownDefinition: rules.Sub(r).GetString("down"),
		}
		err = floatValueRule.Parse()
		if err != nil {
			return err
		}
		config.Rules[r] = floatValueRule
	}
	return nil
}

func parseProbe(config *types.Configuration, r string) (p types.Probe, err error) {
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
		p, err = parseProbeHAP(rules, r, splittedProbe)
	case "cmd":
		p = &probe.Command{
			Cmd: rules.Sub(r).GetString("cmd"),
		}
	case "prom":
		p, err = parseProbeProm(rules, r, splittedProbe)
	case "mock":
		p = &probe.DefaultScalingProbe{}
	default:
		return nil, errors.New("Unknown probe " + splittedProbe[0])
	}
	return p, nil
}

func parseProbeHAP(rules *viper.Viper, r string, splittedProbe []string) (types.Probe, error) {
	if len(splittedProbe) != 3 {
		return nil, errors.New("hap probe need to be like hap.foo.bar")
	}
	if rules.Sub(r).GetString("ha-socket") == "" {
		return nil, errors.New("No hap stat socket specified for " + r + " probe")
	}
	return &probe.HAproxy{
		Socket: rules.Sub(r).GetString("ha-socket"),
		Type:   splittedProbe[1],
		Item:   splittedProbe[2],
	}, nil
}

func parseProbeProm(rules *viper.Viper, r string, splittedProbe []string) (types.Probe, error) {
	if splittedProbe[1] == "http" {
		if rules.Sub(r).GetString("url") == "" {
			return nil, errors.New("No url specified for Prometheus probe")
		}
		if rules.Sub(r).GetString("key") == "" {
			return nil, errors.New("No url specified for Prometheus probe")
		}
		return &probe.Prometheus{
			URL: rules.Sub(r).GetString("url"),
			Key: rules.Sub(r).GetString("key"),
		}, nil
	}
	return nil, fmt.Errorf("Unknow prometheus probe type \"%s\"", splittedProbe[1])
}
