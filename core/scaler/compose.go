package scaler

import (
	"encoding/json"
	"os/exec"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

// ComposeScaler executer docker-compose CLI
type ComposeScaler struct {
	ServiceName       string `json:"service"`
	ConfigFile        string `json:"config"`
	RunningContainers int    `json:"running"`
}

// NewComposeScaler buil a scaler
func NewComposeScaler(name string, ConfigFilePath string) Scaler {
	// TODO need to gather containers, add an INIT ?
	// TODO check for file at provided location
	return &ComposeScaler{
		ServiceName:       name,
		ConfigFile:        ConfigFilePath, // need check
		RunningContainers: 3,              // should be discovered
	}
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
