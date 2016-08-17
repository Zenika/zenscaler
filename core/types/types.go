package types

import "time"

// Configuration holder
type Configuration struct {
	Orchestrator OrchestratorConfig
	Scalers      map[string]Scaler
	Rules        map[string]Rule
	Errchan      chan error
}

// OrchestratorConfig hold all necessary connection informations
type OrchestratorConfig struct {
	Engine        string // docker, kubernetes, mesos...
	Endpoint      string // http adress:port or unix socket
	TLSCACertPath string // ca cert path, PEM formated
	TLSCertPath   string // user cert path, PEM formated
	TLSKeyPath    string // user private key path, PEM formated
	TLS           bool   // TLS activation status
}

// A Rule must be able to perform a check
type Rule interface {
	Check() error                 // performe a check on the target and act if needed
	CheckInterval() time.Duration // time to wait between each check
	JSON() ([]byte, error)        // return a JSON output
}

// Scaler control the service
type Scaler interface {
	Describe() string
	Up() error
	Down() error
	JSON() ([]byte, error)
}

// Probe interface
type Probe interface {
	Name() string
	Value() (float64, error)
}
