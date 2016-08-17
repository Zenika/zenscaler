package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Open all configuration files and parse them
func openAndParseConfig(t *testing.T, path string) {
	viper.AddConfigPath(path)
	_, err := parseConfig()
	assert.Nil(t, err)
}

func TestConfigDockerComposeHAP(t *testing.T) {
	openAndParseConfig(t, "./../examples/docker-compose/hap")
}
func TestConfigDockerComposeProm(t *testing.T) {
	openAndParseConfig(t, "./../examples/docker-compose/prom")
}
func TestConfigDockerComposeTraefik(t *testing.T) {
	openAndParseConfig(t, "./../examples/docker-compose/traefik")
}
func TestConfigDockerServiceHelloWorld(t *testing.T) {
	openAndParseConfig(t, "./../examples/docker-service/helloworld")
}
