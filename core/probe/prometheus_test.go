package probe

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const data = `node_network_transmit_multicast{device="wlp3s0"} 0
node_network_transmit_multicast{device="wwp0s20u10i8"} 0
# HELP node_network_transmit_packets Network device statistic transmit_packets.
# TYPE node_network_transmit_packets gauge
node_network_transmit_packets{device="br-1cb928ac0d1a"} 0
node_network_transmit_packets{device="docker0"} 0
node_network_transmit_packets{device="docker_gwbridge"} 206
node_network_transmit_packets{device="eno1"} 1.339262e+06
node_network_transmit_packets{device="lo"} 139612
node_network_transmit_packets{device="wlp3s0"} 874
node_network_transmit_packets{device="wwp0s20u10i8"} 0
# HELP node_sockstat_UDP_inuse Number of UDP sockets in state inuse.
# TYPE node_sockstat_UDP_inuse gauge
node_sockstat_UDP_inuse 11`

// Retrieve a value
func TestPrometheusFindValue(t *testing.T) {
	p := Prometheus{
		Key: `node_network_transmit_packets{device="lo"}`,
	}
	val, err := p.findValue(strings.NewReader(data))
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	assert.Equal(t, 139612.0, val)
}

// Check if comments are ignored
func TestPrometheusFindValueCommented(t *testing.T) {
	p := Prometheus{
		Key: `node_sockstat_UDP_inuse`,
	}
	val, err := p.findValue(strings.NewReader(data))
	if err != nil {
		t.Errorf("Error: %s", err)
		t.FailNow()
	}
	assert.Equal(t, 11.0, val)
}

func TestPrometheusGetValue(t *testing.T) {
	p := Prometheus{
		URL: "http://localhost:9100/metrics",
		Key: `node_sockstat_UDP_inuse`,
	}
	val, err := p.Value()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("%f\n", val)
}
