package scaler

import log "github.com/Sirupsen/logrus"

// Scaler control the service
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
	log.Info("SCALE UP")
	return nil
}

// Down mock
func (s *MockScaler) Down() error {
	log.Info("SCALE DOWN")
	return nil
}

// NewComposeScaler buil a scaler
func NewComposeScaler(name string, configFilePath string) Scaler {
	// TODO need to gather containers, add an INIT ?
	// TODO check for file at provided location
	return &ComposeScaler{
		serviceName:       name,
		configFile:        configFilePath, // need check
		runningContainers: 3,              // should be discovered
	}
}
