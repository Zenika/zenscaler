package scaler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/Zenika/zscaler/core"
	"github.com/Zenika/zscaler/core/tls"
)

// ComposeScaler executer docker-compose CLI
type ComposeScaler struct {
	ServiceName       string `json:"service"`
	ConfigFile        string `json:"config"`
	ProjectName       string `json:"project"`
	RunningContainers int    `json:"running"`
	UpperCountLimit   int    `json:"UpperCountLimit"`
	LowerCountLimit   int    `json:"LowerCountLimit"`
	withTLS           bool
	tlsCertsPath      string
	env               []string
}

// NewComposeScaler build a scaler
func NewComposeScaler(name, project, configFilePath string) (*ComposeScaler, error) {
	// TODO need to gather containers, add an INIT ?
	// TODO check for file at provided location
	switch "" { // check missing parameter
	case name:
		return nil, errors.New("No target specified")
	case configFilePath:
		return nil, errors.New("No configuration file path specified")
	case project:
		return nil, errors.New("No project specified")
	}
	cs := &ComposeScaler{
		ServiceName:       name,
		ConfigFile:        configFilePath, // need check
		ProjectName:       project,
		RunningContainers: 3,     // should be discovered
		withTLS:           false, // enforcing default
		UpperCountLimit:   -1,    // default to unlimited
		LowerCountLimit:   1,     // default to one, ensuring service avaibility
	}
	// TLS configuration is checked beforehand but we need to perform additional checks because of docker-compose limitations
	if core.Config.Orchestrator.TLS {
		var err error
		cs.tlsCertsPath, err = tls.CheckTLSConfigPath()
		if err != nil {
			return nil, fmt.Errorf("bad tls config: %s", err)
		}
		cs.withTLS = true
	}
	cs.buildEnv()
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
	upCmd.Env = s.env
	out, err := upCmd.CombinedOutput()
	if err != nil {
		log.Errorf("out: %s\nerr: %s", out, err)
		return err
	}
	s.RunningContainers++
	log.WithFields(log.Fields{
		"service": s.ServiceName,
		"count":   s.RunningContainers,
	}).Infof("scale up")
	return nil
}

// Down using doker compose scale
func (s *ComposeScaler) Down() error {
	if s.RunningContainers < 2 {
		log.WithFields(log.Fields{
			"service": s.ServiceName,
			"count":   s.RunningContainers,
		}).Debug("cannot scale down: minimum count achieved")
		return nil
	}
	// #nosec TODO replace with libcompose API
	downCmd := exec.Command("docker-compose", "-f", s.ConfigFile, "scale", s.ServiceName+"="+strconv.Itoa(s.RunningContainers-1))
	downCmd.Env = s.env
	out, err := downCmd.CombinedOutput()
	if err != nil {
		log.Errorf("out: %s\nerr: %v", out, err)
		return err
	}
	s.RunningContainers--
	log.WithFields(log.Fields{
		"service": s.ServiceName,
		"count":   s.RunningContainers,
	}).Infof("scale down")
	return nil
}

// build commands environnement
func (s *ComposeScaler) buildEnv() {
	s.env = os.Environ()
	s.env = append(s.env, fmt.Sprintf("DOCKER_HOST=%s", core.Config.Orchestrator.Endpoint))
	s.env = append(s.env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", s.ProjectName))
	if s.withTLS { // all certs are in the same path and named correctly
		s.env = append(s.env, fmt.Sprintf("DOCKER_CERT_PATH=%s", s.tlsCertsPath))
		s.env = append(s.env, fmt.Sprintf("DOCKER_TLS_VERIFY=%s", "yes"))
	}
}
