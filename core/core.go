// Package core provide interface definition and execution flow
package core

import (
	"zscaler/core/rule"
	"zscaler/core/scaler"

	log "github.com/Sirupsen/logrus"
)

const bufferSize = 10

// Config store the current running configuration
var Config *Configuration

// Configuration holder
type Configuration struct {
	Scalers map[string]scaler.Scaler
	Rules   map[string]rule.Rule
	Errchan chan error
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
			log.Warningf("%s", err)
		}
	}
}
