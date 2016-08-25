package api

import (
	"encoding/json"
	"testing"

	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/scaler"
	"github.com/Zenika/zenscaler/core/types"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

var MockConf = &types.Configuration{
	Scalers: map[string]types.Scaler{
		"whoami-compose": &scaler.ComposeCmdScaler{
			ServiceName: "whoami",
		},
	},
	Rules: map[string]types.Rule{},
}

func failIfErr(t *testing.T, err error) {
	if err != nil {
		assert.FailNow(t, err.Error())
	}
}

func configAndBuildRule(t *testing.T, input string) error {
	core.Config = MockConf
	var frb FloatValueBuilder
	err := json.Unmarshal([]byte(input), &frb)
	failIfErr(t, err)
	_, err = frb.Build()
	return err
}

func TestBuildRuleSwarmProbe(t *testing.T) {
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

func TestBuildRuleMissingScaler(t *testing.T) {
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

func TestBuildRuleMissingDockerProbeArgs(t *testing.T) {
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

func TestBuildRuleBadProbeType(t *testing.T) {
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

func TestBuildRuleBadProbeFormat(t *testing.T) {
	const input = `{
    "probe": "bad-format",
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

func TestBuildMissingServiceName(t *testing.T) {
	const input = `{
    "probe": "bad.probe",
    "rule": "custom",
    "scaler": "whoami-compose",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	err := configAndBuildRule(t, input)
	assert.Error(t, err)
}

func TestBuildHapProbe(t *testing.T) {
	const input = `{
    "probe": "hap.backend.rtime",
	"probeArgs": {
		"socket": "dummy/uri/that/will/be/parsed/later"
	},
    "rule": "custom",
	"service": "whoami",
    "scaler": "whoami-compose",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	err := configAndBuildRule(t, input)
	assert.Nil(t, err)
}

func TestBuildHapProbeDescriptionMissingEnd(t *testing.T) {
	const input = `{
    "probe": "hap.backend",
    "rule": "custom",
	"service": "whoami",
    "scaler": "whoami-compose",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	err := configAndBuildRule(t, input)
	assert.Error(t, err)
}

func TestBuildHapProbeNoArgs(t *testing.T) {
	const input = `{
    "probe": "hap.backend.rtime",
    "rule": "custom",
	"service": "whoami",
    "scaler": "whoami-compose",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	err := configAndBuildRule(t, input)
	assert.Error(t, err)
}

func TestBuildCmdProbe(t *testing.T) {
	const input = `{
    "probe": "cmd.execute",
	"probeArgs": {
		"cmd": "this-will-be-run-later"
	},
    "rule": "custom",
	"service": "whoami",
    "scaler": "whoami-compose",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	viper.Set("allow-cmd-probe", true)
	err := configAndBuildRule(t, input)
	assert.Nil(t, err)
	viper.Set("allow-cmd-probe", false)
	err = configAndBuildRule(t, input)
	assert.Error(t, err)
}

func TestBuildPromProbe(t *testing.T) {
	const input = `{
    "probe": "prom.http",
	"probeArgs": {
		"url": "http://localhost:9001/metrics",
		"key": "node_sockstat_UDP_inuse"
	},
    "rule": "custom",
	"service": "whoami",
    "scaler": "whoami-compose",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
}`
	err := configAndBuildRule(t, input)
	assert.Nil(t, err)
}
