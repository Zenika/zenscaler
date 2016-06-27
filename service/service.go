package service

import (
	"time"
	"zscaler/probe"
)

// Config parameters
type Config struct {
	Services []Service
	Probes   map[string]probe.Probe
}

// Service descibes the object to scale
type Service struct {
	Name    string
	Scale   Scaler
	Timeout time.Duration
}

// Scaler interface
type Scaler interface {
	Up() error
	Down() error
}
