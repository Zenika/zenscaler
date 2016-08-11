package scaler

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/Zenika/zscaler/core"
	"github.com/Zenika/zscaler/core/tls"
	"github.com/Zenika/zscaler/core/types"
)

// ComposeScaler executer docker-compose CLI
type ComposeScaler struct {
	ServiceName       string `json:"service"`
	ConfigFile        string `json:"config"`
	RunningContainers int    `json:"running"`
	withTLS           bool
	tlsCertsPath      string
}

// NewComposeScaler build a scaler
func NewComposeScaler(name string, ConfigFilePath string) (types.Scaler, error) {
	// TODO need to gather containers, add an INIT ?
	// TODO check for file at provided location
	cs := &ComposeScaler{
		ServiceName:       name,
		ConfigFile:        ConfigFilePath, // need check
		RunningContainers: 3,              // should be discovered
		withTLS:           false,          // enforcing default
	}
	// TLS configuration
	if core.Config.Orchestrator.TLS {
		var err error
		cs.tlsCertsPath, err = tls.CheckTLSConfigPath()
		if err != nil {
			return nil, fmt.Errorf("bad tls config: %s", err)
		}
		cs.withTLS = true
	}
	return cs, nil
}

// Describe scaler
func (s *ComposeScaler) Describe() string {
	return "Exec docker-compose scaler"
}

// JSON encoding
func (s *ComposeScaler) JSON() ([]byte, error) {
	encoded, err := json.Marshal(*s)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// Up using doker compose scale
func (s *ComposeScaler) Up() error {
	// #nosec TODO replace with libcompose API
	upCmd := exec.Command("docker-compose", "-f", s.ConfigFile, "scale", s.ServiceName+"="+strconv.Itoa(s.RunningContainers+1))
	log.Infof("Scale "+s.ServiceName+" up to %d", s.RunningContainers+1)
	out, err := upCmd.CombinedOutput()
	if err != nil {
		log.Errorf("out: %s\nerr: %v", out, err)
		return err
	}
	s.RunningContainers++
	return nil
}

// Down using doker compose scale
func (s *ComposeScaler) Down() error {
	if s.RunningContainers < 2 {
		log.Debug("Cannot scale down below one container")
		return nil
	}
	// #nosec TODO replace with libcompose API
	downCmd := exec.Command("docker-compose", "-f", s.ConfigFile, "scale", s.ServiceName+"="+strconv.Itoa(s.RunningContainers-1))
	log.Infof("Scale "+s.ServiceName+" down to %d", s.RunningContainers-1)
	out, err := downCmd.CombinedOutput()
	if err != nil {
		log.Errorf("out: %s\nerr: %v", out, err)
		return err
	}
	s.RunningContainers--
	return nil
}
