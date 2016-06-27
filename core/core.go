package core

import (
	"time"
	"zscaler/service"
)

// Initialize core module
func Initialize(config *service.Config) {
	// gather parameters
	// do some check
	loop(config)
	// exit cleanup
}

// event loop of the scaler
func loop(config *service.Config) {
	for {
		// TODO this should be launched async
		for _, s := range config.Services {
			if s.Probe.Value() > 0.8 {
				s.Scale.Up()
			}
			if s.Probe.Value() < 0.2 {
				s.Scale.Down()
			}
		}
		time.Sleep(time.Second)
	}
}
