package rule

import (
	"os/exec"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

// Scaler control the service
type Scaler interface {
	Describe() string
	Up() error
	Down() error
}

// MockService is a wrapper for the MockScaler
func MockService(name string) Service {
	return Service{
		Name:  name,
		Scale: new(MockScaler),
	}
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

// ComposeService create a new service
func ComposeService(name string) Service {
	return Service{
		Name: name,
		Scale: &ComposeScaler{
			serviceName:       name,
			configFile:        "/home/maximilien/zenika/zscaler/deploy/swarm/docker-compose.yaml",
			runningContainers: 3,
		},
	}
}

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
