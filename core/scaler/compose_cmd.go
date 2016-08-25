package scaler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/tls"
)

// ComposeCmdScaler executer docker-compose CLI
type ComposeCmdScaler struct {
	ServiceName       string `json:"service"`
	ConfigFile        string `json:"config"`
	ProjectName       string `json:"project"`
	RunningContainers uint64 `json:"running"`
	UpperCountLimit   uint64 `json:"upperCountLimit"`
	LowerCountLimit   uint64 `json:"lowerCountLimit"`
	withTLS           bool
	tlsCertsPath      string
	env               []string
}

// NewComposeCmdScaler build a scaler
func NewComposeCmdScaler(name, project, configFilePath string) (*ComposeCmdScaler, error) {
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
	cs := &ComposeCmdScaler{
		ServiceName:       name,
		ConfigFile:        configFilePath, // need check
		ProjectName:       project,
		RunningContainers: 2,     // should be discovered
		withTLS:           false, // enforcing default
		UpperCountLimit:   0,     // default to unlimited
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
func (s *ComposeCmdScaler) Describe() string {
	return "Exec docker-compose scaler"
}

// JSON encoding
func (s *ComposeCmdScaler) JSON() ([]byte, error) {
	encoded, err := json.Marshal(*s)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// Up using doker compose scale
func (s *ComposeCmdScaler) Up() error {
	upRunningContainers := s.RunningContainers + 1
	if upRunningContainers > s.UpperCountLimit && s.UpperCountLimit != 0 {
		s.getLogger().Debug("cannot scale up: maximum count achieved")
		return nil
	}
	err := s.execComposeCmd(upRunningContainers)
	if err != nil {
		return err
	}
	s.RunningContainers++
	s.getLogger().Infof("scale up")
	return nil
}

// Down using doker compose scale
func (s *ComposeCmdScaler) Down() (err error) {
	downRunningContainers := s.RunningContainers - 1
	if downRunningContainers < s.LowerCountLimit {
		s.getLogger().Debug("cannot scale down: minimum count achieved")
		return
	}
	err = s.execComposeCmd(downRunningContainers)
	if err != nil {
		return
	}
	s.RunningContainers--
	s.getLogger().Infof("scale down")
	return
}

// build commands environnement
func (s *ComposeCmdScaler) buildEnv() {
	s.env = os.Environ()
	s.env = append(s.env, fmt.Sprintf("DOCKER_HOST=%s", core.Config.Orchestrator.Endpoint))
	s.env = append(s.env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", s.ProjectName))
	if s.withTLS { // all certs are in the same path and named correctly
		s.env = append(s.env, fmt.Sprintf("DOCKER_CERT_PATH=%s", s.tlsCertsPath))
		s.env = append(s.env, fmt.Sprintf("DOCKER_TLS_VERIFY=%s", "yes"))
	}
}

func (s *ComposeCmdScaler) execComposeCmd(targetRunningContainers uint64) error {
	// #nosec TODO replace with libcompose API
	downCmd := exec.Command("docker-compose", "-f", s.ConfigFile, "scale", s.ServiceName+"="+fmt.Sprintf("%d", targetRunningContainers))
	downCmd.Env = s.env
	out, err := downCmd.CombinedOutput()
	if err != nil {
		log.Errorf("out: %s\nerr: %v", out, err)
		return err
	}
	return nil
}

func (s *ComposeCmdScaler) getLogger() *log.Entry {
	return log.WithFields(log.Fields{
		"service": s.ServiceName,
		"count":   s.RunningContainers,
	})
}
