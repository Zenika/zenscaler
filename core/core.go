// Package core provide interface definition and execution flow
package core

import (
	"fmt"

	"github.com/Zenika/zscaler/core/rule"
	"github.com/Zenika/zscaler/core/scaler"

	log "github.com/Sirupsen/logrus"
)

const bufferSize = 10

// Config store the current running configuration
var Config *Configuration

// Configuration holder
type Configuration struct {
	Orchestrator OrchestratorConfig
	Scalers      map[string]scaler.Scaler
	Rules        map[string]rule.Rule
	Errchan      chan error
}

// OrchestratorConfig hold all necessary connection informations
type OrchestratorConfig struct {
	Kind          string // docker, kubernetes, mesos...
	Endpoint      string // http adress:port or unix socket
	TLSCACertPath string // ca cert path, PEM formated
	TLSCertPath   string // user cert path, PEM formated
	TLSKeyPath    string // user private key path, PEM formated
	tlsStatus     bool
}

// CheckTLS for missing certificate and key
func (o OrchestratorConfig) CheckTLS() error {
	o.tlsStatus = false // enforcing default value
	if o.TLSCACertPath+o.TLSCertPath+o.TLSKeyPath == "" {
		return nil
	}
	switch "" {
	case o.TLSCACertPath:
		return fmt.Errorf("tls-cacert path not provided")
	case o.TLSCertPath:
		return fmt.Errorf("tls-cert path not provided")
	case o.TLSKeyPath:
		return fmt.Errorf("tls-key path not provided")
	}
	// all set, TLS seems ok !
	o.tlsStatus = true
	return nil
}

// TLSActivated activation status
func TLSActivated() bool {
	return Config.Orchestrator.tlsStatus
}

// Initialize core module
func (c Configuration) Initialize() {
	c.Errchan = make(chan error, bufferSize)
	c.loop()
}

// event loop
func (c Configuration) loop() {
	log.Debug("Enter control loop...")
	// lanch a watcher on each rule
	for _, r := range c.Rules {
		go rule.Watcher(c.Errchan, r)
	}
	// watch for errors
	for {
		err := <-c.Errchan
		if err != nil {
			log.Errorf("%s", err)
		}
	}
}
