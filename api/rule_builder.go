package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/probe"
	"github.com/Zenika/zenscaler/core/rule"
	"github.com/Zenika/zenscaler/swarm"
	"github.com/spf13/viper"
)

// FloatValueBuilder contains all the information to build a rule
type FloatValueBuilder struct {
	RuleName       string          `json:"rule"`
	ServiceName    string          `json:"service"`
	ScalerID       string          `json:"scaler"`
	ProbeID        string          `json:"probe"`
	ProbeArgs      json.RawMessage `json:"probeArgs"`
	RefreshRate    time.Duration   `json:"resfreshRate"`
	UpDefinition   string          `json:"up"`
	DownDefinition string          `json:"down"`
}

// Build validate inputed data and return a FloatValue rule
// It checks missing fields and coherence why the loaded configuration
func (builder *FloatValueBuilder) Build() (*rule.FloatValue, error) {
	var fv rule.FloatValue

	err := builder.checkEmptyFields()
	if err != nil {
		return nil, err
	}

	v, present := core.Config.Scalers[builder.ScalerID]
	if !present {
		return nil, fmt.Errorf("Scaler %s not found", builder.ScalerID)
	}

	fv.RuleName = builder.RuleName
	fv.ServiceName = builder.ServiceName
	fv.ScalerID = builder.ScalerID
	fv.Scale = v
	fv.ProbeID = builder.ProbeID

	err = builder.parseProbe(&fv)
	if err != nil {
		return nil, err
	}

	fv.RefreshRate = builder.RefreshRate
	fv.UpDefinition = builder.UpDefinition
	fv.DownDefinition = builder.DownDefinition

	err = fv.Parse()
	if err != nil {
		return nil, fmt.Errorf("Error parsing rules: %s", err)
	}

	return &fv, nil
}

func (builder *FloatValueBuilder) checkEmptyFields() error {
	// non-empty fields
	if builder.ServiceName == "" {
		return fmt.Errorf("Service name field not set (empty string)")
	}

	// fields related to map entries
	if builder.RuleName == "" {
		return fmt.Errorf("No specified RuleName")
	}
	return nil
}

// probes parsing
// see same algo in cmd.parseProbe
func (builder *FloatValueBuilder) parseProbe(fv *rule.FloatValue) (err error) {

	splittedProbe := strings.Split(builder.ProbeID, ".")
	if len(splittedProbe) < 2 {
		return fmt.Errorf("Badly formated probe: %s", builder.ProbeID)
	}

	// determine probe type and unmarshal ProbeArgs JSON with matching structure
	switch splittedProbe[0] {
	case "swarm":
		err = builder.parseProbeSwarm(fv)
	case "hap":
		err = builder.parseProbeHAP(fv, splittedProbe)
	case "cmd":
		err = builder.parseProbeCmd(fv)
	case "prom":
		err = builder.parseProbeProm(fv, splittedProbe)
	case "mock":
		return fmt.Errorf("Mock probe unsupported in API")
	default:
		return fmt.Errorf("Unknown probe %s", splittedProbe[0])
	}
	return err
}

func (builder *FloatValueBuilder) parseProbeHAP(fv *rule.FloatValue, splittedProbe []string) (err error) {
	if len(splittedProbe) != 3 {
		return fmt.Errorf("hap probe need to be like hap.<type>.<item>")
	}

	var hp probe.HAproxy
	err = json.Unmarshal(builder.ProbeArgs, &hp) // only socket field is required
	if err != nil {
		return fmt.Errorf("Badly formated JSON for HAproxy probe: %s", err)
	}
	if hp.Socket == "" {
		return fmt.Errorf("Missing soket in JSON for HAproxy probe")
	}
	hp.Type = splittedProbe[1]
	hp.Item = splittedProbe[2]
	fv.Probe = hp
	return nil
}

func (builder *FloatValueBuilder) parseProbeProm(fv *rule.FloatValue, splittedProbe []string) (err error) {
	if splittedProbe[1] == "http" {
		var pp probe.Prometheus
		err = json.Unmarshal(builder.ProbeArgs, &pp)
		if err != nil {
			return fmt.Errorf("Badly formated JSON for Prometheus probe: %s", err)
		}
		if pp.URL == "" {
			return fmt.Errorf("Missing URL in JSON for Prometheus probe")
		}
		if pp.Key == "" {
			return fmt.Errorf("Missing Key in JSON for Prometheus probe")
		}
		fv.Probe = pp
		return nil
	}
	return fmt.Errorf("Bad prom encoding type (only http available)")
}

func (builder *FloatValueBuilder) parseProbeSwarm(fv *rule.FloatValue) error {
	var sp swarm.AverageCPU
	err := json.Unmarshal(builder.ProbeArgs, &sp)
	if err != nil {
		return fmt.Errorf("Badly formated JSON for swarm probe: %s", err)
	}
	fv.Probe = sp
	return nil
}

func (builder *FloatValueBuilder) parseProbeCmd(fv *rule.FloatValue) (err error) {
	if !viper.GetBool("allow-cmd-probe") {
		return fmt.Errorf("Cannot create cmd probe using API if started without --allow-cmd-probe")
	}

	var cp probe.Command
	err = json.Unmarshal(builder.ProbeArgs, &cp)
	if err != nil {
		return fmt.Errorf("Badly formated JSON for Cmd probe: %s", err)
	}
	fv.Probe = cp
	return nil
}
