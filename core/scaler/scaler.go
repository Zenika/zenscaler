package scaler

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

// Scaler control the service
type Scaler interface {
	Describe() string
	Up() error
	Down() error
	JSON() ([]byte, error)
}

// MockScaler write "scale up" or "scale down" to stdout
type MockScaler struct{}

// Describe scaler
func (s *MockScaler) Describe() string {
	return "A mock scaler writing to stdout"
}

// JSON encode
func (s *MockScaler) JSON() ([]byte, error) {
	encoded, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// Up mock
func (s *MockScaler) Up() error {
	log.Info("SCALE UP")
	return nil
}

// Down mock
func (s *MockScaler) Down() error {
	log.Info("SCALE DOWN")
	return nil
}
