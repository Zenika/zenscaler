package service

import (
	"fmt"
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
	Name  string
	Scale Scaler
	Rule  Rule
	Timer time.Duration
}

// Rule interface
type Rule func() bool

// Scaler interface
type Scaler interface {
	Describe() string
	Up() error
	Down() error
}

// MockScaler write "scale up" or "scale down" to stdout
type MockScaler struct{}

// Describe scaler
func (s *MockScaler) Describe() string {
	return "A mock scaler writing to stdout"
}

// Up mock
func (s *MockScaler) Up() error {
	fmt.Println("SCALE UP")
	return nil
}

// Down mock
func (s *MockScaler) Down() error {
	fmt.Println("SCALE DOWN")
	return nil
}
