package scaler

import (
	"os/exec"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

// ComposeScaler executer docker-compose CLI
type ComposeScaler struct {
	serviceName       string
	configFile        string
	runningContainers int
}

// Describe scaler
func (s *ComposeScaler) Describe() string {
	return "Exec docker-compose scaler"
}

// Up using doker compose scale
func (s *ComposeScaler) Up() error {
	upCmd := exec.Command("docker-compose", "-f", s.configFile, "scale", s.serviceName+"="+strconv.Itoa(s.runningContainers+1))
	log.Infof("Scale "+s.serviceName+" up to %d", s.runningContainers+1)
	out, err := upCmd.CombinedOutput()
	if err != nil {
		log.Errorf("out: %s\nerr: %v", out, err)
		return err
	}
	s.runningContainers++
	return nil
}

// Down using doker compose scale
func (s *ComposeScaler) Down() error {
	if s.runningContainers < 2 {
		log.Debug("Cannot scale down below one container")
		return nil
	}
	downCmd := exec.Command("docker-compose", "-f", s.configFile, "scale", s.serviceName+"="+strconv.Itoa(s.runningContainers-1))
	log.Infof("Scale "+s.serviceName+" down to %d", s.runningContainers-1)
	out, err := downCmd.CombinedOutput()
	if err != nil {
		log.Errorf("out: %s\nerr: %v", out, err)
		return err
	}
	s.runningContainers--
	return nil
}
