package scaler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeScalerCount(t *testing.T) {
	_, err := NewComposeScaler("whomami", "traefik", "/home/maximilien/.go/src/github.com/Zenika/zscaler/examples/docker-compose/traefik/docker-compose.yaml")
	assert.Nil(t, err)
}
