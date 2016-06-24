package service

import "time"

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

// Probe interface
type Probe interface {
	Name() string
	Value() int
}
