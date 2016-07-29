package api

import (
	"encoding/json"
	"testing"
	"zscaler/core"
	"zscaler/core/rule"
	"zscaler/core/scaler"

	"github.com/stretchr/testify/assert"
)

var MockConf = &core.Configuration{
	Scalers: map[string]scaler.Scaler{
		"whoami-compose": &scaler.ComposeScaler{
			ServiceName: "whoami",
		},
	},
	Rules: map[string]rule.Rule{},
}

func setupMockAndDecode(t *testing.T, input string, target interface{}) {
	core.Config = MockConf
	err := json.Unmarshal([]byte(input), &target)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
}

func configAndBuildRule(t *testing.T, input string) error {
	var frb FloatValueBuilder
	setupMockAndDecode(t, input, frb)
	_, err := frb.Build()
	return err
}

func TestCreateComposeRuleSwarmProbe(t *testing.T) {
	const input = `{
    "probe": "swarm.cpu_average",
    "rule": "custom",
    "scaler": "whoami-compose",
    "service": "whoami",
    "probeArgs": {
        "Tag": "whoami"
    },
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	err := configAndBuildRule(t, input)
	assert.Nil(t, err)
}

func TestCreateComposeRuleMissingScaler(t *testing.T) {
	const input = `{
    "probe": "swarm.cpu_average",
    "rule": "custom",
    "scaler": "missing",
    "service": "whoami",
    "probeArgs": {
        "Tag": "whoami"
    },
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	err := configAndBuildRule(t, input)
	assert.Error(t, err)
}

func TestCreateComposeRuleMissingDockerProbeArgs(t *testing.T) {
	const input = `{
    "probe": "swarm.cpu_average",
    "rule": "custom",
    "scaler": "whoami-compose",
    "service": "whoami",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	err := configAndBuildRule(t, input)
	assert.Error(t, err)
}

func TestCreateComposeRuleBadProbe(t *testing.T) {
	const input = `{
    "probe": "bad.probe",
    "rule": "custom",
    "scaler": "whoami-compose",
    "service": "whoami",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	err := configAndBuildRule(t, input)
	assert.Error(t, err)
}
