package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"zscaler/core"
	"zscaler/core/probe"
	"zscaler/core/rule"

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
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
	var rule rule.FloatValue
	err := c.BindJSON(&rule)
	if err != nil {
	core.Config.Rules[rule.RuleName] = &rule

// ruleValidateCheckJSON inputed to create a rule
// It checks missing fields and coherency why the loaded configuration
func ruleValidateCheckJSON(r *rule.FloatValue) error {
	// non-empty fields
	if r.ServiceName == "" {
		return fmt.Errorf("Service name field not set (empty string)")
	}

	// fields related to map entries
	if _, present := core.Config.Rules[r.RuleName]; present {
		return fmt.Errorf("Rule %s already exist", r.RuleName)
	}
	v, present := core.Config.Scalers[r.ScalerID]
	if !present {
		return fmt.Errorf("Scaler %s not found", r.ScalerID)
	}
	r.Scale = v

	// probes parsing
	// see same algo in cmd.parseProbe
	splittedProbe := strings.Split(r.ProbeID, ".")
	if len(splittedProbe) < 2 {
		return fmt.Errorf("Badly formated probe: %s", r.ProbeID)
	}

	// determine probe type and unmarshal ProbeArgs JSON with matching structure
	switch splittedProbe[0] {
	case "swarm":
		// handle swarm probe
		var sp swarm.AverageCPU
		err := json.Unmarshal(r.ProbeArgs, &sp)
		if err != nil {
			return fmt.Errorf("Badly formated JSON for swarm probe: %s", err)
		}
	case "hap":
		// HAproxy probes
		if len(splittedProbe) != 3 {
			return fmt.Errorf("hap probe need to be like hap.<type>.<item>")
		}

		var hp probe.HAproxy
		err := json.Unmarshal(r.ProbeArgs, &hp)
		if err != nil {
			return fmt.Errorf("Badly formated JSON for HAproxy probe: %s", err)
		}
		if hp.Socket == "" {
			return fmt.Errorf("Missing soket in JSON for HAproxy probe")
		}
		hp.Type = splittedProbe[1]
		hp.Item = splittedProbe[2]
	case "cmd":
		var cp probe.Command
		err := json.Unmarshal(r.ProbeArgs, &cp)
		if err != nil {
			return fmt.Errorf("Badly formated JSON for Cmd probe: %s", err)
		}
		r.Probe = cp
	case "prom":
		if splittedProbe[1] == "http" {
			var pp probe.Prometheus
			err := json.Unmarshal(r.ProbeArgs, &pp)
			if err != nil {
				return fmt.Errorf("Badly formated JSON for Prometheus probe: %s", err)
			}
			if pp.URL == "" {
				return fmt.Errorf("Missing URL in JSON for Prometheus probe")
			}
			if pp.Key == "" {
				return fmt.Errorf("Missing Key in JSON for Prometheus probe")
			}
			r.Probe = pp
		}
		return fmt.Errorf("Bad prom encoding type (only http available)")
	case "mock":
		return fmt.Errorf("Mock probe unsupported in API")
	default:
		return fmt.Errorf("Unknown probe %s", splittedProbe[0])
	}
	return nil
}

func patchRule(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func deleteRule(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}
