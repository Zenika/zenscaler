package api

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/types"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const endpoint = "http://localhost:3000/"

// Simple pre-configured wrapper for POST http requests
func postRequest(t *testing.T, route, json string) *http.Response {
	reader := strings.NewReader(json)
	req, err := http.NewRequest("POST", endpoint+route, reader)
	assert.Nil(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	return res
}

// Simple pre-configured wrapper for POST http requests
func getRequest(t *testing.T, route string) *http.Response {
	req, err := http.NewRequest("GET", endpoint+route, nil)
	assert.Nil(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	return res
}

func TestIntegration(t *testing.T) {
	gin.SetMode(gin.ReleaseMode) // supress GIN logging
	// empty config file
	core.Config = &types.Configuration{
		Scalers: map[string]types.Scaler{},
		Rules:   map[string]types.Rule{},
	}
	viper.Set("api-port", ":3000")
	viper.Set("allow-cmd-probe", true)
	// start API server
	go Start()

	// empty queries
	res := getRequest(t, "v1/scalers")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	res = getRequest(t, "v1/rules")
	assert.Equal(t, http.StatusOK, res.StatusCode)

	promRuleJSON := `{
    "probe": "prom.http",
	"probeArgs": {
		"url": "http://localhost:9001/metrics",
		"key": "node_sockstat_UDP_inuse"
	},
    "rule": "prom-rule",
	"service": "dummy-service",
    "scaler": "compose",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
    }`

	// add docker-compose scaler named 'compose'
	composeScalerJSON := `{
    "type":"docker-compose-cmd",
    "name":"compose",
    "args": {
        "service":"whoami",
        "config":"/dummy/path"
        }
    }`
	res = postRequest(t, "v1/scalers", composeScalerJSON)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	// add the same scaler again, expecting error
	res = postRequest(t, "v1/scalers", composeScalerJSON)
	assert.Equal(t, http.StatusConflict, res.StatusCode)
	// query scaler list
	res = getRequest(t, "v1/scalers")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	// query specified scaler
	res = getRequest(t, "v1/scalers/compose")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	// query non-existing scaler
	res = getRequest(t, "v1/scalers/not-found")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	// add a rule
	res = postRequest(t, "v1/rules", promRuleJSON)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	// add the same rule again, expecting error
	res = postRequest(t, "v1/rules", promRuleJSON)
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, res.StatusCode, string(body))

	// add another rule
	cmdRuleJSON := `{
    "probe": "cmd.execute",
	"probeArgs": {
		"cmd": "this-will-be-run-later"
	},
    "rule": "cmd-rule",
	"service": "dummy-service",
    "scaler": "compose",
    "down": "< 1.5",
    "resfreshRate": 10000000000,
    "up": "> 2"
    }`
	res = postRequest(t, "v1/rules", cmdRuleJSON)
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode, string(body))

	// add another scaler
	serviceScalerJSON := `{
    "type":"docker-service",
    "name":"service",
    "args": {
        "service":"whoami"
    }
    }`
	res = postRequest(t, "v1/scalers", serviceScalerJSON)
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode, string(body))

	// query existing rule
	res = getRequest(t, "v1/rules/cmd-rule")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	// query non-existing rule
	res = getRequest(t, "v1/rules/not-found")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
