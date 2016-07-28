package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"zscaler/core"
	"zscaler/core/probe"
	"zscaler/core/rule"
	"zscaler/swarm"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func getRules(c *gin.Context) {
	var ruleNames = make([]string, 0)
	for k := range core.Config.Rules {
		ruleNames = append(ruleNames, k)
	}
	c.JSON(http.StatusOK, gin.H{
		"rules": ruleNames,
	})
}

func getRule(c *gin.Context) {
	name := c.Param("name")
	// does this rule exist ?
	if rule, ok := core.Config.Rules[name]; ok {
		encoded, err := rule.JSON()
		if err == nil {
			c.Data(http.StatusOK, "application/json", encoded)
			return
		}
		log.Errorf("Encode error %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": name + " does not encode to JSON",
		})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": name + " not found",
	})
}

func createRule(c *gin.Context) {
	var floatRuleBuilder FloatValueBuilder
	err := c.BindJSON(&floatRuleBuilder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON object cannot be parsed: " + err.Error(),
		})
		return
	}
	fvRule, err := floatRuleBuilder.Build()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON failed validation check: " + err.Error(),
		})
		return
	}
	// TODO check data race
	core.Config.Rules[fvRule.RuleName] = fvRule
	go rule.Watcher(core.Config.Errchan, fvRule)

	c.JSON(http.StatusCreated, gin.H{
		"rule": fvRule.RuleName,
	})
}

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
// It checks missing fields and coherency why the loaded configuration
func (r *FloatValueBuilder) Build() (*rule.FloatValue, error) {
	var fv rule.FloatValue
	// non-empty fields
	if r.ServiceName == "" {
		return nil, fmt.Errorf("Service name field not set (empty string)")
	}

	// fields related to map entries
	if _, present := core.Config.Rules[r.RuleName]; present {
		return nil, fmt.Errorf("Rule %s already exist", r.RuleName)
	}
	v, present := core.Config.Scalers[r.ScalerID]
	if !present {
		return nil, fmt.Errorf("Scaler %s not found", r.ScalerID)
	}

	fv.RuleName = r.RuleName
	fv.ServiceName = r.ServiceName
	fv.ScalerID = r.ScalerID
	fv.Scale = v
	fv.ProbeID = r.ProbeID

	// probes parsing
	// see same algo in cmd.parseProbe
	splittedProbe := strings.Split(r.ProbeID, ".")
	if len(splittedProbe) < 2 {
		return nil, fmt.Errorf("Badly formated probe: %s", r.ProbeID)
	}

	// determine probe type and unmarshal ProbeArgs JSON with matching structure
	switch splittedProbe[0] {
	case "swarm":
		// handle swarm probe
		var sp swarm.AverageCPU
		err := json.Unmarshal(r.ProbeArgs, &sp)
		if err != nil {
			return nil, fmt.Errorf("Badly formated JSON for swarm probe: %s", err)
		}
		fv.Probe = sp
	case "hap":
		// HAproxy probes
		if len(splittedProbe) != 3 {
			return nil, fmt.Errorf("hap probe need to be like hap.<type>.<item>")
		}

		var hp probe.HAproxy
		err := json.Unmarshal(r.ProbeArgs, &hp)
		if err != nil {
			return nil, fmt.Errorf("Badly formated JSON for HAproxy probe: %s", err)
		}
		if hp.Socket == "" {
			return nil, fmt.Errorf("Missing soket in JSON for HAproxy probe")
		}
		hp.Type = splittedProbe[1]
		hp.Item = splittedProbe[2]
		fv.Probe = hp
	case "cmd":
		var cp probe.Command
		err := json.Unmarshal(r.ProbeArgs, &cp)
		if err != nil {
			return nil, fmt.Errorf("Badly formated JSON for Cmd probe: %s", err)
		}
		fv.Probe = cp
	case "prom":
		if splittedProbe[1] == "http" {
			var pp probe.Prometheus
			err := json.Unmarshal(r.ProbeArgs, &pp)
			if err != nil {
				return nil, fmt.Errorf("Badly formated JSON for Prometheus probe: %s", err)
			}
			if pp.URL == "" {
				return nil, fmt.Errorf("Missing URL in JSON for Prometheus probe")
			}
			if pp.Key == "" {
				return nil, fmt.Errorf("Missing Key in JSON for Prometheus probe")
			}
			fv.Probe = pp
		}
		return nil, fmt.Errorf("Bad prom encoding type (only http available)")
	case "mock":
		return nil, fmt.Errorf("Mock probe unsupported in API")
	default:
		return nil, fmt.Errorf("Unknown probe %s", splittedProbe[0])
	}

	fv.RefreshRate = r.RefreshRate
	fv.UpDefinition = r.UpDefinition
	fv.DownDefinition = r.DownDefinition

	err := fv.Parse()
	if err != nil {
		return nil, fmt.Errorf("Error parsing rules: %s", err)
	}

	return &fv, nil
}

func patchRule(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func deleteRule(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}
