package api

import (
	"encoding/json"
	"testing"

	"github.com/Zenika/zenscaler/core"
	"github.com/stretchr/testify/assert"
)

func configAndBuildScaler(t *testing.T, input string) error {
	var sb ScalerBuilder
	core.Config = MockConf
	err := json.Unmarshal([]byte(input), &sb)
	failIfErr(t, err)
	_, err = sb.Build()
	return err
}

func TestBuildComposeCmdScaler(t *testing.T) {
	const input = `{
    "type":"docker-compose-cmd",
    "name":"testing",
    "args": {
        "service":"whoami",
        "project":"test",
        "config":"/dummy/path",
		"upperCountLimit":0,
		"lowerCountLimit":1
    }
}`
	err := configAndBuildScaler(t, input)
	assert.Nil(t, err)
}

func TestBuildComposeCmdScalerBadType(t *testing.T) {
	const input = `{
    "type":"badtype",
    "name":"testing",
    "args": {
        "service":"whoami",
        "project":"test",
        "config":"/dummy/path",
		"upperCountLimit":0,
		"lowerCountLimit":1
    }
}`
	err := configAndBuildScaler(t, input)
	assert.Error(t, err)
}
