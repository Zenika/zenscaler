package core

import (
	"time"
	"zscaler/provider/swarm"
	"zscaler/service"
)

// Initialize core module
func Initialize(config *service.Config) {
	// preliminary check
	checked = swarm.CheckServices(config.Services)
	loop(config)
	// exit cleanup
}

// event loop
func loop(config *service.Config) {
	for {

	}
}

func serviceWatcher(errchan chan error, service service.Service) {

	time.Sleep(service.Timer)
}
