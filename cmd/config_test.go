package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Open all configuration files and parse them
func openAndParseConfig(t *testing.T, path string) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path + "/config.yaml")
	_, err := parseConfig()
	assert.Nil(t, err)
}

func TestConfigDockerComposeCmdHAP(t *testing.T) {
	openAndParseConfig(t, "./../examples/docker-compose/hap")
}
func TestConfigDockerComposeCmdProm(t *testing.T) {
	openAndParseConfig(t, "./../examples/docker-compose/prom")
}
func TestConfigDockerComposeCmdTraefik(t *testing.T) {
	openAndParseConfig(t, "./../examples/docker-compose/traefik")
}
func TestConfigDockerComposeCmdTraefikTLS(t *testing.T) {
	openAndParseConfig(t, "./../examples/docker-compose/traefik-tls")
}
func TestConfigDockerServiceHelloWorld(t *testing.T) {
	openAndParseConfig(t, "./../examples/docker-service/helloworld")
}
