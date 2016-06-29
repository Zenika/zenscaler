// Package core provide interface definition and execution flow
package core

import (
	"time"
	"zscaler/core/service"
	"zscaler/swarm"
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

}

func serviceWatcher(errchan chan error, service service.Service) {

	time.Sleep(service.Timer)
}
